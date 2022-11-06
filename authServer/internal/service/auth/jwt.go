package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

func GenerateJwt(email string) (string, error) {
	privateKey, err := getRSAPrivateKey()
	if err != nil {
		return "", err
	}
	claims := make(jwt.MapClaims)
	claims["email"] = email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(72 * time.Hour).Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate a new jwt token. Error: %s", err)
	}
	return token, nil
}

func getRSAPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile("internal/config/keys/id_rsa")
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes of RSA private key file. Error: %s", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key from PEM. Error: %s", err)
	}
	return privateKey, nil
}
