package users

import "net/http"

func (uh *userHandler) RegisterRoutes() {
	uh.app.Router.HandleFunc("/users", uh.SignUp).Methods(http.MethodOptions, http.MethodPost)
	uh.app.Router.HandleFunc("/users/login", uh.SignUp).Methods(http.MethodOptions, http.MethodPost)
	uh.app.Router.Handle("/users/current", AuthenticateToken(http.HandlerFunc(uh.GetCurrentUser))).Methods(http.MethodOptions, http.MethodGet)
}
