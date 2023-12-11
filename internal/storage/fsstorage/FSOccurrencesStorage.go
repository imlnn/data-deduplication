package fsstorage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type FSOccurrencesStorage struct {
	directory string
}

func NewFSStorage(fileName string) *FSOccurrencesStorage {
	return &FSOccurrencesStorage{directory: fileName}
}

func (fs *FSOccurrencesStorage) GetDirectory() string {
	return fs.directory
}

func (fs *FSOccurrencesStorage) SetDirectory(fileName string) {
	fs.directory = fileName
}

func (fs *FSOccurrencesStorage) Put(hash []byte, occurrence int) error {
	filePath := filepath.Join(fs.directory, fmt.Sprintf("%x", hash))
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", occurrence))
	return err
}

func (fs *FSOccurrencesStorage) Get(hash []byte) ([]int, error) {
	filePath := filepath.Join(fs.directory, fmt.Sprintf("%x", hash))
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var occurrences []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		occurrence, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		occurrences = append(occurrences, occurrence)
	}
	return occurrences, scanner.Err()
}

func (fs *FSOccurrencesStorage) Delete(hash []byte) error {
	filePath := filepath.Join(fs.directory, fmt.Sprintf("%x", hash))
	return os.Remove(filePath)
}

func (fs *FSOccurrencesStorage) Clean() error {
	return os.RemoveAll(fs.directory)
}
