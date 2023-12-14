package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	Alg              string `json:"alg"`
	BatchSize        int    `json:"batchSize"`
	BatchStoragePath string `json:"batchStoragePath"`
	RestorationPath  string `json:"restorationPath"`
}

func LoadConfig(cfgFileName string) *Config {
	const fn = "config/config/LoadConfig"

	cfgPath := fmt.Sprintf(cfgFileName)
	log.Printf("Trying to load configuration file at %s", cfgPath)
	f, err := os.Open(cfgPath)
	if err != nil {
		log.Printf("[%s] %s", fn, err)
		log.Fatalf("[%s] Cannot open configuration file", fn)
	}
	defer f.Close()

	var cfgFile []byte
	cfgFile, err = io.ReadAll(f) // Read the entire file into cfgFile
	if err != nil {
		log.Fatalf("[%s] Error reading configuration file: %s", fn, err)
	}

	var cfg Config
	err = json.Unmarshal(cfgFile, &cfg)
	if err != nil {
		log.Fatalf("[%s] Configuration parsing failed: %s", fn, err)
	}

	return &cfg
}
