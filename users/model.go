package users

import (
	"github.com/sonereker/simple-auth/grpc/v1"
	"gorm.io/gorm"
)

//UserDBModel is the DB model for User
type UserDBModel struct {
	gorm.Model
	Email        string `gorm:"index not null"`
	Password     string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

//AsResponse returns a simplified DTO object for the given UserDBModel
func (u *UserDBModel) AsResponse() *grpc.UserResponse {
	return &grpc.UserResponse{
		Email: u.Email,
	}
}
