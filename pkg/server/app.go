package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"star-wars/pkg/planet"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	router    *mux.Router
	container *container
	server    *http.Server
}

func (a *App) Start(ctx context.Context) {
	port := viper.GetString("port")
	log.Printf("Application started at port: %s", port)
	a.server = &http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     a,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}
	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
func (a App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func NewApp() *App {
	var app App

	configureLog(viper.GetString("log_level"))
	planetService := planet.NewService(mongoCollection(), viper.GetDuration("MONGO_TIMEOUT"))
	container := NewContainer(planetService)
	app.container = container

	app.RegisterRoutes()

	return &app
}

func mongoCollection() *driver.Collection {

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	client, err := driver.Connect(ctx, options.Client().ApplyURI(viper.GetString("MONGO_URI")))

	if err != nil {
		log.Fatal("Error trying to connect to the database.", err)
	}

	return client.Database(viper.GetString("MONGO_DB")).Collection(viper.GetString("MONGO_COLLECTION"))
}

func (a *App) RegisterRoutes() {
	router := mux.Router{}

	router.Use(a.HTTPServerMetricMiddleware)
	router.Use(a.RequestIdMiddleware)

	router.Handle("/health", a.healthHandler()).Methods(http.MethodGet)
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	router.Handle("/v1/planets", a.handleCreatePlanet(a.container.planetInserter)).Methods(http.MethodPost)
	router.Handle("/v1/planets/{id}", a.handleGetPlanetByID(a.container.planetGetter)).Methods(http.MethodGet)
	router.Handle("/v1/planets/{id}", a.handleUpdatePlanet(a.container.planetUpdater)).Methods(http.MethodPut)
	a.router = &router
}

func (a App) healthHandler() http.HandlerFunc {
	type response struct {
		Message string
	}
	return func(w http.ResponseWriter, request *http.Request) {
		writeJsonResponse(w, http.StatusOK, response{
			Message: "UP",
		})
	}
}
