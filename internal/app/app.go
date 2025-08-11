package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/GazpachoGit/sso/internal/app/grpc"
)

// App is the main application structure that holds the gRPC server, storage, and other components.
type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	storagePath string,
	tokenTTL time.Duration) *App {
	//init storage

	//init auth service

	//init grpc app
	grpcApp := grpcapp.New(log, port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
