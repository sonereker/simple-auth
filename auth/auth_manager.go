package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthManager struct {
	TokenSecret string
}

func NewAuthManager(tokenSecret string) *AuthManager {
	return &AuthManager{
		TokenSecret: tokenSecret,
	}
}

type UserClaims struct {
	jwt.StandardClaims
	UserID uint
	Email  string
}

func (am *AuthManager) verifyToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(am.TokenSecret), nil
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
	} else {
		return nil, err
	}
}

func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
