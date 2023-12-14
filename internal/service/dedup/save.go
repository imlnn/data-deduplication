package dedup

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
)

func generateRandomString(length int) (string, error) {
	const fn = "internal/service/dedup/service/save/generateRandomString"

	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Log that the function is starting
	log.Printf("[%s] Generating a random string of length %d", fn, length)

	result := make([]byte, length)

	for i := range result {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// Log the error and return it
			log.Printf("[%s] Error generating random number: %s", fn, err)
			return "", err
		}
		result[i] = letters[randIndex.Int64()]
	}

	randomStr := string(result)

	// Log that the function has completed successfully
	log.Printf("[%s] Generated random string: %s", fn, randomStr)

	return randomStr, nil
}

func (svc *Svc) Save(path string) (marker string, err error) {
	const fn = "internal/service/dedup/service/Save"

	// Generate a random marker for the saved file
	storedFileName, err := generateRandomString(5)
	if err != nil {
		log.Printf("[%s] Error generating random marker: %s", fn, err)
		return "", err
	}

	// Create a directory with the generated marker
	wd, err := os.Getwd()
	dir := fmt.Sprintf("%s\\%s", wd, storedFileName)
	err = os.Mkdir(dir, 0777)
	if err != nil {
		log.Printf("[%s] Error creating directory: %s", fn, err)
		return "", err
	}

	// Based on the selected hash function, perform the save operation
	switch svc.hashFunc {
	case sha1:
		log.Printf("[%s] Saving file using SHA-1 hash function...", fn)
		err = svc.saveSHA1(path, storedFileName)
		if err != nil {
			// Log and clean up on error
			log.Printf("[%s] Error saving file: %s", fn, err)
			_ = os.RemoveAll(storedFileName)
			return "", err
		}
	default:
		log.Printf("[%s] Unsupported hash function selected: %v", fn, svc.hashFunc)
		return "", fmt.Errorf("unsupported hash function")
	}

	log.Printf("[%s] File saved successfully with marker: %s", fn, storedFileName)

	return storedFileName, nil
}
