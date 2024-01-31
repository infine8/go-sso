package suite

import (
	"context"
	"net"
	"sso/internal/config"
	"strconv"
	"testing"
	"time"

	ssov1 "github.com/infine8/go-sso-proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
    configPath = "/Users/infine_1/Documents/_PROJECTS/golang/sso-grpc/sso/config.yaml"
    grpcHost = "localhost"
)

type Suite struct {
    *testing.T                  // Потребуется для вызова методов *testing.T
    Cfg        *config.Config   // Конфигурация приложения
    AuthClient ssov1.AuthClient // Клиент для взаимодействия с gRPC-сервером Auth
}

func New(t *testing.T) (context.Context, *Suite) {
    t.Helper()   // Функция будет восприниматься как вспомогательная для тестов
    t.Parallel() // Разрешаем параллельный запуск тестов

    cfg := config.MustLoadPath(configPath)

    // Основной родительский контекст   
    ctx, cancelCtx := context.WithTimeout(context.Background(), time.Duration(1 * time.Second))

    // Когда тесты пройдут, закрываем контекст
    t.Cleanup(func() {
        t.Helper()
        cancelCtx()
    })

    // Адрес нашего gRPC-сервера
    grpcAddress := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
    cc, err := grpc.DialContext(context.Background(),
        grpcAddress,
        // Используем insecure-коннект для тестов
        grpc.WithTransportCredentials(insecure.NewCredentials())) 
    if err != nil {
        t.Fatalf("grpc server connection failed: %v", err)
    }

    authClient := ssov1.NewAuthClient(cc)

    return ctx, &Suite{
        T:          t,
        Cfg:        cfg,
        AuthClient: authClient,
    }
}