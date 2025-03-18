package app

import (
	"context"
	"log/slog"
	grpcapp "sso/intenal/app/grpc"
	"sso/intenal/config"
	"sso/intenal/services/auth"
	"sso/intenal/storage/postgres"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	db config.DBConfig,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.New(context.Background(), db)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
