package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
}

func NewJWTService(secretKey string) *JWTService {
	if secretKey == "" {
		panic("JWT secret key must be provided")
	}

	return &JWTService{
		secretKey: secretKey,
	}
}

func (s *JWTService) Sign(name string) (string, error) {
	claims := jwt.MapClaims{
		"name": name,
		"exp":  time.Now().Add(2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) Verify(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется ожидаемый метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secretKey), nil
	})
}
