package users

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

var (
	tokenSecret        = []byte("2AB89F28-0DF2-4D47-93AD-97810483C515")
	refreshTokenSecret = []byte("53d6cfb8-4c78-11eb-ae93-0242ac130002")
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID uint
	Email  string
}

func AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := extractJWTToken(r)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusExpectationFailed)
		}

		claims, err := parseJWT(jwtToken, string(tokenSecret))

		if err != nil {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
		} else {
			ctx := context.WithValue(r.Context(), "UserID", claims.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	})
}

func parseJWT(tokenString string, key string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		}

		return nil, err
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, err
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, err
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func extractJWTToken(req *http.Request) (string, error) {
	tokenString := req.Header.Get("Authorization")

	if tokenString == "" {
		return "", fmt.Errorf("could not find token")
	}

	tokenString, err := stripTokenPrefix(tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func stripTokenPrefix(tok string) (string, error) {
	tokenParts := strings.Split(tok, " ")

	if len(tokenParts) < 2 {
		return tokenParts[0], nil
	}

	return tokenParts[1], nil
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
