package dedup

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func (svc *Svc) Restore(marker string) (err error) {
	const fn = "internal/service/dedup/service/Restore"

	log.Printf("[%s] Restoring data using marker: %s", fn, marker)

	infoFilePath := fmt.Sprintf("%s/%s/info", svc.fsStorage.GetDirectory(), marker)
	info, err := os.Open(infoFilePath)
	if err != nil {
		log.Fatalf("[%s] Error opening info file: %s", fn, err)
	}
	defer info.Close()

	fileName, segmentsNum, err := parseInfoFile(info)
	if err != nil {
		log.Fatalf("[%s] Error parsing info file: %s", fn, err)
	}

	log.Printf("[%s] Restoring to file: %s with %d segments", fn, fileName, segmentsNum)

	resFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("[%s] Error creating result file: %s", fn, err)
	}
	defer resFile.Close()

	fileSize := int64(segmentsNum) * int64(svc.batchSize)
	err = setFileSize(resFile, fileSize)
	if err != nil {
		log.Fatalf("[%s] Error setting result file size: %s", fn, err)
	}

	log.Printf("[%s] Result file size set to %d bytes", fn, fileSize)

	err = rewriteSegments(resFile, marker, svc)
	if err != nil {
		log.Fatalf("[%s] Error rewriting segments: %s", fn, err)
	}

	log.Printf("[%s] Data restoration completed successfully", fn)

	return nil
}

func parseInfoFile(info *os.File) (string, int, error) {
	const fn = "internal/service/dedup/service/Restore/parseInfoFile"

	log.Printf("[%s] Parsing info file", fn)

	data := make([]byte, 1024) // Adjust buffer size as needed
	_, err := info.Read(data)
	if err != nil {
		return "", 0, err
	}

	parts := strings.Split(string(data), "\n")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("[%s] Info File format is incorrect", fn)
	}

	fileName := parts[0]
	segmentsNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("[%s] Error parsing number: %s", fn, err)
	}

	log.Printf("[%s] Parsed info: File: %s, Segments: %d", fn, fileName, segmentsNum)

	return fileName, segmentsNum, nil
}

func setFileSize(file *os.File, size int64) error {
	const fn = "internal/service/dedup/service/Restore/setFileSize"

	log.Printf("[%s] Setting result file size to %d bytes", fn, size)

	return file.Truncate(size)
}

func rewriteSegments(resFile *os.File, marker string, svc *Svc) error {
	const fn = "internal/service/dedup/service/Restore/rewriteSegments"

	log.Printf("[%s] Rewriting segments using marker: %s", fn, marker)

	dir, err := os.Open(marker)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		hash := []byte(fileInfo.Name())

		log.Printf("[%s] Restoring segment for hash: %x", fn, hash)

		segments, err := svc.fsStorage.Get(hash)
		if err != nil {
			return err
		}

		batch, err := svc.batchStorage.Get(hash)
		if err != nil {
			return err
		}

		err = writeSegments(resFile, batch, segments, svc.batchSize)
		if err != nil {
			return err
		}
	}

	log.Printf("[%s] Segments rewriting completed", fn)

	return nil
}

func writeSegments(file *os.File, data []byte, segments []int, batchSize int) error {
	const fn = "internal/service/dedup/service/Restore/writeSegments"

	for _, segment := range segments {
		offset := int64(batchSize * segment)
		_, err := file.Seek(offset, 0)
		if err != nil {
			return err
		}

		_, err = file.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
