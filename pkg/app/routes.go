package app

import "net/http"

type Route struct {
	desc    string
	method  string
	path    string
	auth    bool
	handler http.HandlerFunc
}

func NewRoute(desc, method, path string, auth bool, handler http.HandlerFunc) Route {
	return Route{
		desc:    desc,
		method:  method,
		path:    path,
		auth:    auth,
		handler: handler,
	}
}

func SecureRoute(desc, method, path string, handler http.HandlerFunc) Route {
	return NewRoute(desc, method, path, true, handler)
}

func PublicRoute(desc, method, path string, handler http.HandlerFunc) Route {
	return NewRoute(desc, method, path, false, handler)
}

// var routes = []Route{
// 	PublicRoute("Login", "POST", "/api/login", handleLogin)
// }
