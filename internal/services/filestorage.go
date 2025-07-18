package services

import (
	"io"
	"os"
	"path/filepath"
)

type fileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) *fileStorage {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}
	return &fileStorage{basePath: basePath}
}

func (s *fileStorage) Save(filename string, data io.Reader) error {
	path := filepath.Join(s.basePath, filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

func (s *fileStorage) GetPath(filename string) string {
	return filepath.Join(s.basePath, filename)
}
