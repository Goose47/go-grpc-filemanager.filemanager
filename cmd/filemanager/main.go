package main

import (
	"filemanager/internal/app"
	"filemanager/internal/config"
	"filemanager/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath)
	application.GRPCServer.MustRun()
}
