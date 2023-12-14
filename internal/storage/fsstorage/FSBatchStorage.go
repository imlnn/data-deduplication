package fsstorage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FSBatchStorage struct {
	directory string
	mu        sync.Mutex // Добавляем мьютекс для синхронизации доступа к файлам
}

func NewFSBatchStorage(directory string) *FSBatchStorage {
	const fn = "internal/storage/fsstorage/FSBatchStorage/NewFSBatchStorage"
	log.Printf("[%s] Creating new FSBatchStorage with directory: %s", fn, directory)

	storage := &FSBatchStorage{}

	storage.SetDirectory(directory)

	return storage
}

func (fs *FSBatchStorage) GetDirectory() string {
	return fs.directory
}

func (fs *FSBatchStorage) SetDirectory(directory string) {
	const fn = "internal/storage/fsstorage/FSBatchStorage/SetDirectory"
	log.Printf("[%s] Setting directory to: %s", fn, directory)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[%s] %s", fn, err)
	}

	dir := filepath.Join(wd, directory)
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(dir, 0777); err != nil {
				log.Fatalf("[%s] %s", fn, err)
			}
		} else {
			log.Fatalf("[%s] %s", fn, err)
		}
	}

	fs.directory = directory
}

// createFilename создает безопасное имя файла на основе хэша.
func (fs *FSBatchStorage) createFilename(hash []byte) string {
	return filepath.Join(fs.directory, fmt.Sprintf("%x", hash))
}

func (fs *FSBatchStorage) Put(hash []byte, batch []byte) error {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Put"
	filename := fs.createFilename(hash)

	fs.mu.Lock() // Блокируем доступ к файлам
	defer fs.mu.Unlock()

	if _, err := os.Stat(filename); err == nil {
		log.Printf("[%s] Batch already exists with hash: %s", fn, fmt.Sprintf("%x", hash))
		return errors.New("batch already exists")
	} else if err != nil && !os.IsNotExist(err) {
		return err // Возвращаем ошибку файловой системы
	}

	log.Printf("[%s] Storing batch with hash: %s", fn, fmt.Sprintf("%x", hash))
	return os.WriteFile(filename, batch, 0777)
}

func (fs *FSBatchStorage) Get(hash string) ([]byte, error) {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Get"
	filename := filepath.Join(fs.directory, fmt.Sprintf("%s", hash))

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("[%s] Batch does not exist with hash: %s", fn, fmt.Sprintf("%x", hash))
		return nil, errors.New("batch does not exist")
	} else if err != nil {
		return nil, err // Возвращаем ошибку файловой системы
	}

	log.Printf("[%s] Retrieving batch with hash: %s", fn, fmt.Sprintf("%x", hash))
	return os.ReadFile(filename)
}

func (fs *FSBatchStorage) Delete(hash []byte) error {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Delete"
	filename := fs.createFilename(hash)

	fs.mu.Lock() // Блокируем доступ к файлам
	defer fs.mu.Unlock()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("[%s] Batch does not exist with hash: %s", fn, fmt.Sprintf("%x", hash))
		return errors.New("batch does not exist")
	} else if err != nil {
		return err // Возвращаем ошибку файловой системы
	}

	log.Printf("[%s] Deleting batch with hash: %s", fn, fmt.Sprintf("%x", hash))
	return os.Remove(filename)
}

func (fs *FSBatchStorage) Clean() error {
	const fn = "internal/storage/fsstorage/FSBatchStorage/Clean"
	log.Printf("[%s] Cleaning directory: %s", fn, fs.directory)

	fs.mu.Lock() // Блокируем доступ к файлам
	defer fs.mu.Unlock()

	// Удаляем только файлы, созданные хранилищем
	files, err := os.ReadDir(fs.directory)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := os.Remove(filepath.Join(fs.directory, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
