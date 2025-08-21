package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/GazpachoGit/sso/internal/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPC application (server)
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))
	log.Info("starting grpc server")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("stopping grpc server")
	a.gRPCServer.GracefulStop()
	log.Info("grpc server was stopped gracefully")
}
