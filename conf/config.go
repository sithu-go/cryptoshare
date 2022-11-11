package conf

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var (
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string

	INFURA_BASE_URL string
	INFURA_API_KEY  string
)

func init() {
	// Load env file
	err := godotenv.Load("./conf/.env")
	if err != nil {
		log.Println("error opening .env file")
		log.Fatalf(err.Error(), "FGDD")
		return
	}

	// Load rsa [private]
	privateBytes, err := os.ReadFile(os.Getenv("RSA_PRIVATE"))
	if err != nil {
		log.Println("Error on loading private key: ", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Println("Error on parsing private key: ", err)
	}
	PrivateKey = privateKey

	// Load rsa [public]
	publicBytes, err := os.ReadFile(os.Getenv("RSA_PUBLIC"))
	if err != nil {
		log.Println("Error on loading public key: ", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Println("Error on parsing key: ", err)
	}
	PublicKey = publicKey

	// Load rsa [secret]
	RefreshSecret = os.Getenv("RSA_SECRET")

	path := "storage/qrcode"
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "no such file or directory") {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}

		}
	}

}
