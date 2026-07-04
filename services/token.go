package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var SecretKey = []byte(os.Getenv("TOKEN_SECRET"))

func GenerateAccessToken(userID int64, email string) (string, error) {

	validityMins, err := strconv.Atoi(os.Getenv("TOKEN_VALIDITY_MINS"))
	if err != nil {
		validityMins = 15
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Duration(validityMins) * time.Minute)),
		"iat":     jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SecretKey)
}

func GenerateRefreshToken() string {
	refreshToken := uuid.NewString()
	if refreshToken == "" {
		refreshToken = "9b827290-321f-429d-b072-288721423f77"
	}

	return refreshToken
}

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid claims.")
	}

	return claims, nil
}
