package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sonereker/simple-auth/internal"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type userHandler struct {
	app *internal.Server
}

//NewHandler returns a new userHandler instance
func NewHandler(s *internal.Server) *userHandler {
	return &userHandler{s}
}

//GetCurrentUser returns currently authenticated user
func (uh *userHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("UserID")
	var user UserDBModel
	uh.app.DB.Take(&user, "id = ?", id)

	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(user.AsDTO())
	_, _ = w.Write(res)
}

//SignUp registers new user and logs-in newly registered user, and returns a JWT token
func (uh *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user UserDBModel
	result := uh.app.DB.Take(&user, "email = ?", authRequest.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

//LogIn logs in user with given credentials, and returns a JWT token
func (uh *userHandler) LogIn(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user UserDBModel
	result := uh.app.DB.Take(&user, "email = ?", authRequest.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !comparePasswords(user.Password, []byte(authRequest.Password)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &JWTClaims{
		Email:  authRequest.Email,
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshExpirationTime := time.Now().Add(24 * time.Hour)
	refreshClaims := &JWTClaims{
		Email:  authRequest.Email,
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshTokenSecret)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(AuthResponse{Token: tokenString, RefreshToken: refreshTokenString, User: user.AsDTO()})
}
