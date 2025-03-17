package main

import (
	"log/slog"
	"os"
	"sso/intenal/config"
	"sso/intenal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	// TODO: инициализация конфигурации
	cfg := config.MustLoad()

	// TODO: инициализация логгера
	log := setupLogger(cfg.Env)

	log.Info("starting sso service",
		slog.String("env", cfg.Env),
		slog.Any("config", cfg),
	)

	// TODO: инициализация приложения(app)

	// TODO: запуск gRPC приложения(app)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
