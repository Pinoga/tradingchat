package app

import (
	"fmt"
	"net/http"
	"tradingchat/pkg/model"
	"tradingchat/pkg/util"

	"github.com/gorilla/sessions"
)

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := app.SessionStore.Get(r, "cookie-name")

	r.ParseMultipartForm(1 << 16)
	user := r.FormValue("user")
	password := r.FormValue("password")

	if user == "" || password == "" {
		http.Error(w, "user or password can't be empty", http.StatusBadRequest)
		return
	}

	foundUser, err := app.UserService.FindByUsername(user)
	if err != nil || !util.ComparePasswords(foundUser.Hash, []byte(password)) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	session.Values["user"] = foundUser.ID
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
		return
	}

	found, _ := app.UserService.FindByUsername(user)
	if found != nil {
		http.Error(w, "user already registered", http.StatusOK)
		return
	}

	uID, err := app.UserService.Register(user, password)
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	session.Values["user"] = uID
	saveSession(w, r, session)
	SendJSONResponse(w, "registration successful", http.StatusCreated, nil)
}

func (app *App) handleAuthentication(w http.ResponseWriter, r *http.Request) {
	session, err := app.SessionStore.Get(r, "cookie-name")
	if err != nil {
		fmt.Printf("couldn't decode session. Err: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// userIDString, ok := userID.(string)
	// if !ok {
	// 	http.Error(w, "not authenticated", http.StatusUnauthorized)
	// }

	if auth := app.isAuth(session); !auth {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	uID, ok := session.Values["user"].(string)
	if !ok {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	user, err := app.UserService.FindByID(uID)
	if err != nil {
		http.Error(w, "couldn't find user", http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, "authenticated", http.StatusOK, *user)
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

func (app *App) GetUserFromSession(session *sessions.Session) (*model.User, error) {
	uID, ok := session.Values["user"].(string)
	if !ok {
		return nil, fmt.Errorf("user not found in session")
	}

	user, err := app.UserService.FindByID(uID)
	if err != nil {
		return nil, fmt.Errorf("user not found in database")
	}
	return user, nil
}
