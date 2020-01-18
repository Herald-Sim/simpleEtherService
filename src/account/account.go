package account

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func GenKey() (string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	//fmt.Println("PrivKey: ", hexutil.Encode(privateKeyBytes)[2:])
	privateKeyStr := fmt.Sprintf("%s", hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	//fmt.Println("PubKey: ", hexutil.Encode(hash.Sum(nil)[12:]))
	pubKeyStr := fmt.Sprintf("%s", hexutil.Encode(hash.Sum(nil)[12:]))

	return privateKeyStr, pubKeyStr
}
