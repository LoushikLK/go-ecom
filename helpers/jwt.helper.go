package helpers

import (
	"ecommerce/configs"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	Exp    int64  `json:"exp"`
}

func GenerateToken(claims TokenClaims) (string, error) {

	configs := configs.AppEnv()

	secret := []byte(configs.JWT_SECRET_KEY)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": claims.UserId,
		"mobile": claims.Mobile,
		"email":  claims.Email,
		"exp":    claims.Exp,
	}).SignedString(secret)

	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseToken(jwtToken string) (*TokenClaims, error) {

	configs := configs.AppEnv()

	secret := []byte(configs.JWT_SECRET_KEY)

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &TokenClaims{
			UserId: claims["userId"].(string),
			Email:  claims["email"].(string),
			Mobile: claims["mobile"].(string),
			Exp:    int64(claims["exp"].(float64)),
		}, nil
	}

	return nil, errors.New("invalid token")

}
