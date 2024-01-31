package grpcapp

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	ssov1 "github.com/infine8/go-sso-proto/gen/go/sso"
)

func interceptorLogger(l *slog.Logger) logging.Logger {
    return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := fields[len(fields) - 1]

		if req, ok := f.(*ssov1.RegisterRequest); ok {
			fields = fields[:len(fields) - 2]
			newReq := *req
			newReq.Password = ""
			fields = append(fields, newReq)
		}

		if req, ok := f.(*ssov1.LoginRequest); ok {
			fields = fields[:len(fields) - 2]
			newReq := *req
			newReq.Password = ""
			fields = append(fields, newReq)
		}

        l.Log(ctx, slog.Level(lvl), msg, fields...)
    })
}