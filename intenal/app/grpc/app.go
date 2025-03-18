package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/intenal/grpc/auth"
	"sso/intenal/lib/logger/sl"
)

type App struct {
	log         *slog.Logger
	gRPCServer  *grpc.Server
	authService authgrpc.AuthSrv
	port        int
}

func New(log *slog.Logger, port int, authService authgrpc.AuthSrv) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Reg(gRPCServer, authService)

	return &App{
		log:         log,
		gRPCServer:  gRPCServer,
		authService: authService,
		port:        port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error("failed to listen", sl.Error(err))
		return fmt.Errorf("failed to listen: %s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("failed to serve", sl.Error(err))
		return fmt.Errorf("failed to serve: %s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
