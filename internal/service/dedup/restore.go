package dedup

import (
	"log"
	"os"
)

func (svc *Svc) Restore(marker string) (err error) {
	const fn = "internal/service/dedup/service/Restore"

	log.Printf("[%s] Restoring data using marker: %s", fn, marker)

	fileName, segmentsNum, lastBatchSize, err := svc.occurrencesStorage.GetMetadata()
	if err != nil {
		log.Fatalf("[%s] Error parsing info file: %s", fn, err)
	}

	log.Printf("[%s] Restoring to file: %s with %d segments", fn, fileName, segmentsNum)

	resFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("[%s] Error creating result file: %s", fn, err)
	}
	defer resFile.Close()

	fileSize := int64(segmentsNum)*int64(svc.batchSize) + int64(lastBatchSize)
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

func setFileSize(file *os.File, size int64) error {
	const fn = "internal/service/dedup/service/Restore/setFileSize"

	log.Printf("[%s] Setting result file size to %d bytes", fn, size)

	return file.Truncate(size)
}

func rewriteSegments(resFile *os.File, marker string, svc *Svc) error {
	const fn = "internal/service/dedup/service/Restore/rewriteSegments"

	log.Printf("[%s] Rewriting segments using marker: %s", fn, marker)

	dir, err := os.Open(svc.occurrencesStorage.GetDirectory())
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		if fileInfo.Name() == "info" {
			continue
		}

		hash := fileInfo.Name()

		log.Printf("[%s] Restoring segment for hash: %x", fn, hash)

		segments, err := svc.occurrencesStorage.Get(hash)
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
