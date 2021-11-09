package app

import (
	"context"
	"log"
	"net/http"
	"time"
)

type UserContextKey string

func (app *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (app *App) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.SessionStore.Get(r, "cookie-name")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		user, err := app.GetUserFromSession(session)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
		ctx = context.WithValue(ctx, UserContextKey("user"), *user)
		defer cancel()
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
