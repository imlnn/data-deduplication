package dedup

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func (svc *Svc) Restore(marker string) (err error) {
	info, err := os.Open(fmt.Sprintf("%s/%s/info", svc.fsStorage.GetDirectory(), marker))
	var data []byte
	_, err = info.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(string(data), "\n")
	if len(parts) != 2 {
		return fmt.Errorf("Info File format is incorrect")
	}

	fileName := parts[0]
	segmentsNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("Error parsing number:", err)
	}

	resFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer resFile.Close()

	// Step 2: Set the file  size
	var fileSize int64
	fileSize = int64(segmentsNum) * int64(svc.batchSize) // Size in bytes
	err = resFile.Truncate(fileSize)
	if err != nil {
		log.Fatal(err)
	}

	// Step 3: Rewrite a portion of the file at a selected offset
	dir, err := os.Open(marker)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(-1) // -1 to read all files and directories
	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range files {
		hash := []byte(fileInfo.Name())
		segments, err := svc.fsStorage.Get(hash)
		if err != nil {
			return err
		}

		batch, err := svc.batchStorage.Get(hash)
		if err != nil {
			return err
		}

		for _, segment := range segments {
			offset := int64(svc.batchSize * segment) // For example, start rewriting at the 500th byte
			_, err = resFile.Seek(offset, 0)         // Seek to the offset
			if err != nil {
				return err
			}

			_, err = resFile.Write(batch) // Write the new data
			if err != nil {
				return err
			}
		}
	}
	return nil
}
