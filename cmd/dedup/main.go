package main

import (
	"dedup/config"
	"dedup/internal/service/dedup"
	"dedup/internal/storage/fsstorage"
	"fmt"
	"log"
	"os"
)

// Commands format:
// ./dedup save %FILE_PATH%
// ./dedup restore %FILE_MARKER%
func main() {
	cfg := config.LoadConfig("config.json")
	batchStorage := fsstorage.NewFSBatchStorage("")
	occurrencesStorage := fsstorage.NewFSStorage("")
	svc, err := dedup.NewSvc(cfg, batchStorage, occurrencesStorage)
	if err != nil {
		log.Fatal()
	}

	args := os.Args[1:]

	switch args[0] {
	case "save":
		if len(args) < 2 {
			fmt.Println("Error: File path not provided")
			return
		}
		filePath := args[1]
		file, err := svc.Save(filePath)
		if err != nil {
			fmt.Printf("Error: %s", err)
		} else {
			fmt.Printf("%s", file)
		}

	case "restore":
		if len(args) < 2 {
			fmt.Println("Error: File marker not provided")
			return
		}
		fileMarker := args[1]
		svc.Restore(fileMarker)
		if err != nil {
			fmt.Printf("Error: %s", err)
		} else {
			fmt.Printf("File %s restoration completed ", fileMarker)
		}

	default:
		fmt.Println("Error: Unknown command")
	}
}
