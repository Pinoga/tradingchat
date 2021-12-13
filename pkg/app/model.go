package app

import (
	"fmt"
	"net/http"
	"tradingchat/pkg/chat"
	"tradingchat/pkg/mongodb"
	"tradingchat/pkg/service"
	"tradingchat/pkg/store"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/streadway/amqp"
)

type AppConfig struct {
	LenBgs          int
	CookieSecretKey string
	Port            string
	DatabaseURI     string
	DatabaseName    string
	RabbitMQURI     string
}

type App struct {
	Router             *mux.Router
	Bgs                []*chat.BroadcastGroup
	SessionStore       *sessions.CookieStore
	UserService        service.UserService
	Store              store.Store
	RabbitMQConnection *amqp.Connection
	AppConfig
}

func NewApp() *App {
	return &App{}
}

func (app *App) Initialize(c AppConfig) *App {
	app.AppConfig = c
	app.Bgs = make([]*chat.BroadcastGroup, c.LenBgs)

	for i := 0; i < c.LenBgs; i++ {
		app.Bgs[i] = chat.NewBroadCastGroup()
	}

	key := []byte(c.CookieSecretKey)
	app.SessionStore = sessions.NewCookieStore(key)
	app.SessionStore.Options = &sessions.Options{
		MaxAge:   0,
		HttpOnly: false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	app.Router = mux.NewRouter()

	apiRouter := app.Router.PathPrefix("/api").Subrouter()
	apiRouter.Use(app.loggingMiddleware)

	apiRouter.HandleFunc("/login", app.handleLogin).Methods("POST")
	apiRouter.HandleFunc("/register", app.handleRegister).Methods("POST")
	apiRouter.HandleFunc("/authenticate", app.handleAuthentication).Methods("POST")

	chatRouter := apiRouter.PathPrefix("/chat").Subrouter()
	chatRouter.Use(app.authenticationMiddleware)

	chatRouter.HandleFunc("/enter/{room}", app.handleEnterRoom).Methods("GET")

	app.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	return app
}

func (app *App) Run() error {
	store, err := mongodb.NewStore(mongodb.MongoOptions{
		DB:  app.DatabaseName,
		URI: app.DatabaseURI,
	})
	if err != nil {
		return err
	}
	fmt.Println("Connected to MongoDB successfully")

	app.Store = store
	app.UserService = service.NewUserService(store)

	rabbitMQConnection, err := amqp.Dial(app.RabbitMQURI)
	if err != nil {
		return err
	}
	fmt.Println("Connected to RabbitMQ successfully")

	app.RabbitMQConnection = rabbitMQConnection
	ch, err := rabbitMQConnection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := app.ConsumeQueueMessages(ch)
	if err != nil {
		return err
	}
	go app.HandleConsumedMessages(msgs)

	for _, bg := range app.Bgs {
		go bg.HandleBroadcasts()
	}

	http.Handle("/", app.Router)
	fmt.Printf("Listening on %s\n", app.AppConfig.Port)

	err = http.ListenAndServe(":"+app.AppConfig.Port, nil)
	return err
}

func (app *App) Stop() {
	app.Store.Disconnect()
}
