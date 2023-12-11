package fsstorage

import (
	"errors"
	"os"
	"path/filepath"
)

type FSBatchStorage struct {
	directory string
}

func NewFSBatchStorage(directory string) *FSBatchStorage {
	return &FSBatchStorage{directory: directory}
}

func (fs *FSBatchStorage) GetDirectory() string {
	return fs.directory
}

func (fs *FSBatchStorage) SetDirectory(fileName string) {
	fs.directory = fileName
}

func (fs *FSBatchStorage) Put(hash []byte, batch []byte) error {
	filename := filepath.Join(fs.directory, string(hash))

	if _, err := os.Stat(filename); err == nil {
		return errors.New("batch already exists")
	}

	return os.WriteFile(filename, batch, 0644)
}

func (fs *FSBatchStorage) Get(hash []byte) ([]byte, error) {
	filename := filepath.Join(fs.directory, string(hash))

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New("batch does not exist")
	}

	return os.ReadFile(filename)
}

func (fs *FSBatchStorage) Delete(hash []byte) error {
	filename := filepath.Join(fs.directory, string(hash))

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return errors.New("batch does not exist")
	}

	return os.Remove(filename)
}

func (fs *FSBatchStorage) Clean() error {
	return os.RemoveAll(fs.directory)
}
