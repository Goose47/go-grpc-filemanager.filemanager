package main

import (
	"filemanager/internal/app"
	"filemanager/internal/config"
	"filemanager/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath)
	go func() {
		application.GRPCServer.MustRun()
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	application.GRPCServer.Stop()

	log.Info("Gracefully stopped")
}
