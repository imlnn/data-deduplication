package main

import (
	"dedup/config"
	"dedup/internal/service/dedup"
	"dedup/internal/storage/mem"
	"fmt"
	"log"
	"os"
)

// Commands format:
// ./dedup save %FILE_PATH%
// ./dedup restore %FILE_MARKER%

func main() {
	cfg := config.LoadConfig("config.json")
	strg := mem.NewStorage()
	svc, err := dedup.NewSvc(cfg, strg)
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
		svc.Save(filePath)

	case "restore":
		if len(args) < 2 {
			fmt.Println("Error: File marker not provided")
			return
		}
		fileMarker := args[1]
		svc.Restore(fileMarker)
		if err != nil {

		}

	default:
		fmt.Println("Error: Unknown command")
	}
}
