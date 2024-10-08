package app

import (
	grpcapp "filemanager/internal/app/grpc"
	"filemanager/internal/services/storage"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
) *App {

	storageService := storage.New(log, nil, nil)
	grpcApp := grpcapp.New(log, storageService, grpcPort)

	return &App{
		grpcApp,
	}
}
