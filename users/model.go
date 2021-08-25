package users

import "gorm.io/gorm"

//UserDBModel is the DB model for User
type UserDBModel struct {
	gorm.Model
	Email        string `gorm:"index not null"`
	Password     string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
}

//AuthRequest is the DTO object for auth requests
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//UserDTO is the DTO object for user
type UserDTO struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

//AsDTO returns a simplified DTO object for the given UserDBModel
func (u *UserDBModel) AsDTO() *UserDTO {
	return &UserDTO{
		Email: u.Email,
		ID:    u.ID,
	}
}

//AuthResponse is the DTO object having generated tokens and user info
type AuthResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken"`
	User         *UserDTO `json:"account"`
}
