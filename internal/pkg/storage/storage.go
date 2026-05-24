package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sync"
)

var ErrInvalidJSON = errors.New("Invalid JSON")

type Storage interface {
	Load(data map[string]any) error
	LoadFromPath() error
	save(newStorage map[string]any) error
	SetPath(path string)

	Set(key string, value any) error
	Delete(key string) error
	Get(key string) (any, error)
	List() map[string]any
}

type FileStore struct {
	path    string
	storage map[string]any
	mu      sync.RWMutex
}

func NewFileStore() *FileStore {
	return &FileStore{
		path:    "",
		storage: make(map[string]any),
	}
}

func (fileStore *FileStore) SetPath(path string) {
	fileStore.mu.Lock()
	fileStore.path = path
	fileStore.mu.Unlock()
}

func (fileStore *FileStore) LoadData(data map[string]any) error {
	newStorage := make(map[string]any, len(data))
	maps.Copy(newStorage, data)

	if err := fileStore.save(newStorage); err != nil {
		return err
	}

	fileStore.mu.Lock()
	fileStore.storage = newStorage
	fileStore.mu.Unlock()

	return nil
}

func (fileStore *FileStore) LoadFromPath() error {
	fileStore.mu.RLock()
	path := fileStore.path
	fileStore.mu.RUnlock()

	if path == "" {
		return fmt.Errorf("Path is not set")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var newStorage map[string]any
	if err := json.Unmarshal(data, &newStorage); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	if err := fileStore.LoadData(newStorage); err != nil {
		return err
	}

	return nil
}

func (fileStore *FileStore) save(newStorage map[string]any) error {
	fileStore.mu.RLock()
	path := fileStore.path
	fileStore.mu.RUnlock()

	if path == "" {
		return fmt.Errorf("Path is not set")
	}

	data, err := json.Marshal(newStorage)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "filestore-*.tmp")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()
	cleanup := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	if _, err := tmpFile.Write(data); err != nil {
		cleanup()
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		cleanup()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		cleanup()
		return err
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}

func (fileStore *FileStore) Get(key string) (any, error) {
	fileStore.mu.RLock()
	value, ok := fileStore.storage[key]
	fileStore.mu.RUnlock()

	if !ok {
		return nil, nil
	}

	return value, nil
}

func (fileStore *FileStore) Set(key string, value any) error {
	fileStore.mu.RLock()
	newStorage := make(map[string]any, len(fileStore.storage))
	maps.Copy(newStorage, fileStore.storage)
	newStorage[key] = value
	fileStore.mu.RUnlock()

	if err := fileStore.save(newStorage); err != nil {
		return err
	}

	fileStore.mu.Lock()
	fileStore.storage = newStorage
	fileStore.mu.Unlock()
	return nil
}

func (fileStore *FileStore) Delete(key string) error {
	fileStore.mu.RLock()

	if _, ok := fileStore.storage[key]; !ok {
		fileStore.mu.RUnlock()
		return nil
	}

	newStorage := make(map[string]any, len(fileStore.storage))
	maps.Copy(newStorage, fileStore.storage)
	delete(newStorage, key)
	fileStore.mu.RUnlock()

	if err := fileStore.save(newStorage); err != nil {
		return err
	}

	fileStore.mu.Lock()
	fileStore.storage = newStorage
	fileStore.mu.Unlock()
	return nil
}

func (fileStore *FileStore) List() map[string]any {
	fileStore.mu.RLock()
	defer fileStore.mu.RUnlock()
	storage := make(map[string]any, len(fileStore.storage))
	maps.Copy(storage, fileStore.storage)
	return storage
}
