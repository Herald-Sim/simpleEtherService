package api

import (
	"account"
	"appcontext"
	"contract"
	"driver"
	"encoding/json"
	"erctoken"
	"eth"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"router"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var MySQL driver.SQLinfo
var MyAES driver.AESObject

var TokenAddress string
var MyToken erctoken.ERC20Token

var ContractAddress string
var MyContract contract.MyContract

func init() {
	MySQL.InitSQL()
	// For prevent MariaDB EOF Error
	MySQL.SQLClient.SetMaxIdleConns(0)

	MyAES.InitAES()

	// HRToken
	TokenAddress = "0x25904468f630ad9a3937be11e96a6ded913abc71"
	MyToken.InitToken(TokenAddress)

	// My contract account
	ContractAddress = "0x8364fd2B18B15c0ABd7A86A042c532964FcDb8B1"
	MyContract.InitContract(ContractAddress)
}

type MiddleWare func(next router.HandlerFunc) router.HandlerFunc

// SignUp : read user info(get from client) and make Ether account and save DB
func SignUp(c *router.Context) {
	retJSON := appcontext.ReturnJSON{}
	dbobject := appcontext.UserObject{}

	recvData, _ := ioutil.ReadAll(c.Request.Body)

	json.Unmarshal(recvData, &dbobject)
	fmt.Println("Receive from client: ", dbobject)

	// Insert DB
	userInsert, err := MySQL.SQLClient.Exec("INSERT `myEtherWallet`.`user` (`id`, `passwd`) VALUES (?, ?);",
		dbobject.ID,
		dbobject.Passwd)

	// for recover handler
	defer func() {
		s := recover()
		fmt.Println(s)
	}()

	if err != nil {
		c.ResponseWriter.WriteHeader(500)
		retJSON.Status = "fail"
		jsonData, _ := json.Marshal(retJSON)

		c.ResponseWriter.Write(jsonData)

		fmt.Println(err)
		panic("INSERT ERROR")
	}

	n, err := userInsert.RowsAffected()
	if n == 1 {
		// created successful
		c.ResponseWriter.WriteHeader(201)

		retJSON.Status = "success"
		jsonData, _ := json.Marshal(retJSON)

		c.ResponseWriter.Write(jsonData)

		userPriv, userPub := account.GenKey()
		walletInsert, err := MySQL.SQLClient.Exec("INSERT `myEtherWallet`.`wallet` (`pubKey`, `privKey`, `ownerId`) VALUES (?, ?, ?);",
			userPub,
			userPriv,
			dbobject.ID,
		)
		if err != nil {
			log.Fatal(err)
			panic("walletInsert ERROR")
		}

		n, err := walletInsert.RowsAffected()
		if n == 1 {
			fmt.Println("walletInsert success")
		} else {
			panic("Need Check walletInsert RowsAffected")
		}
	} else {
		// internal server error
		c.ResponseWriter.WriteHeader(500)
		retJSON.Status = "fail"
		jsonData, _ := json.Marshal(retJSON)

		c.ResponseWriter.Write(jsonData)
		panic("Need Check userInsert RowsAffected")
	}

	fmt.Println(dbobject)
}

func SignIn(c *router.Context) {
	loginobject := appcontext.LoginInfo{}
	var wallets appcontext.Wallets
	var wallet appcontext.Wallet

	recvData, _ := ioutil.ReadAll(c.Request.Body)

	loginobject.ID = c.Params["user_id"].(string)
	loginobject.Passwd = string(recvData)

	var userId string
	var userPasswd string

	query := `SELECT id FROM myEtherWallet.user WHERE id=` + `"` + loginobject.ID + `"` + `;`
	queryRows, queryErr := MySQL.SQLClient.Query(query)
	if queryErr != nil {
		fmt.Println("id query error")
	}
	for queryRows.Next() {
		err := queryRows.Scan(&userId)
		if err != nil {
			fmt.Println("id select error")
		}
	}
	defer queryRows.Close()

	fmt.Println("id: ", userId)

	if userId == "" {
		fmt.Println("non-Exist Email")
		c.ResponseWriter.WriteHeader(501)

		return
	}

	sqlQuery := "select passwd from myEtherWallet.user where id =  " + "'" + loginobject.ID + "'" + ";"
	rows, err := MySQL.SQLClient.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userPasswd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// all passwd should not be empty
	if (loginobject.Passwd == "" || userPasswd == "") || (loginobject.Passwd != userPasswd) {
		fmt.Println("Wrong Passworld")
		c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
	} else {
		// issue a access token(JWT)
		accessTokenExpirationTime := time.Now().Add(1 * time.Minute)
		accessClaims := &driver.Claims{
			Username: loginobject.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: accessTokenExpirationTime.Unix(),
			},
		}
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
		accessTokenString, err := accessToken.SignedString(driver.JwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			fmt.Println("Access token generation error")
			c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(c.ResponseWriter, &http.Cookie{
			Name:    "accessToken",
			Value:   accessTokenString,
			Expires: accessTokenExpirationTime,
		})

		// issue a refresh token(JWT)
		refreshTokenExpirationTime := time.Now().Add(1440 * time.Minute)
		refreshClaims := &driver.Claims{
			Username: loginobject.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: refreshTokenExpirationTime.Unix(),
			},
		}
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		refreshTokenString, err := refreshToken.SignedString(driver.JwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			fmt.Println("Refresh token generation error")
			log.Fatal(err)
			c.ResponseWriter.WriteHeader(http.StatusInternalServerError)

			return
		}

		http.SetCookie(c.ResponseWriter, &http.Cookie{
			Name:    "refreshToken",
			Value:   refreshTokenString,
			Expires: refreshTokenExpirationTime,
		})

		fmt.Println("accessTokenString: ", accessTokenString)
		fmt.Println("refreshTokenString: ", refreshTokenString)

		// Need querying all wallet infos
		sqlQuery := "select pubKey from myEtherWallet.wallet where ownerId =  " + "'" + loginobject.ID + "'" + ";"
		rows, err := MySQL.SQLClient.Query(sqlQuery)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&wallet.WalletAddress)
			if err != nil {
				log.Fatal(err)
			}
			wallets.List = append(wallets.List, wallet)
		}

		// make wallets json object
		jsonData, _ := json.Marshal(wallets)

		c.ResponseWriter.WriteHeader(200)
		c.ResponseWriter.Write(jsonData)
	}
}

