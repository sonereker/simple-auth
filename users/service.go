package users

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sonereker/simple-auth/grpc/v1"
	"gorm.io/gorm"
	"time"
)

type userService struct {
	grpc.UnimplementedUsersServer
	authManager *authManager
	DB          *gorm.DB
}

func NewUserService(db *gorm.DB, am *authManager) *userService {
	return &userService{DB: db, authManager: am}
}

func (s *userService) Register(ctx context.Context, rr *grpc.RegistrationRequest) (*grpc.AuthenticationResponse, error) {
	var user UserDBModel
	result := s.DB.Take(&user, "email = ?", rr.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user.Password = hashAndSalt(user.Password)
		result := s.DB.Create(&user)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	authenticationResponse, err := s.Login(ctx, &grpc.LoginRequest{
		Email:    rr.Email,
		Password: rr.Password,
	})
	if err != nil {
		return nil, err
	}

	return authenticationResponse, nil
}

func (s *userService) Login(ctx context.Context, lr *grpc.LoginRequest) (*grpc.AuthenticationResponse, error) {
	var user UserDBModel
	result := s.DB.Take(&user, "email = ?", lr.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	if !comparePasswords(user.Password, []byte(lr.Password)) {
		return nil, errors.New("email or password is incorrect")
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &UserClaims{
		Email:  lr.Email,
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.authManager.tokenSecret)

	refreshExpirationTime := time.Now().Add(24 * time.Hour)
	refreshClaims := &UserClaims{
		Email:  lr.Email,
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.authManager.tokenSecret)
	if err != nil {
		return nil, err
	}

	return &grpc.AuthenticationResponse{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
		User:         user.AsResponse(),
	}, nil
}

func (s *userService) GetCurrent(ctx context.Context, in *grpc.EmptyParams) (*grpc.UserResponse, error) {
	var user UserDBModel
	s.DB.Take(&user, "id = ?", 1)

	return &grpc.UserResponse{Email: user.Email}, nil
}
