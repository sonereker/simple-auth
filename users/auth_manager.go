package users

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type authManager struct {
	tokenSecret string
}

func NewAuthManager(tokenSecret string) *authManager {
	return &authManager{
		tokenSecret: tokenSecret,
	}
}

type UserClaims struct {
	jwt.StandardClaims
	UserID uint
	Email  string
}

func (am *authManager) verifyToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(am.tokenSecret), nil
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

func hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
