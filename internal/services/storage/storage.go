package storage

import (
	"context"
	"filemanager/internal/models"
	"log/slog"
)

type FileSaver interface{}
type FileProvider interface{}

type Storage struct {
	log          *slog.Logger
	fileSaver    FileSaver
	fileProvider FileProvider
}

func New(
	log *slog.Logger,
	fileSaver FileSaver,
	fileProvider FileProvider,
) *Storage {
	return &Storage{
		log:          log,
		fileSaver:    fileSaver,
		fileProvider: fileProvider,
	}
}

func (s *Storage) SaveFile(fileData []byte) (bytes int, filename string) {
	return 1, "test"
}
func (s *Storage) ListFiles(ctx context.Context) ([]models.File, error) {
	files := []models.File{
		{Name: "name", CreationDate: 1, UpdateDate: 1},
	}
	return files, nil
}
func (s *Storage) File(ctx context.Context, filename string) ([]byte, error) {
	return make([]byte, 0), nil
}
