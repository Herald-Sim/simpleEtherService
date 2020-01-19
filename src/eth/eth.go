package eth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

type Ethereum struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func calculateKeccak256(addr []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(addr)
	return hash.Sum(nil)
}

func checksumByte(addr byte, hash byte) string {
	result := strconv.FormatUint(uint64(addr), 16)
	if hash >= 8 {
		return strings.ToUpper(result)
	} else {
		return result
	}
}

func ChecksumAddr(addrStr string) (string, error) {
	addrStr = addrStr[2:]
	addr, err := hex.DecodeString(addrStr)
	if err != nil {
		return "", err
	}
	hash := calculateKeccak256([]byte(strings.ToLower(addrStr)))

	result := "0x"

	// fmt.Println("addr:", hex.EncodeToString(addr), addr)
	// fmt.Println("hash:", hex.EncodeToString(hash), hash)
	for i, b := range addr {
		result += checksumByte(b>>4, hash[i]>>4)
		result += checksumByte(b&0xF, hash[i]&0xF)
	}

	return result, nil
}

func GetBalacnce(walletAddress string) *big.Float {
	ether := Ethereum{}

	realUnit := new(big.Int)
	realResult := big.NewFloat(0)

	convert := big.NewFloat(0.000000000000000001)

	resp, err := http.Get("https://api-ropsten.etherscan.io/api?module=account&action=balance&address=" + walletAddress + "&tag=latest&apikey=PJC7PSH73MEI5J5AD8URMAVFB1G2FMGHYJ")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &ether)
	fmt.Println("ether: ", ether)

	realUnit, _ = realUnit.SetString(ether.Result, 10)

	floatResult := new(big.Float).SetInt(realUnit)
	realResult.Mul(floatResult, convert)

	return realResult
}
