package app

import (
	"fmt"
	"net/http"
	"tradingchat/pkg/service"

	"github.com/gorilla/sessions"
)

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := app.SessionStore.Get(r, "cookie-name")

	r.ParseMultipartForm(1 << 16)
	user := r.FormValue("user")
	password := r.FormValue("password")

	if user == "" || password == "" {
		http.Error(w, "user or password can't be empty", http.StatusBadRequest)
	}

	// if user != u || password != p {
	// 	http.Error(w, "user or password invalid", http.StatusForbidden)
	// }

	foundUser, err := app.UserService.FindByName(user)
	if err != nil {
		http.Error(w, "user or password invalid", http.StatusUnauthorized)
	}

	if passwordsMatch := service.ComparePasswords(foundUser.Hash, password); !passwordsMatch {
		http.Error(w, "invalid password", http.StatusUnauthorized)
	}

	session.Values["user"] = foundUser
	saveSession(w, r, session)
	SendJSONResponse(w, "login successful", http.StatusOK, nil)
}

func (app *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	session, _ := app.SessionStore.Get(r, "cookie-name")

	r.ParseMultipartForm(1 << 16)
	user := r.FormValue("user")
	password := r.FormValue("password")

	if user == "" || password == "" {
		http.Error(w, "user or password can't be empty", http.StatusBadRequest)
	}

	foundUser, err := app.UserService.FindByName(user)
	if err == nil {
		http.Error(w, "user already registered", http.StatusBadRequest)
	}

	created, err := app.UserService.Register(user, password)

	session.Values["user"] = created
	saveSession(w, r, session)
	SendJSONResponse(w, "registration successful", http.StatusCreated, nil)
}

func (app *App) handleAuthentication(w http.ResponseWriter, r *http.Request) {
	session, err := app.SessionStore.Get(r, "cookie-name")
	if err != nil {
		fmt.Printf("couldn't decode session. Err: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	// user, ok := session.Values["user"]V
	// if !ok {
	// 	http.Error(w, "not authenticated", http.StatusUnauthorized)
	// }
	// userIDString, ok := userID.(string)
	// if !ok {
	// 	http.Error(w, "not authenticated", http.StatusUnauthorized)
	// }

	if auth := app.isAuth(session); !auth {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
	}

	// user, err := app.UserRepo.Find(userIDString)
	// if err != nil {
	// 	http.Error(w, "not authenticated", http.StatusUnauthorized)
	// }

	SendJSONResponse(w, user, http.StatusOK, user)
}

func saveSession(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func (app *App) isAuth(session *sessions.Session) bool {
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
}
