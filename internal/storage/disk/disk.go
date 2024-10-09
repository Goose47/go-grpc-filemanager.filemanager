package disk

import (
	"bufio"
	"errors"
	"filemanager/internal/models"
	"filemanager/internal/storage"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"path"
	"path/filepath"
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
func (d *Disk) SaveFile(filename string, data []byte) (int, error) {
	const op = "disk.SaveFile"

	fullPath := path.Join(d.BasePath, filename)

	if _, err := os.Stat(fullPath); err != nil {
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
func (d *Disk) File(filename string) ([]byte, error) {
	const op = "disk.SaveFile"

	fullPath := path.Join(d.BasePath, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrFileNotFound)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrFileNotFound)
	}
	defer file.Close()

	var fileData []byte
	reader := bufio.NewReader(file)
	buf := make([]byte, 256)

	for {
		_, err := reader.Read(buf)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		fileData = append(fileData, buf...)
	}

	return fileData, nil
}

// Files retrieves all uploaded files creation and update date
func (d *Disk) Files() ([]models.File, error) {
	const op = "disk.Files"

	var files []models.File

	err := filepath.Walk(d.BasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			modTime, changeTime, err := getFileTimes(path)
			if err != nil {
				return fmt.Errorf("error retrieving times for file %s: %v", path, err)
			}

			files = append(files, models.File{
				Name:         path,
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
	stat := fileInfo.Sys().(*unix.Stat_t)
	changeTime = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)) // closest approximation to creation time

	return modTime, changeTime, nil
}
