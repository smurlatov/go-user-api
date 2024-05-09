package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"user-api-service/internals/config"
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

	fmt.Println(cfg.Database)
	_, _ = db, err

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	//TODO init server
}
