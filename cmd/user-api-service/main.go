package main

import (
	"fmt"
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
	log.Info("init user Api service")

	// connect postgres
	db, err := postgres.Connect(cfg.Database)

	fmt.Println(cfg.Database)
	_, _ = db, err
	//TODO init router
	//TODO init server
}
