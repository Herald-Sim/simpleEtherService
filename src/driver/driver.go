/* ToDo
1. CryptionRead interface 적용
*/

package driver

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

// SQLinfo : MySql DB object
type SQLinfo struct {
	SQLClient *sql.DB
}

// InitSQL : Initialize SQL driver
func (mySQL *SQLinfo) InitSQL() {
	client, err := sql.Open("mysql", "root:titantech0070@tcp(neocool.iptime.org:3306)/myEtherWallet")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("DB Connection Succeed")
		mySQL.SQLClient = client
	}
	var version string

	mySQL.SQLClient.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)
}

// JwtKey : JSON Web Token key
var JwtKey = []byte("herald-0070")

// Credentials : Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims : struct that will be encoded to a JWT.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// AESObject : For AES256 Enc and Dec
type AESObject struct {
	CipherBlock cipher.Block
}

// InitAES : To init cipherBlock
func (myAES *AESObject) InitAES() {
	cipherKey := "000000000heraldsim000000000-0070"
	var err error

	myAES.CipherBlock, err = aes.NewCipher([]byte(cipherKey))
	if err != nil {
		fmt.Println("NewCipher() Error")
		log.Fatal(err)
	}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]

	return encrypt[:len(encrypt)-int(padding)]
}

func (myAES *AESObject) EncryptCBC(plaintext []byte) []byte {
	fmt.Println("aes BlockSize: ", aes.BlockSize)
	plaintext = PKCS5Padding(plaintext, aes.BlockSize)
	fmt.Printf("[plain text to Enc]: %x\n", plaintext)

	// 초기화 벡터 공간(aes.BlockSize)만큼 더 생성
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	// 부분 슬라이스로 초기화 벡터 공간을 가져옴
	iv := ciphertext[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(myAES.CipherBlock, iv)

	// 암호화 블록 모드 인스턴스로 암호화
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	fmt.Printf("[cipher data]: %x\n", ciphertext)
	return ciphertext
}

func (myAES *AESObject) DecryptCBC(ciphertext []byte) []byte {
	// 블록 크기의 배수가 아니면 리턴
	if len(ciphertext)%aes.BlockSize != 0 {
		fmt.Println("암호화된 데이터의 길이는 블록 크기의 배수가 되어야합니다.")
		return nil
	}

	fmt.Printf("cipher: %x\n", ciphertext)

	iv := ciphertext[:aes.BlockSize]        // 부분 슬라이스로 초기화 벡터 공간을 가져옴
	ciphertext = ciphertext[aes.BlockSize:] // 부분 슬라이스로 암호화된 데이터를 가져옴

	fmt.Printf("iv: %x\n", iv)
	fmt.Printf("realCipher: %x\n", ciphertext)

	plaintext := make([]byte, len(ciphertext))            // 평문 데이터를 저장할 공간 생성
	mode := cipher.NewCBCDecrypter(myAES.CipherBlock, iv) // 암호화 블록과 초기화 벡터를 넣어서
	// 복호화 블록 모드 인스턴스 생성
	mode.CryptBlocks(plaintext, ciphertext) // 복호화 블록 모드 인스턴스로 복호화
	fmt.Println("Befor Trimming: ", plaintext)

	return PKCS5Trimming(plaintext)
}

func CryptionRead(myAES AESObject, data []byte) []byte {
	jsonString := string(data)

	encData := jsonString[1 : len(jsonString)-1]
	encByte, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	decByte := myAES.DecryptCBC(encByte)

	return decByte
}