func Refresh(c *router.Context) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	fmt.Println("Start refresh")

	cookie, err := c.Request.Cookie("refreshToken")
	if err != nil {
		if err == http.ErrNoCookie {
			c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fmt.Println("Cookie is fine")
	}

	tknStr := cookie.Value
	claims := &driver.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return driver.JwtKey, nil
	})
	if !tkn.Valid {
		c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		fmt.Println("Token is good")
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fmt.Println("Refresh Token is no problem")
	}

	fmt.Println("username: ", claims.Username)

	// Test time = 1 minute
	accessTokenExpirationTime := time.Now().Add(1 * time.Minute)
	accessClaims := &driver.Claims{
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(driver.JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		fmt.Println("Access token generation error")
		c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.ResponseWriter, &http.Cookie{
		Name:    "accessToken",
		Value:   accessTokenString,
		Expires: accessTokenExpirationTime,
	})

	// Test time = 1 minute
	refreshTokenExpirationTime := time.Now().Add(1440 * time.Minute)
	refreshClaims := &driver.Claims{
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(driver.JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		fmt.Println("Refresh token generation error")
		log.Fatal(err)
		c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.ResponseWriter, &http.Cookie{
		Name:    "refreshToken",
		Value:   refreshTokenString,
		Expires: refreshTokenExpirationTime,
	})

	fmt.Println("accessTokenString: ", accessTokenString)
	fmt.Println("refreshTokenString: ", refreshTokenString)

	c.ResponseWriter.WriteHeader(200)
}

func JWTauth(next router.HandlerFunc) router.HandlerFunc {
	return func(c *router.Context) {
		fmt.Println("Check JWT")

		cookie, err := c.Request.Cookie("accessToken")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("Cookie Error")
				c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
				return
			}
			fmt.Println("What value?")
			c.ResponseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := cookie.Value
		claims := &driver.Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return driver.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				fmt.Println("JWT Error")
				c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
				return
			}
			c.ResponseWriter.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			fmt.Println("Valid fail")
			c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println("Token in Valid")
		c.TempData = claims.Username

		// Token auth가 완료되면, HandlerFunc 호출
		next(c)
	}
}

