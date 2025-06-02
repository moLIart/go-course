package services

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTService_Sign_ReturnsValidToken(t *testing.T) {
	secret := "testsecret"
	service := NewJWTService(secret)
	name := "testuser"

	tokenString, err := service.Sign(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokenString == "" {
		t.Fatal("expected a token string, got empty string")
	}

	// Token should have three parts separated by '.'
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		t.Errorf("expected token to have 3 parts, got %d", len(parts))
	}
}

func TestJWTService_Sign_TokenContainsCorrectClaims(t *testing.T) {
	secret := "anothersecret"
	service := NewJWTService(secret)
	name := "alice"

	tokenString, err := service.Sign(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected jwt.MapClaims")
	}

	if claims["name"] != name {
		t.Errorf("expected name claim to be %q, got %q", name, claims["name"])
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("expected exp claim to be float64")
	}
	expTime := time.Unix(int64(exp), 0)
	if expTime.Before(time.Now().Add(1*time.Hour)) || expTime.After(time.Now().Add(2*time.Hour+time.Minute)) {
		t.Errorf("expected exp to be about 2 hours from now, got %v", expTime)
	}
}
func TestJWTService_VerifyToken_ValidToken(t *testing.T) {
	secret := "mysecret"
	service := NewJWTService(secret)
	name := "bob"

	tokenString, err := service.Sign(name)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	token, err := service.Verify(tokenString)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !token.Valid {
		t.Error("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected jwt.MapClaims")
	}
	if claims["name"] != name {
		t.Errorf("expected name claim %q, got %q", name, claims["name"])
	}
}

func TestJWTService_VerifyToken_InvalidSignature(t *testing.T) {
	secret := "secret1"
	service := NewJWTService(secret)
	otherSecret := "secret2"
	name := "eve"

	// Sign with a different secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"exp":  time.Now().Add(2 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(otherSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = service.Verify(tokenString)
	if err == nil {
		t.Error("expected error due to invalid signature, got nil")
	}
}

func TestJWTService_VerifyToken_InvalidTokenFormat(t *testing.T) {
	secret := "secret"
	service := NewJWTService(secret)

	invalidToken := "not.a.valid.token"
	_, err := service.Verify(invalidToken)
	if err == nil {
		t.Error("expected error for invalid token format, got nil")
	}
}
