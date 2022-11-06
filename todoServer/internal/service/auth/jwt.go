package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

func ValidateToken(tokenStr string, key *rsa.PublicKey) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("wrong singning method, expected RSA but got: %s", jwtToken.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to validate token. Error: %s", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token is not valid")
}

func GetRSAPublicKey() (*rsa.PublicKey, error) {
	publicKeyBytes, err := os.ReadFile("internal/config/keys/id_rsa.pub")
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes of RSA public key file. Error: %s", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key from PEM. Error: %s", err)
	}
	return publicKey, nil
}
