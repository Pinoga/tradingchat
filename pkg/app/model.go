package app

import (
	"log"
	"net/http"
	"tradingchat/pkg/chat"
	"tradingchat/pkg/mongodb"
	"tradingchat/pkg/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type AppConfig struct {
	LenBgs          int
	CookieSecretKey string
	Port            string
	DatabaseURI     string
	DatabaseName    string
}

type App struct {
	Router       *mux.Router
	Port         string
	Bgs          []chat.BroadcastGroup
	SessionStore *sessions.CookieStore
	UserService  service.UserService
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize(c AppConfig) *App {
	app.Port = c.Port
	app.Bgs = make([]chat.BroadcastGroup, c.LenBgs)
	key := []byte(c.CookieSecretKey)
	app.SessionStore = sessions.NewCookieStore(key)
	app.SessionStore.Options = &sessions.Options{
		MaxAge:   0,
		HttpOnly: false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	store := mongodb.NewStore(mongodb.MongoOptions{
		DB:  c.DatabaseName,
		URI: c.DatabaseURI,
	})

	app.UserService = service.NewUserService(store)

	app.Router = mux.NewRouter()

	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	apiRouter.Use(app.loggingMiddleware)

	apiRouter.HandleFunc("/login", app.handleLogin).Methods("POST")
	apiRouter.HandleFunc("/register", app.handleRegister).Methods("POST")

	chatRouter := apiRouter.PathPrefix("/chat").Subrouter()
	chatRouter.Use(app.authenticationMiddleware)

	chatRouter.HandleFunc("/authenticate", app.handleAuthentication).Methods("POST")
	chatRouter.HandleFunc("/enter/{room}", app.handleEnterRoom).Methods("GET")
	// chatRouter.HandleFunc("/leave", app.handleLeaveRoom).Methods("POST")

	app.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	return app
}

func (app *App) Run() error {
	for _, bg := range app.Bgs {
		go bg.HandleBroadcasts()
	}
	http.Handle("/", app.Router)
	log.Printf("Listening on %s", app.Port)
	err := http.ListenAndServe(app.Port, nil)
	return err
}