func AddWallet(c *router.Context) {
	userid := c.Params["user_id"].(string)

	var wallet appcontext.Wallet
	var wallets appcontext.Wallets

	userPriv, userPub := account.GenKey()
	walletInsert, err := MySQL.SQLClient.Exec("INSERT `myEtherWallet`.`wallet` (`pubKey`, `privKey`, `ownerId`) VALUES (?, ?, ?);",
		userPub,
		userPriv,
		userid,
	)
	if err != nil {
		log.Fatal(err)
		panic("walletInsert ERROR")
	}

	n, err := walletInsert.RowsAffected()
	if n == 1 {
		fmt.Println("walletInsert success")
		// Need querying all wallet infos
		sqlQuery := "select pubKey from myEtherWallet.wallet where ownerId =  " + "'" + userid + "'" + ";"
		rows, err := MySQL.SQLClient.Query(sqlQuery)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&wallet.WalletAddress)
			if err != nil {
				log.Fatal(err)
			}
			wallets.List = append(wallets.List, wallet)
		}

		// make wallets json object
		jsonData, _ := json.Marshal(wallets)

		c.ResponseWriter.Write(jsonData)
		c.ResponseWriter.WriteHeader(201)
	} else {
		panic("Need Check walletInsert RowsAffected")
	}
}

func ReadContract(c *router.Context) {
	result := contract.Call(MyContract)
	c.ResponseWriter.Write([]byte(result))
}

func WriteContract(c *router.Context) {
	var req appcontext.ContractReq
	var ownerPrivKey string

	recvData, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(recvData, &req)

	sqlQuery := "select privKey from myEtherWallet.wallet where pubKey =  " + "'" + req.WalletAddress + "'" + ";"
	rows, err := MySQL.SQLClient.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&ownerPrivKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	transaction := contract.Transaction{}
	nonce := MyContract.GetNonce(ownerPrivKey)
	value, _ := strconv.ParseInt(req.Value, 10, 64)

	transaction.SetTransaction(uint(value), ownerPrivKey, nonce)

	txResult := contract.Send(transaction, MyContract)
	c.ResponseWriter.Write([]byte(txResult))
}

func GetWalletInfo(c *router.Context) {
	var metaData appcontext.WalletMetaData

	walletAddress := c.Params["walletAddress"].(string)

	// Make faster using go-rutine
	etherBalance, _ := strconv.ParseFloat(fmt.Sprint(eth.GetBalacnce(walletAddress)), 32)
	ercBalance, _ := strconv.ParseFloat(fmt.Sprint(erctoken.GetBalacnce(MyToken, walletAddress)), 32)
	ercHistroy := erctoken.GetRecentHistory(walletAddress)

	metaData.ETHbalance = fmt.Sprintf("%f", etherBalance)
	metaData.HRTbalance = fmt.Sprintf("%f", ercBalance)
	metaData.WalletHistroy = ercHistroy

	jsonData, _ := json.Marshal(metaData)

	c.ResponseWriter.Write(jsonData)
	c.ResponseWriter.WriteHeader(200)
}

func TransferToken(c *router.Context) {
	var transferReq appcontext.Transfer
	var ownerPrivKey string

	recvData, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(recvData, &transferReq)

	ercBalance, _ := strconv.ParseFloat(fmt.Sprint(erctoken.GetBalacnce(MyToken, transferReq.FromWallet)), 32)
	quantity, _ := strconv.ParseFloat(transferReq.Quantity, 32)

	if ercBalance > quantity {
		fmt.Println("Can send token!")

		sqlQuery := "select privKey from myEtherWallet.wallet where pubKey =  " + "'" + transferReq.FromWallet + "'" + ";"
		rows, err := MySQL.SQLClient.Query(sqlQuery)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&ownerPrivKey)
			if err != nil {
				log.Fatal(err)
			}
		}

		transcation := erctoken.Transaction{}
		transcation.SetTransaction(transferReq.ToWallet, quantity, ownerPrivKey)

		txResult := erctoken.TransferToken(transcation, MyToken)

		c.ResponseWriter.Write([]byte(txResult))
		c.ResponseWriter.WriteHeader(200)
	}
	fmt.Println("Not enough HRT")
	c.ResponseWriter.WriteHeader(500)
}
