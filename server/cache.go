package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrCacheNotFound = errors.New("cache not found")

type Cache interface {
	Init() error
	Exists(key string) bool
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

type FileStoreCache struct {
	path string
}

func (fsc *FileStoreCache) Init() error {
	tempDir, err := os.MkdirTemp("", "file-previews-cache-*")
	if err != nil {
		return fmt.Errorf("error creating cache temporary directory: %w", err)
	}

	fsc.path = tempDir

	return nil
}

func (fsc *FileStoreCache) joinPath(key string) string {
	return filepath.Join(fsc.path, key)
}

func (fsc *FileStoreCache) Exists(key string) bool {
	_, err := os.Stat(fsc.joinPath(key))
	return os.IsExist(err)
}

func (fsc *FileStoreCache) Get(key string) ([]byte, error) {
	if !fsc.Exists(key) {
		return nil, ErrCacheNotFound
	}

	return nil, nil
}

func (fsc *FileStoreCache) Set(key string, value []byte) error {
	return nil
}

func NewFileStoreCache() Cache {
	return &FileStoreCache{}
}
