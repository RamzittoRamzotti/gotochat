package jwt

import (
	"errors"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(name string) (string, error) {
	privKeyBytes, err := os.ReadFile("./keys/id_ed25519")
	if err != nil {
		return "", err
	}
	privKey, err := jwt.ParseEdPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"name": name,
	})
	return token.SignedString(privKey)
}

func ValidateToken(tokenStr string) (string, error) {
	pubKeyBytes, err := os.ReadFile("./keys/id_ed25519.pub")
	if err != nil {
		return "", err
	}
	pubKey, err := jwt.ParseEdPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return "", err
	}
	tokenObj, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return pubKey, nil
	})
	if err != nil || !tokenObj.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to extract claims")
	}
	name, ok := claims["name"].(string)
	if !ok {
		return "", errors.New("failed to extract name")
	}
	return name, nil
}
