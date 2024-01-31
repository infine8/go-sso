package main

import (
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	logger "sso/lib/logger/setup"
	"syscall"
)


func main(){
	cfg := config.MustLoad()
	
	log := logger.SetupLogger(cfg.Env)

	app := app.New(log, cfg.StoragePath, cfg.GRPC.Port, cfg.TokenTTL)

	go func ()  {
		app.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()
}