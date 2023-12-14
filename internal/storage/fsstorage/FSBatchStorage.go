package fsstorage

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type FSBatchStorage struct {
	directory string
}

func NewFSBatchStorage(directory string) *FSBatchStorage {
	const fn = "internal/storage/fsstorage/FSBatchStorage/NewFSBatchStorage"
	log.Printf("[%s] Creating new FSBatchStorage with directory: %s", fn, directory)
	return &FSBatchStorage{directory: directory}
}

func (fs *FSBatchStorage) GetDirectory() string {
	return fs.directory
}

func (fs *FSBatchStorage) SetDirectory(directory string) {
	const fn = "internal/storage/fsstorage/FSBatchStorage/SetDirectory"
	log.Printf("[%s] Setting directory to: %s", fn, directory)
	fs.directory = directory
}

func (fs *FSBatchStorage) Put(hash []byte, batch []byte) error {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Put"
	filename := filepath.Join(fs.directory, string(hash))

	if _, err := os.Stat(filename); err == nil {
		log.Printf("[%s] Batch already exists with hash: %s", fn, string(hash))
		return errors.New("batch already exists")
	}

	log.Printf("[%s] Storing batch with hash: %s", fn, string(hash))
	return os.WriteFile(filename, batch, 0644)
}

func (fs *FSBatchStorage) Get(hash []byte) ([]byte, error) {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Get"
	filename := filepath.Join(fs.directory, string(hash))

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("[%s] Batch does not exist with hash: %s", fn, string(hash))
		return nil, errors.New("batch does not exist")
	}

	log.Printf("[%s] Retrieving batch with hash: %s", fn, string(hash))
	return os.ReadFile(filename)
}

func (fs *FSBatchStorage) Delete(hash []byte) error {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Delete"
	filename := filepath.Join(fs.directory, string(hash))

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("[%s] Batch does not exist with hash: %s", fn, string(hash))
		return errors.New("batch does not exist")
	}

	log.Printf("[%s] Deleting batch with hash: %s", fn, string(hash))
	return os.Remove(filename)
}

func (fs *FSBatchStorage) Clean() error {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Clean"
	log.Printf("[%s] Cleaning directory: %s", fn, fs.directory)
	return os.RemoveAll(fs.directory)
}
