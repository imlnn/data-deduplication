package dedup

import (
	SHA1 "crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func (svc *Svc) saveSHA1(path string, savedDirectory string) error {
	const fn = "internal/service/dedup/service/save/saveSHA1"

	log.Printf("[%s] Saving file '%s' to directory '%s'", fn, path, savedDirectory)

	// Open the source file
	f, err := os.Open(path)
	if err != nil {
		log.Printf("[%s] Error opening source file: %s", fn, err)
		return err
	}
	defer f.Close()

	// Create an info file in the saved directory
	infoFilePath := fmt.Sprintf("%s/info", savedDirectory)
	info, err := os.Create(infoFilePath)
	if err != nil {
		log.Printf("[%s] Error creating info file: %s", fn, err)
		return err
	}
	defer info.Close()

	segments, err := svc.processFile(f, savedDirectory)
	if err != nil {
		log.Printf("[%s] Error processing file: %s", fn, err)
		return err
	}

	err = svc.writeInfoFile(info, path, segments)
	if err != nil {
		log.Printf("[%s] Error writing info file: %s", fn, err)
		return err
	}

	log.Printf("[%s] File saved successfully with %d segments", fn, segments)
	return nil
}

func (svc *Svc) processFile(f *os.File, savedDirectory string) (int, error) {
	const fn = "internal/service/dedup/service/save/processFile"

	var segments = 0

	buf := make([]byte, svc.batchSize)
	for i := 0; ; i++ {
		n, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("[%s] Error reading file: %s", fn, err)
			}
			return segments, err
		} else if n == 0 {
			break
		}

		hashFunc := SHA1.New()
		hash := hashFunc.Sum(buf)

		err = svc.batchStorage.Put(hash, buf)
		if err != nil && err.Error() != "batch already exists" {
			log.Printf("[%s] Error storing batch: %s", fn, err)
			return segments, err
		}

		err = svc.fsStorage.Put(hash, i)
		if err != nil {
			log.Printf("[%s] Error storing data: %s", fn, err)
			return segments, err
		}

		segments++
	}

	log.Printf("[%s] Processed file successfully with %d segments", fn, segments)
	return segments, nil
}

func (svc *Svc) writeInfoFile(info *os.File, path string, segments int) error {
	const fn = "internal/service/dedup/service/save/writeInfoFile"

	infoStr := fmt.Sprintf("%s\n%d", filepath.Base(path), segments)
	_, err := info.WriteString(infoStr)
	if err != nil {
		log.Printf("[%s] Error writing info file: %s", fn, err)
		return err
	}

	log.Printf("[%s] Info file written successfully", fn)
	return nil
}
