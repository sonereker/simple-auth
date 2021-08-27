package users

import (
	"context"
	"github.com/sonereker/simple-auth/internal/pb/v1"
	"testing"
)

func TestRegisterWOEmail(t *testing.T) {
	service := NewUserService(nil, nil)
	request := &pb.RegistrationRequest{
		Email:    "",
		Password: "hello123",
	}
	_, err := service.Register(context.Background(), request)
	requestValidationError, ok := err.(pb.RegistrationRequestValidationError)
	if err != nil && (!ok || requestValidationError.Reason() != "value must be a valid email address") {
		t.Errorf("User should not be able to register with an empty email")
	}
}

func TestRegisterWOValidEmail(t *testing.T) {
	service := NewUserService(nil, nil)
	request := &pb.RegistrationRequest{
		Email:    "somedummytext",
		Password: "hello123",
	}
	_, err := service.Register(context.Background(), request)
	requestValidationError, ok := err.(pb.RegistrationRequestValidationError)
	if err != nil && (!ok || requestValidationError.Reason() != "value must be a valid email address") {
		t.Errorf("User should not be able to register with an invalid email")
	}
}

func TestRegisterWShortPassword(t *testing.T) {
	service := NewUserService(nil, nil)
	request := &pb.RegistrationRequest{
		Email:    "dummy@email.com",
		Password: "123",
	}
	_, err := service.Register(context.Background(), request)
	requestValidationError, ok := err.(pb.RegistrationRequestValidationError)
	if err != nil && (!ok || requestValidationError.Reason() != "value length must be at least 6 runes") {
		t.Errorf("User should not be able to register with a short password")
	}
}

func TestLoginWOEmail(t *testing.T) {
	service := NewUserService(nil, nil)
	request := &pb.LoginRequest{
		Email:    "",
		Password: "hello123",
	}
	_, err := service.Login(context.Background(), request)
	requestValidationError, ok := err.(pb.LoginRequestValidationError)
	if err != nil && (!ok || requestValidationError.Reason() != "value must be a valid email address") {
		t.Errorf("User should not be able to login with an empty email")
	}
}

func TestLoginWOValidEmail(t *testing.T) {
	service := NewUserService(nil, nil)
	request := &pb.LoginRequest{
		Email:    "somedummytext",
		Password: "hello123",
	}
	_, err := service.Login(context.Background(), request)
	requestValidationError, ok := err.(pb.LoginRequestValidationError)
	if err != nil && (!ok || requestValidationError.Reason() != "value must be a valid email address") {
		t.Errorf("User should not be able to login with an invalid email")
	}
}

func TestLoginWShortPassword(t *testing.T) {
	service := NewUserService(nil, nil)
	request := &pb.LoginRequest{
		Email:    "dummy@email.com",
		Password: "123",
	}
	_, err := service.Login(context.Background(), request)
	requestValidationError, ok := err.(pb.LoginRequestValidationError)
	if err != nil && (!ok || requestValidationError.Reason() != "value length must be at least 6 runes") {
		t.Errorf("User should not be able to login with a short password")
	}
}
