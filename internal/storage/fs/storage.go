package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

type DiskStorage struct {
	storagePath string
}

func NewDiskStorage(storagePath string) *DiskStorage {
	return &DiskStorage{storagePath}
}

func (ds *DiskStorage) Put(hash []byte) error {
	filePath := filepath.Join(ds.storagePath, fmt.Sprintf("%x", hash))
	_, err := os.Create(filePath)
	return err
}

func (ds *DiskStorage) Get(hash []byte) (string, error) {
	filePath := filepath.Join(ds.storagePath, fmt.Sprintf("%x", hash))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Delete удаляет данные с заданным хешем из хранилища на диске.
func (ds *DiskStorage) Delete(hash []byte) error {
	filePath := filepath.Join(ds.storagePath, fmt.Sprintf("%x", hash))
	return os.Remove(filePath)
}

// Clean удаляет все данные из хранилища на диске.
func (ds *DiskStorage) Clean() error {
	files, err := os.ReadDir(ds.storagePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Remove(filepath.Join(ds.storagePath, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
