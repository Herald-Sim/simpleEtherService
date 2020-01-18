package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"store"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Transaction struct {
	SenderKey string
	// main object to excute transaction
	Auth  *bind.TransactOpts
	Value *big.Int
	//CallOpts *bind.CallOpts
}

type MyContract struct {
	EthClient        *ethclient.Client
	ContractAddress  string
	ContractInstance *store.Store
}

func (myContract *MyContract) InitContract(myContractAddress string) {
	var err error

	myContract.EthClient, err = ethclient.Dial("https://ropsten.infura.io/v3/857781c59662445bb934c92e10433b54")
	if err != nil {
		fmt.Println("ethclient Dial Error")
		log.Fatal(err)
	}

	myContract.ContractAddress = myContractAddress
	contractAddress := common.HexToAddress(myContract.ContractAddress)

	myContract.ContractInstance, err = store.NewStore(contractAddress, myContract.EthClient)
	if err != nil {
		fmt.Println("ContractInstance generation error")
		log.Fatal(err)
	}
}

func (myContract *MyContract) GetNonce(ownerPrivKey string) uint64 {
	privKey, err := crypto.HexToECDSA(ownerPrivKey)
	if err != nil {
		fmt.Println("userKey converting error")
	}
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := myContract.EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	return nonce
}

func (myTransaction *Transaction) SetTransaction(value uint, ownerPrivKey string, nonce uint64) {
	// Set transaction subject
	myTransaction.SenderKey = ownerPrivKey
	privKey, err := crypto.HexToECDSA(myTransaction.SenderKey)
	if err != nil {
		fmt.Println("userKey converting error")
	}

	myTransaction.Auth = bind.NewKeyedTransactor(privKey)
	myTransaction.Auth.GasLimit = uint64(350000)
	//myTransaction.Auth.GasPrice = big.NewInt(40000000000)
	myTransaction.Auth.Value = big.NewInt(0)
	myTransaction.Auth.Nonce = big.NewInt(int64(nonce))

	myTransaction.Value = big.NewInt(int64(value))
}

func Call(myContract MyContract) string {
	rawResult, err := myContract.ContractInstance.Get(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	result := rawResult.String()
	return result
}

func Send(transaction Transaction, myContract MyContract) string {
	gasPrice, err := myContract.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	transaction.Auth.GasPrice = gasPrice

	tx, err := myContract.ContractInstance.Set(transaction.Auth, transaction.Value)
	if err != nil {
		log.Fatal(err)
	}
	result := fmt.Sprintf("%s", tx.Hash().Hex())
	return result
}
