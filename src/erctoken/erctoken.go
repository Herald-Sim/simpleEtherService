// EhtherScanAPI: PJC7PSH73MEI5J5AD8URMAVFB1G2FMGHYJ
/*ToDo
1. GetBalance 단위 수정
2. 최근 거래내역 추적 추가
*/

package erctoken

import (
	"appcontext"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
	"token"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

type ERC20Token struct {
	EthClient     *ethclient.Client
	TokenAddress  string
	TokenInstance *token.Token
}

func (myERC *ERC20Token) InitToken(myTokenAddress string) {
	var err error

	myERC.EthClient, err = ethclient.Dial("https://ropsten.infura.io/v3/857781c59662445bb934c92e10433b54")
	if err != nil {
		fmt.Println("ethclient Dial Error")
		log.Fatal(err)
	}

	myERC.TokenAddress = myTokenAddress
	tokenAddress := common.HexToAddress(myERC.TokenAddress)

	myERC.TokenInstance, err = token.NewToken(tokenAddress, myERC.EthClient)
	if err != nil {
		fmt.Println("TokenInstance generation error")
		log.Fatal(err)
	}
}

type Transaction struct {
	SenderKey string
	// main object to excute transaction
	Auth      bind.TransactOpts
	ToAddress common.Address
	Amount    *big.Int
}

type Trans struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Result  []TransactionRecord `json:"result"`
}

type TransactionRecord struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDeciaml      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	Confirmations     string `json:"confirmations"`
}

func DecimalsToWei(amount decimal.Decimal, decimals int) *big.Int {
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

func (myTransaction *Transaction) SetTransaction(toAddress string, amount float64, sender string) {
	// Set amout of token to send
	myTransaction.Amount = DecimalsToWei(decimal.NewFromFloat(amount), int(18))

	// Set transaction subject
	myTransaction.SenderKey = sender
	key, err := crypto.HexToECDSA(myTransaction.SenderKey)
	if err != nil {
		fmt.Println("userKey converting error")
	}
	myTransaction.Auth = *bind.NewKeyedTransactor(key)
	myTransaction.Auth.GasLimit = uint64(350000)
	myTransaction.Auth.GasPrice = big.NewInt(40000000000)
	myTransaction.Auth.Value = big.NewInt(0)
	myTransaction.Auth.Nonce = nil

	myTransaction.ToAddress = common.HexToAddress(toAddress)
}

func TransferToken(transaction Transaction, myToken ERC20Token) string {
	tx, err := myToken.TokenInstance.Transfer(&transaction.Auth, transaction.ToAddress, transaction.Amount)
	var result string

	if err != nil {
		//log.Fatal(err)
		panic(err)

	} else if err == nil {
		log.Printf("tx: %s", tx.Hash().Hex())
		result = fmt.Sprintf("%s", tx.Hash().Hex())
	}

	return result
}

func GetBalacnce(myToken ERC20Token, walletAddress string) *big.Float {
	address := common.HexToAddress(walletAddress)
	opts := bind.CallOpts{}

	realResult := big.NewFloat(0)
	convert := big.NewFloat(0.000000000000000001)

	balance, err := myToken.TokenInstance.TokenCaller.BalanceOf(&opts,
		address)

	defer func() {
		s := recover()
		fmt.Println(s)
	}()

	if err != nil {
		fmt.Println("GetBalance error")
		panic(err)
	}

	floatResult := new(big.Float).SetInt(balance)
	realResult.Mul(floatResult, convert)

	return realResult
}

func GetRecentHistory(walletAddress string) appcontext.History {
	trans := Trans{}

	history := appcontext.History{}
	recent := appcontext.RecentTransaction{}

	realUnit := new(big.Int)
	realResult := big.NewFloat(0)

	convert := big.NewFloat(0.000000000000000001)

	resp, err := http.Get(`https://api-ropsten.etherscan.io/api?module=account&action=tokentx&contractaddress=0x1e3f9caf34340f8313ed416d71959aa209c2114f&address=` + walletAddress + `&page=1&offset=5&sort=desc&apikey=PJC7PSH73MEI5J5AD8URMAVFB1G2FMGHYJ`)

	if err != nil {
		fmt.Println(nil)
	}

	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(data, &trans)
	fmt.Println(trans)

	for _, tran := range trans.Result {
		i, err := strconv.ParseInt(tran.TimeStamp, 10, 64)
		if err != nil {
			panic(err)
		}

		recent.TimeStamp = fmt.Sprint(time.Unix(i, 0))
		recent.From = tran.From
		recent.To = tran.To

		realUnit, _ = realUnit.SetString(tran.Value, 10)
		floatResult := new(big.Float).SetInt(realUnit)
		realResult.Mul(floatResult, convert)

		recent.Value = fmt.Sprint(realResult)
		recent.TokenName = tran.TokenName

		history.List = append(history.List, recent)
	}

	return history
}