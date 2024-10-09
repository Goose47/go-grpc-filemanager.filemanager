package app

import (
	grpcapp "filemanager/internal/app/grpc"
	"filemanager/internal/services/storage"
	"filemanager/internal/storage/disk"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	diskStorage := disk.New(storagePath)
	storageService := storage.New(log, diskStorage, diskStorage)
	grpcApp := grpcapp.New(log, storageService, grpcPort)

	return &App{
		grpcApp,
	}
}
