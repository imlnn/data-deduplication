package main

import (
	"data-deduplication/config"
	"data-deduplication/internal/service/dedup"
	"data-deduplication/internal/storage/fsstorage"
	"fmt"
	"log"
	"os"
	"time"
)

// Commands format:
// ./dedup save %FILE_PATH%
// ./dedup restore %FILE_MARKER%
func main() {
	const fn = "cmd/dedup/main"

	// Load the configuration
	cfg := config.LoadConfig("config.json")

	// Create storage instances
	batchStorage := fsstorage.NewFSBatchStorage(cfg.BatchStoragePath)
	fsOccurrencesStorage := fsstorage.NewFSStorage("occurrences")
	// Create the deduplication service
	svc, err := dedup.NewSvc(cfg, batchStorage, fsOccurrencesStorage)
	if err != nil {
		log.Fatalf("[%s] Error creating deduplication service: %s", fn, err)
	}

	// ==================== DEBUG SECTION ====================

	startTime := time.Now()

	file, err := svc.Save("testVIDEO.mp4")
	if err != nil {
		log.Printf("[%s] Error saving file: %s", fn, err)
	} else {
		log.Printf("[%s] File saved: %s", fn, file)
	}
	saveTime := time.Since(startTime)

	startTime = time.Now()
	fileMarker := file
	err = svc.Restore(fileMarker)
	if err != nil {
		log.Printf("[%s] Error restoring file: %s", fn, err)
	} else {
		log.Printf("[%s] File %s restoration completed", fn, fileMarker)
	}
	restoreTime := time.Since(startTime)

	fmt.Printf("Save time: %s, Restore time: %s", saveTime, restoreTime)

	// ==================== DEBUG SECTION ====================

	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatalf("[%s] Error: Command not provided", fn)
	}

	switch args[0] {
	case "save":
		if len(args) < 2 {
			log.Fatalf("[%s] Error: File path not provided", fn)
		}
		filePath := args[1]
		file, err := svc.Save(filePath)
		if err != nil {
			log.Printf("[%s] Error saving file: %s", fn, err)
		} else {
			log.Printf("[%s] File saved: %s", fn, file)
		}

	case "restore":
		if len(args) < 2 {
			log.Fatalf("[%s] Error: File marker not provided", fn)
		}
		fileMarker := args[1]
		err := svc.Restore(fileMarker)
		if err != nil {
			log.Printf("[%s] Error restoring file: %s", fn, err)
		} else {
			log.Printf("[%s] File %s restoration completed", fn, fileMarker)
		}

	default:
		log.Fatalf("[%s] Error: Unknown command", fn)
	}
}
