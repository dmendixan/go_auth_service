package services

import (
	"auth-service/config"
	"auth-service/internal/models"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte(config.JWTSecret)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken создает JWT access token с указанным сроком жизни (например, 1 час)
func GenerateAccessToken(userID uint, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	fmt.Printf("%+v\n", jwtKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("%+v\n", token)

	signedToken, err := token.SignedString([]byte(config.JWTSecret))

	fmt.Printf("%+v\n", err)
	return signedToken, err
}

func GenerateRefreshToken(userID uint) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	refreshToken := hex.EncodeToString(tokenBytes)

	rt := models.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := config.DB.Create(&rt).Error; err != nil {
		return "", err
	}
	return refreshToken, nil
}

func ValidateRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := config.DB.Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	if rt.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return &rt, nil
}
