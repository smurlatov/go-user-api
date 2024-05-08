package main

import (
	"log/slog"
	"os"
	"user-api-service/internals/config"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("init user Api service")
	_ = cfg
	//TODO init storage
	//TODO init router
	//TODO init server
}
