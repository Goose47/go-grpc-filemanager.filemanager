package storage

import "errors"

var (
	ErrFileNotFound = errors.New("file is not found")
	ErrFileExists   = errors.New("file already exists")
)
