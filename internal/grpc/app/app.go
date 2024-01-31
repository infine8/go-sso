package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"sso/internal/grpc/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
    log        *slog.Logger
    gRPCServer *grpc.Server
    port       int
}

func New(log *slog.Logger, authService grpcauth.Auth, port int) *App {

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
        recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(interceptorLogger(log), loggingOpts...),
    ))

	grpcauth.Register(gRPCServer, authService)

    return &App{
        log:        log,
        gRPCServer: gRPCServer,
        port:       port,
    }
}

func (a *App) MustRun() {
    if err := a.Run(); err != nil {
        panic(err)
    }
}

func (a *App) Run() error {
    const op = "grpcapp.Run"

    l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

    a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

    if err := a.gRPCServer.Serve(l); err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

    return nil
}

func (a *App) GracefulStop() {
	a.gRPCServer.GracefulStop()
}