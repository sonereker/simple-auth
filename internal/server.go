package internal

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func NewServer(db *gorm.DB, r *mux.Router) *Server {
	return &Server{
		DB:     db,
		Router: r,
	}
}
