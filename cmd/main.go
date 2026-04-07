package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/nlsnnn/berezhok-notifications/internal/adapters"
	"github.com/nlsnnn/berezhok-notifications/internal/config"
	"github.com/nlsnnn/berezhok-notifications/internal/consumer"
	"github.com/nlsnnn/berezhok-notifications/internal/processor"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("start app", slog.String("env", cfg.Env))

	conn, err := consumer.GetConn(cfg.RabbitMQ.URL)
	if err != nil {
		log.Error("failed to connect to RabbitMQ", "err", err)
		return
	}

	err = conn.Setup()
	if err != nil {
		log.Error("failed to setup RabbitMQ", "err", err)
		return
	}

	resend := adapters.NewResendClient(cfg.Email.ResendApiKey, cfg.Email.From)

	dispatcher := processor.NewDispatcher(log)
	dispatcher.Register(processor.TypeEmail, processor.NewEmailProcessor(resend))

	consumer := consumer.New(conn.Channel, dispatcher, log)

	if err := consumer.Start(ctx); err != nil {
		log.Error("consumer stopped", "err", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
