package storage

import (
	"context"
	"errors"
	"filemanager/internal/lib/random"
	"filemanager/internal/models"
	"filemanager/internal/storage"
	"fmt"
	"io"
	"log/slog"
)

type FileSaver interface {
	SaveFile(ctx context.Context, filename string) (io.WriteCloser, error)
}
type FileProvider interface {
	File(ctx context.Context, filename string) (io.ReadCloser, error)
	Files(ctx context.Context) ([]models.File, error)
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

func (s *Storage) FileWriter(ctx context.Context) (io.WriteCloser, string, error) {
	const op = "storage.SaveFile"

	filename := random.String(20)

	log := s.log.With(slog.String("op", op), slog.String("filename", filename))

	log.Info("trying to create file")

	file, err := s.fileSaver.SaveFile(ctx, filename)

	if err != nil {
		log.Error("failed to create file", slog.Any("error", err))

		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file created successfully")

	return file, filename, nil
}

func (s *Storage) ListFiles(ctx context.Context) ([]models.File, error) {
	const op = "storage.ListFiles"

	log := s.log.With(slog.String("op", op))

	log.Info("retrieving files")

	files, err := s.fileProvider.Files(ctx)
	if err != nil {
		log.Error("failed to retrieve files", slog.Any("error", err))

		return []models.File{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("files retrieved successfully")

	return files, nil
}

func (s *Storage) FileReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	const op = "storage.File"

	log := s.log.With(slog.String("op", op), slog.String("filename", filename))

	log.Info("trying to find file")

	file, err := s.fileProvider.File(ctx, filename)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			log.Warn("file is not found", slog.Any("error", err))

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to retrieve file", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file found successfully")

	return file, nil
}
