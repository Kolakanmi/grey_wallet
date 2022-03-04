package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

//CORS - adds CORS policy to request
func CORS(h http.Handler) http.Handler {
	headersOK := handlers.AllowedHeaders([]string{"*"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	c := handlers.CORS(headersOK, originsOK, methodsOK)
	return c(h)
}
