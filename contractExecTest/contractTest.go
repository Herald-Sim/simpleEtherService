package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"store"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://ropsten.infura.io/v3/857781c59662445bb934c92e10433b54")
	if err != nil {
		log.Fatal(err)
	}

	// Metamask
	privateKey, err := crypto.HexToECDSA("43f28e52e0b619195f73a86eee963840976ead52825150f34f531c34d1f2b995")
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	// My contract account address
	address := common.HexToAddress("0x8364fd2B18B15c0ABd7A86A042c532964FcDb8B1")

	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	// Send value
	
	newBigInt := big.NewInt(9999)
	tx, err := instance.Set(auth, newBigInt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())
	

	// Receive value
	/*
	rawResult, err := instance.Get(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	result := rawResult.String()
	fmt.Println("===Result===")
	fmt.Println(result)
	*/
	/*
		// load bifrost-setValue
		instance, err := bifrost.NewBifrost(address, client)
		if err != nil {
			log.Fatal(err)
		}

		// setValue parameter
		newBigInt := big.NewInt(500)

		tx, err := instance.SetValue(auth, newBigInt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

		rawResult, err :=
		if err != nil {
			log.Fatal(err)
		}
		result := rawResult.String()
		fmt.Println("===Result===")
		fmt.Println(result)

		//instance.BifrostCaller.GetValue(&bind.CallOpts{})
		//instance.GetValue(&bind.CallOpts{})
	*/
}
