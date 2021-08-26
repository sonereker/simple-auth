package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//AuthManager is the manager struct
type AuthManager struct {
	TokenSecret string
}

//NewAuthManager creates a new AuthManager
func NewAuthManager(tokenSecret string) *AuthManager {
	return &AuthManager{
		TokenSecret: tokenSecret,
	}
}

//UserClaims JWT custom claims struct
type UserClaims struct {
	jwt.StandardClaims
	ID uint
}

//GenerateToken generates a new JWT token
func (manager *AuthManager) GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &UserClaims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.TokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//VerifyToken verifies given token
func (manager *AuthManager) VerifyToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.TokenSecret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if token.Valid {
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}
		return claims, err
	}

	return nil, err
}

//Hash generates hash for the given password
func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

//IsCorrectPassword compares given hashed password with the plain password
func IsCorrectPassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
