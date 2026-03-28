package jwt

import (
	"errors"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})
	privKey, err := os.ReadFile("./keys/id_ed25519")
	if err != nil {
		return "", err
	}
	tokenString, err := jwtToken.SignedString(privKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(token string) (string, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		pubKey, err := os.ReadFile("./keys/id_ed25519.pub")
		if err != nil {
			return nil, err
		}
		return jwt.ParseEdPublicKeyFromPEM(pubKey)
	})
	if err != nil {
		return "", err
	}
	if !tokenObj.Valid {
		return "", errors.New("Invalid token")
	}
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Failed to extract claims")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("Failed to extract user ID")
	}
	return userID, nil
}
