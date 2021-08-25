package internal

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

//InitRouter returns a new mux router instance with middlewares configured
func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(applicationJsonHeader)
	r.Use(mux.CORSMethodMiddleware(r))

	headers, origins, methods := configureCORSHandlers()
	r.Use(handlers.CORS(headers, origins, methods))
	return r
}

func applicationJsonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func configureCORSHandlers() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodOptions})
	return headers, origins, methods
}
