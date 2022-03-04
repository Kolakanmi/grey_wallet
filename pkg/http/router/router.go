package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type (
	//Middleware - type alias to avoid casting
	Middleware = func(handler http.Handler) http.Handler

	//Route struct
	Route struct {
		Path        string
		Method      string
		Queries     []string
		Handler     http.Handler
		Middlewares []Middleware
	}

	//Config struct
	Config struct {
		// StaticPaths map[string]http.Handler

		GlobalMiddlewares []Middleware
		Routes            []Route
		NotFoundHandler   http.Handler
	}
)

//New - returns a new handler with all registered routes
func New(conf *Config) (http.Handler, error) {
	r := mux.NewRouter()
	for _, middleware := range conf.GlobalMiddlewares {
		r.Use(middleware)
	}
	for _, route := range conf.Routes {
		h := route.Handler
		//Last middleware will be the innermost middleware and executed last.
		for i := len(route.Middlewares) - 1; i >= 0; i-- {
			h = route.Middlewares[i](h)
		}
		r.Path(route.Path).Methods(route.Method).Handler(h).Queries(route.Queries...)
	}

	if conf.NotFoundHandler != nil {
		r.PathPrefix("/").Handler(conf.NotFoundHandler)
	}

	return r, nil
}

//GetEmptyConfig - returns empty config to be set
func GetEmptyConfig() *Config {
	return &Config{}
}
