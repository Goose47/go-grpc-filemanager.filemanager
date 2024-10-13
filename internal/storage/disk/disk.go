package disk

import (
	"context"
	"filemanager/internal/models"
	"filemanager/internal/storage"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

type Disk struct {
	BasePath string
}

func New(storagePath string) *Disk {
	return &Disk{
		BasePath: storagePath,
	}
}

// SaveFile saves file from byte array
func (d *Disk) SaveFile(ctx context.Context, filename string, data []byte) (int, error) {
	const op = "disk.SaveFile"

	fullPath := path.Join(d.BasePath, filename)

	if _, err := os.Stat(fullPath); err == nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrFileExists)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return n, nil
}

// File reads file specified in filename and returns its contents
func (d *Disk) File(ctx context.Context, filename string) (io.ReadCloser, error) {
	const op = "disk.SaveFile"

	fullPath := path.Join(d.BasePath, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrFileNotFound)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrFileNotFound)
	}

	return file, nil
}

// Files retrieves all uploaded files creation and update date
func (d *Disk) Files(ctx context.Context) ([]models.File, error) {
	const op = "disk.Files"

	var files []models.File

	err := filepath.Walk(d.BasePath, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			modTime, changeTime, err := getFileTimes(fullPath)
			if err != nil {
				return fmt.Errorf("error retrieving times for file %s: %v", fullPath, err)
			}

			files = append(files, models.File{
				Name:         path.Base(fullPath),
				CreationDate: modTime.Unix(),
				UpdateDate:   changeTime.Unix(),
			})
		}
		return nil
	})

	if err != nil {
		return []models.File{}, fmt.Errorf("%s: %w", op, err)
	}

	return files, nil
}

// getFileTimes retrieves the file modification time and estimated creation time (if available).
// Specific for linux systems
func getFileTimes(filePath string) (modTime, changeTime time.Time, err error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	modTime = fileInfo.ModTime()
	stat := fileInfo.Sys().(*syscall.Stat_t)
	changeTime = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)) // closest approximation to creation time

	return modTime, changeTime, nil
}
