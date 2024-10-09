package storage

import (
	"context"
	"errors"
	"filemanager/internal/lib/random"
	"filemanager/internal/models"
	"filemanager/internal/storage"
	"fmt"
	"log/slog"
)

type FileSaver interface {
	SaveFile(filename string, data []byte) (int, error)
}
type FileProvider interface {
	File() ([]byte, error)
	Files() ([]models.File, error)
}

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

func (s *Storage) SaveFile(ctx context.Context, fileData []byte) (int, string, error) {
	const op = "storage.SaveFile"

	filename := random.String(20)

	log := s.log.With(slog.String("op", op), slog.String("filename", filename))

	log.Info("trying to save file")

	bytes, err := s.fileSaver.SaveFile(filename, fileData)
	if err != nil {
		log.Error("failed to save file", slog.Any("error", err))

		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file saved successfully")

	return bytes, filename, nil
}

func (s *Storage) ListFiles(ctx context.Context) ([]models.File, error) {
	const op = "storage.ListFiles"

	log := s.log.With(slog.String("op", op))

	log.Info("retrieving files")

	files, err := s.fileProvider.Files()
	if err != nil {
		log.Error("failed to retrieve files", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("files retrieved successfully")

	return files, nil
}

func (s *Storage) File(ctx context.Context, filename string) ([]byte, error) {
	const op = "storage.File"

	log := s.log.With(slog.String("op", op), slog.String("filename", filename))

	log.Info("trying to find file")

	fileData, err := s.fileProvider.File()
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			log.Warn("file is not found", slog.Any("error", err))

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to retrieve file", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return fileData, nil
}
