package router

import (
	"net/http"
	"nlip/handlers"
)

type Route struct {
	Method  string
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var Routes = []Route{
	{Method: "POST", Path: "/NLIP/", Handler: handlers.InitiationHandler},
	{Method: "POST", Path: "/test/", Handler: handlers.TestHandler},
	{Method: "POST", Path: "/register/", Handler: handlers.Register},
	{Method: "POST", Path: "/login/", Handler: handlers.Login},
}

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range Routes {
		mux.HandleFunc(route.Method+" "+route.Path, route.Handler)
	}

	return mux
}
