package fsstorage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type FSOccurrencesStorage struct {
	directory string
	mu        sync.Mutex
}

func NewFSStorage(directory string) *FSOccurrencesStorage {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/NewFSStorage"
	log.Printf("[%s] Creating new FSOccurrencesStorage with directory: %s", fn, directory)

	storage := &FSOccurrencesStorage{}
	storage.SetDirectory(directory)

	return storage
}

func (fs *FSOccurrencesStorage) GetDirectory() string {
	return fs.directory
}

func (fs *FSOccurrencesStorage) SetDirectory(directory string) {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/SetDirectory"
	log.Printf("[%s] Setting directory to: %s", fn, directory)

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("[%s] %s", fn, err)
		return
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

func (fs *FSOccurrencesStorage) createFilePath(hash []byte) string {
	return filepath.Join(fs.directory, fmt.Sprintf("%x", hash))
}

func (fs *FSOccurrencesStorage) SetWD(fileName string) {
	fs.SetDirectory(fmt.Sprintf("%s/%s", fs.directory, fileName))
}

func (fs *FSOccurrencesStorage) PutMetadata(sourceFilePath string, segments int, lastBatchSize int) error {
	fileName := filepath.Base(sourceFilePath)
	data := fmt.Sprintf("DEDUP_%s %d %d", fileName, segments, lastBatchSize)

	infoFilePath := filepath.Join(fs.directory, "info")
	info, err := os.Create(infoFilePath)
	if err != nil {
		return fmt.Errorf("error creating info file: %w", err)
	}
	defer info.Close()

	_, err = fmt.Fprint(info, data)
	if err != nil {
		return fmt.Errorf("error writing to info file: %w", err)
	}

	return nil
}

func (fs *FSOccurrencesStorage) GetMetadata() (string, int, int, error) {
	const fn = "internal/service/dedup/service/Restore/parseInfoFile"

	infoFilePath := filepath.Join(fs.directory, "info")
	info, err := os.Open(infoFilePath)
	if err != nil {
		return "", 0, 0, fmt.Errorf("error opening info file: %w", err)
	}
	defer info.Close()

	log.Printf("[%s] Parsing info file", fn)

	var data string
	scanner := bufio.NewScanner(info)
	if scanner.Scan() {
		data = scanner.Text()
	} else {
		return "", 0, 0, fmt.Errorf("[%s] Info File is empty", fn)
	}

	parts := strings.Split(data, " ")
	if len(parts) != 3 {
		return "", 0, 0, fmt.Errorf("[%s] Info File format is incorrect", fn)
	}

	fileName := parts[0]
	segmentsNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, 0, fmt.Errorf("[%s] Error parsing number: %w", fn, err)
	}

	lastBatchSize, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", 0, 0, fmt.Errorf("[%s] Error parsing number: %w", fn, err)
	}

	log.Printf("[%s] Parsed info: File: %s, Segments: %d", fn, fileName, segmentsNum)

	return fileName, segmentsNum, lastBatchSize, nil
}

func (fs *FSOccurrencesStorage) Put(hash []byte, occurrence int) error {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Put"
	filePath := fs.createFilePath(hash)

	fs.mu.Lock()
	defer fs.mu.Unlock()

	log.Printf("[%s] Storing occurrence for hash %x: %d", fn, hash, occurrence)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("[%s] Error opening file: %s", fn, err)
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	if _, err = file.WriteString(fmt.Sprintf("%d\n", occurrence)); err != nil {
		log.Printf("[%s] Error writing occurrence: %s", fn, err)
		return fmt.Errorf("error writing occurrence: %w", err)
	}
	return nil
}

func (fs *FSOccurrencesStorage) Get(hash string) ([]int, error) {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Get"
	filename := filepath.Join(fs.directory, fmt.Sprintf("%s", hash))

	fs.mu.Lock()
	defer fs.mu.Unlock()

	log.Printf("[%s] Retrieving occurrences for hash %x", fn, hash)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("[%s] Error opening file: %s", fn, err)
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var occurrences []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		occurrence, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("[%s] Error converting text to int: %s", fn, err)
			return nil, fmt.Errorf("error converting text to int: %w", err)
		}
		occurrences = append(occurrences, occurrence)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[%s] Error scanning file: %s", fn, err)
		return nil, fmt.Errorf("error scanning file: %w", err)
	}
	return occurrences, nil
}

func (fs *FSOccurrencesStorage) Delete(hash []byte) error {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Delete"
	filePath := fs.createFilePath(hash)

	fs.mu.Lock()
	defer fs.mu.Unlock()

	log.Printf("[%s] Deleting file for hash %x", fn, hash)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}
	return nil
}

func (fs *FSOccurrencesStorage) Clean() error {
	const fn = "internal/storage/fsstorage/FSOccurrencesStorage/Clean"
	fs.mu.Lock()
	defer fs.mu.Unlock()

	log.Printf("[%s] Cleaning directory: %s", fn, fs.directory)
	dir, err := os.Open(fs.directory)
	if err != nil {
		return fmt.Errorf("error opening directory: %w", err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("error reading directory names: %w", err)
	}

	for _, file := range files {
		if err := os.Remove(filepath.Join(fs.directory, file)); err != nil {
			return fmt.Errorf("error removing file %s: %w", file, err)
		}
	}
	return nil
}
