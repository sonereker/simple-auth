package users

import (
	"github.com/sonereker/simple-auth/pb/v1"
	"gorm.io/gorm"
)

//UserDBModel is the DB model for User
type UserDBModel struct {
	gorm.Model
	Email    string `gorm:"index:idx_email,unique"`
	Password string `gorm:"not null"`
}

//AsResponse returns a simplified DTO object for the given UserDBModel
func (u *UserDBModel) AsResponse() *pb.UserResponse {
	return &pb.UserResponse{
		Email: u.Email,
	}
}
