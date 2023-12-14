package dedup

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
)

func generateRandomString(length int) (string, error) {
	const fn = "internal/service/dedup/save/generateRandomString"

	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	log.Printf("[%s] Generating a random string of length %d", fn, length)

	result := make([]byte, length)
	for i := range result {
		if randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters)))); err != nil {
			return "", fmt.Errorf("error generating random number: %w", err)
		} else {
			result[i] = letters[randIndex.Int64()]
		}
	}

	return string(result), nil
}

func (svc *Svc) Save(path string) (string, error) {
	const fn = "internal/service/dedup/save/Save"

	storedFileName, err := generateRandomString(5)
	if err != nil {
		return "", fmt.Errorf("error generating random marker: %w", err)
	}

	svc.occurrencesStorage.SetWD(storedFileName)

	defer func() {
		if err != nil {
			os.RemoveAll(svc.occurrencesStorage.GetDirectory())
		}
	}()

	if err = svc.saveFile(path); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	log.Printf("[%s] File saved successfully with marker: %s", fn, storedFileName)
	return storedFileName, nil
}

func (svc *Svc) saveFile(sourceFilePath string) error {
	const fn = "internal/service/dedup/save/saveFile"

	f, err := os.Open(sourceFilePath)
	if err != nil {
		return fmt.Errorf("error opening source file '%s': %w", sourceFilePath, err)
	}
	defer f.Close()

	if segments, lastBatchSize, err := svc.processFile(f); err != nil {
		return fmt.Errorf("error processing file: %w", err)
	} else if err = svc.occurrencesStorage.PutMetadata(sourceFilePath, segments, lastBatchSize); err != nil {
		return fmt.Errorf("error writing info file: %w", err)
	}

	return nil
}

func (svc *Svc) processFile(f *os.File) (int, int, error) {
	const fn = "internal/service/dedup/save/processFile"

	stat, err := f.Stat()
	if err != nil {
		return 0, 0, err
	}

	lastBatchSize := int(stat.Size() % int64(svc.batchSize))
	if lastBatchSize == 0 && stat.Size() > 0 {
		lastBatchSize = svc.batchSize
	}

	segments := int(stat.Size() / int64(svc.batchSize))

	hashFuncGenerator := getHashFunc(svc.hashFunc)
	buf := make([]byte, svc.batchSize)

	for i := segments; i > 0; i-- {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, fmt.Errorf("error reading file: %w", err)
		}

		hashFunc := hashFuncGenerator()
		hash := hashFunc.Sum(buf[:n])[:16]

		if err = svc.batchStorage.Put(hash, buf[:n]); err != nil && err.Error() != "batch already exists" {
			return 0, 0, fmt.Errorf("error storing batch: %w", err)
		}

		if err = svc.occurrencesStorage.Put(hash, segments); err != nil {
			return 0, 0, fmt.Errorf("error storing data: %w", err)
		}

		segments++
	}

	startPos := stat.Size() - int64(lastBatchSize)

	_, err = f.Seek(startPos, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	buf = make([]byte, stat.Size()-startPos)
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return 0, 0, err
	}

	hashFunc := hashFuncGenerator()
	hash := hashFunc.Sum(buf[:lastBatchSize])[:16]

	if err = svc.batchStorage.Put(hash, buf[:lastBatchSize]); err != nil && err.Error() != "batch already exists" {
		return 0, 0, fmt.Errorf("error storing batch: %w", err)
	}

	if err = svc.occurrencesStorage.Put(hash, segments); err != nil {
		return 0, 0, fmt.Errorf("error storing data: %w", err)
	}

	return segments, lastBatchSize, nil
}
