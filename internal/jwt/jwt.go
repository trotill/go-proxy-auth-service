package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-proxy-auth-service/internal/env"
	"go-proxy-auth-service/internal/utils"
)

type CustomClaims struct {
	Login     string `json:"login"`
	Role      string `json:"role"`
	SessionId string `json:"sessionId"`
	Type      string `json:"type"`
	jwt.RegisteredClaims
}

var jwtSecret *rsa.PublicKey

func parseRSAPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS1 encoded public key: %w", err)
	}

	return rsaPub, nil
}

func GetJwtSecret() {
	config := env.GetEnv()
	data, err := utils.ReadFile(config.PublicKeyPath)
	if err != nil {
		panic("Could not get public key")
	}
	secret, err := parseRSAPublicKey(data)
	if err != nil {
		panic("Could not parse RSA AES public key")
	}
	jwtSecret = secret
}
func VerifyToken(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
