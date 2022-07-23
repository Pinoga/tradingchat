package app

import (
	"fmt"
	"log"
	"net/http"
	"tradingchat/pkg/chat"
	"tradingchat/pkg/mongodb"
	"tradingchat/pkg/repo"
	"tradingchat/pkg/service"
	"tradingchat/pkg/store"

	"github.com/go-redis/redis/v8"
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
	RedisHost       string
	RedisPort       string
}

type App struct {
	Router             *mux.Router
	Bgs                []*chat.BroadcastGroup
	Store              store.Store
	SessionStore       *sessions.CookieStore
	RedisStore         *redis.Client
	UserService        service.UserService
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

	apiRouter := app.Router.PathPrefix("/v1").Subrouter()
	apiRouter.Use(app.loggingMiddleware)

	apiRouter.HandleFunc("/login", app.handleLogin).Methods("POST")
	apiRouter.HandleFunc("/register", app.handleRegister).Methods("POST")
	apiRouter.HandleFunc("/authenticate", app.handleAuthentication).Methods("POST")

	chatRouter := apiRouter.PathPrefix("/chat").Subrouter()
	chatRouter.Use(app.authenticationMiddleware)
	chatRouter.HandleFunc("/{room}/enter", app.handleEnterRoom).Methods("POST")

	app.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	http.Handle("/", app.Router)

	app.RedisStore = redis.NewClient(&redis.Options{
		Addr:     c.RedisHost,
		Password: "",
		DB:       0,
	})

	store, err := mongodb.NewStore(mongodb.MongoOptions{
		DB:  app.DatabaseName,
		URI: app.DatabaseURI,
	})
	if err != nil {
		panic("could not connect to MongoDB")
	}
	fmt.Println("Connected to MongoDB successfully")

	app.Store = store
	userRepository := repo.NewUserRepository("user", store)
	app.UserService = service.NewUserService(userRepository)

	rabbitMQConnection, err := amqp.Dial(app.RabbitMQURI)
	if err != nil {
		panic("could not connect to RabbitMQ")
	}
	fmt.Println("Connected to RabbitMQ successfully")

	app.RabbitMQConnection = rabbitMQConnection
	ch, err := rabbitMQConnection.Channel()
	if err != nil {
		panic("could not create RabbitMQ channel")
	}
	defer ch.Close()

	msgs, err := app.ConsumeQueueMessages(ch)
	if err != nil {
		panic("could not connect to RabbitMQ queue")
	}
	go app.HandleConsumedMessages(msgs)

	for _, bg := range app.Bgs {
		go bg.HandleBroadcasts()
	}

	log.Fatal(http.ListenAndServe(":"+app.AppConfig.Port, app.Router))
	fmt.Printf("Listening on %s\n", app.AppConfig.Port)

	return app
}

func (app *App) Stop() {
	app.Store.Disconnect()
}
