package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/behummble/Questions-answers/internal/config"
	"github.com/behummble/Questions-answers/internal/handlers/http"
	"github.com/behummble/Questions-answers/internal/service"
	"github.com/behummble/Questions-answers/internal/storage/postgres"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cfg := config.NewConfig()
	log := newLog(cfg.Log)
	storage := postgres.NewStorage(ctx, log, cfg.Storage)
	service := service.NewService(log, storage, storage)
	server := http.NewServer(ctx, log, &cfg.Server, service)
	go server.Start()
	<- ctx.Done()
	srvContext, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	server.Shutdown(srvContext)
}

func newLog(config config.LogConfig) *slog.Logger {
	var output *os.File
	if config.Path != "" {
		file, err := os.OpenFile(config.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		output = file
	} else {
		output = os.Stdout
	}
	return slog.New(
		slog.NewJSONHandler(
			output,
			&slog.HandlerOptions{Level: slog.Level(config.Level)},
		),
	)
}
