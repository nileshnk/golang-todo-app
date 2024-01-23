package auth_controller

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	Types "github.com/nileshnk/golang-todo-app/types"
)

func CreateAccessToken(tokenPayload Types.TokenPayload) (string, error) {
	createClaim := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userId": tokenPayload.UserId,
	})

	pemString := os.Getenv("JWT_PRIVATE_KEY")

	if pemString == "" {
		fmt.Println("JWT_PRIVATE_KEY environment variable not set")
		return "", fmt.Errorf("JWT_PRIVATE_KEY environment variable not set")
	}

	block, rest := pem.Decode([]byte(pemString))

	if block == nil {
		fmt.Println("Failed to decode PEM block")
		fmt.Println("Remaining data after decoding:", string(rest))
		return "", fmt.Errorf("Failed to decode PEM block")
	}
	parseResult, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		fmt.Println("Failed to parse private key:", err)
		return "", err
	}

	privateKey := parseResult

	createClaim.Valid = true

	tokenString, tokenSignErr := createClaim.SignedString(privateKey)
	if tokenSignErr != nil {
		fmt.Println("Error signing token")
		fmt.Println(tokenSignErr)
		return "", tokenSignErr
	}

	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (bool, jwt.MapClaims, error) {
	pemString := os.Getenv("JWT_PUBLIC_KEY")

	if pemString == "" {
		fmt.Println("JWT_PUBLIC_KEY environment variable not set")
		return false, nil, fmt.Errorf("JWT_PUBLIC_KEY environment variable not set")
	}

	block, rest := pem.Decode([]byte(pemString))

	if block == nil {
		fmt.Println("Failed to decode PEM block")
		fmt.Println("Remaining data after decoding:", string(rest))
		return false, nil, fmt.Errorf("Failed to decode PEM block")
	}

	parseResult, parseErr := x509.ParsePKIXPublicKey(block.Bytes)

	if parseErr != nil {
		fmt.Println("Failed to parse public key:", parseErr)
		return false, nil, parseErr
	}

	publicKey := parseResult.(*rsa.PublicKey)

	token, tokenParseErr := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if tokenParseErr != nil {
		fmt.Println("Error parsing token:", tokenParseErr)
		return false, nil, tokenParseErr

	}
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)

	return true, claims, nil
	// token, tokenParseErr := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return rsa.PublicKey(), nil
	// })
}
