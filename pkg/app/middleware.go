package app

import (
	"fmt"
	"log"
	"net/http"
)

func (app *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (app *App) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.SessionStore.Get(r, "cookie-name")
		fmt.Println(r.Cookies())

		fmt.Println(r.Header.Values("Cookie"))
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusUnauthorized)
			return
		}
		fmt.Println("passou do auth")

		next.ServeHTTP(w, r)
	})
}
