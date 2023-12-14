package fsstorage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type FSOccurrencesStorage struct {
	directory string
}

func NewFSStorage(directory string) *FSOccurrencesStorage {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/NewFSStorage"
	log.Printf("[%s] Creating new FSOccurrencesStorage with directory: %s", fn, directory)
	return &FSOccurrencesStorage{directory: directory}
}

func (fs *FSOccurrencesStorage) GetDirectory() string {
	return fs.directory
}

func (fs *FSOccurrencesStorage) SetDirectory(directory string) {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/SetDirectory"
	log.Printf("[%s] Setting directory to: %s", fn, directory)
	fs.directory = directory
}

func (fs *FSOccurrencesStorage) Put(hash []byte, occurrence int) error {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Put"
	filePath := filepath.Join(fs.directory, fmt.Sprintf("%x", hash))

	log.Printf("[%s] Storing occurrence for hash %x: %d", fn, hash, occurrence)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("[%s] Error opening file: %s", fn, err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", occurrence))
	if err != nil {
		log.Printf("[%s] Error writing occurrence: %s", fn, err)
	}

	return err
}

func (fs *FSOccurrencesStorage) Get(hash []byte) ([]int, error) {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Get"
	filePath := filepath.Join(fs.directory, fmt.Sprintf("%x", hash))

	log.Printf("[%s] Retrieving occurrences for hash %x", fn, hash)

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("[%s] Error opening file: %s", fn, err)
		return nil, err
	}
	defer file.Close()

	var occurrences []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		occurrence, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("[%s] Error converting text to int: %s", fn, err)
			return nil, err
		}
		occurrences = append(occurrences, occurrence)
	}

	if scanner.Err() != nil {
		log.Printf("[%s] Error scanning file: %s", fn, scanner.Err())
	}

	return occurrences, scanner.Err()
}

func (fs *FSOccurrencesStorage) Delete(hash []byte) error {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Delete"
	filePath := filepath.Join(fs.directory, fmt.Sprintf("%x", hash))

	log.Printf("[%s] Deleting file for hash %x", fn, hash)

	return os.Remove(filePath)
}

func (fs *FSOccurrencesStorage) Clean() error {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Clean"
	log.Printf("[%s] Cleaning directory: %s", fn, fs.directory)
	return os.RemoveAll(fs.directory)
}
