package app

import (
	"log/slog"
	grpcapp "sso/internal/grpc/app"
	"sso/internal/service/auth"
	"sso/internal/storage/sqlite"
	"time"
)

type App struct {
	log *slog.Logger
	GRPCServer *grpcapp.App
	port int
	dbCloseConnection DbCloseConnection
}

type DbCloseConnection interface {
	CloseConnection()
}

func New(log *slog.Logger, storagePath string, port int, tokenTTL time.Duration) *App {

	storage, err := sqlite.New(storagePath)
	if (err != nil) {
		log.Error("Couldn't create storate")
	}

	auth := auth.New(log, storage, storage, storage, tokenTTL)

	grpcapp := grpcapp.New(log, auth, port)

	return &App{
		log: log,
		GRPCServer: grpcapp,
		port: port,
		dbCloseConnection: storage,
	}
}

func (a *App) Stop() {
    const op = "grpcapp.Stop"

    a.log.With(slog.String("op", op)).
        Info("stopping gRPC server", slog.Int("port", a.port))

    a.GRPCServer.GracefulStop()

	a.dbCloseConnection.CloseConnection()

	a.log.Info("Gracefully stopped")
}