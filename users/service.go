package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/sonereker/simple-auth/auth"
	"github.com/sonereker/simple-auth/pb/v1"
	"gorm.io/gorm"
)

type userService struct {
	pb.UnimplementedUserServiceServer
	authManager *auth.AuthManager
	DB          *gorm.DB
}

//NewUserService creates a new userService with provided params
func NewUserService(db *gorm.DB, am *auth.AuthManager) *userService {
	return &userService{DB: db, authManager: am}
}

//Register creates the new user and returns a token with created user info
func (service *userService) Register(ctx context.Context, rr *pb.RegistrationRequest) (*pb.AuthenticationResponse, error) {
	var user UserDBModel
	result := service.DB.Take(&user, "email = ?", rr.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user.Email = rr.Email
		hashedPassword, err := auth.Hash(rr.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
		fmt.Printf("Creating new user with email %s and password %s\n", rr.Email, hashedPassword)
		result := service.DB.Create(&user)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	authenticationResponse, err := service.Login(ctx, &pb.LoginRequest{
		Email:    rr.Email,
		Password: rr.Password,
	})
	if err != nil {
		return nil, err
	}

	return authenticationResponse, nil
}

//Login returns a JWT token if a user exists with given credentials
func (service *userService) Login(_ context.Context, lr *pb.LoginRequest) (*pb.AuthenticationResponse, error) {
	var user UserDBModel
	result := service.DB.Take(&user, "email = ?", lr.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user with email " + lr.Email + " not found")
	}

	if !auth.IsCorrectPassword(user.Password, lr.Password) {
		return nil, errors.New("password is incorrect")
	}

	tokenString, err := service.authManager.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &pb.AuthenticationResponse{
		Token: tokenString,
		User:  user.AsResponse(),
	}, nil
}

//GetCurrent returns current user with the token
func (service *userService) GetCurrent(ctx context.Context, _ *pb.Empty) (*pb.UserResponse, error) {
	id := ctx.Value("id")
	var user UserDBModel
	service.DB.Take(&user, "id = ?", id)

	return &pb.UserResponse{Email: user.Email}, nil
}