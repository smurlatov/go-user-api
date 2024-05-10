package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"user-api-service/internals/config"
	"user-api-service/internals/http-server/handlers/user/get"
	"user-api-service/internals/http-server/handlers/user/save"
	"user-api-service/internals/http-server/handlers/user/update"
	"user-api-service/internals/storage/postgres"
)

func main() {
	// read config
	cfg := config.MustLoad()
	// init logger
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("init logger")

	// connect postgres
	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}
	defer db.Close()

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/users", save.New(log, db))
	router.Get("/user/{id}", get.New(log, db))
	router.Patch("/user/{id}", update.New(log, db))

	// init server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}
