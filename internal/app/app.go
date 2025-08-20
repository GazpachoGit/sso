package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/GazpachoGit/sso/internal/app/grpc"
	"github.com/GazpachoGit/sso/internal/service/auth"
	"github.com/GazpachoGit/sso/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	//init auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)
	//init grpc app
	grpcApp := grpcapp.New(log, authService, port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
