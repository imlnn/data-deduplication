package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Alg       string `json:"alg"`
	BatchSize int    `json:"batchSize"`
}

func LoadConfig(cfgPath string) *Config {
	const fn = "config/config/LoadConfig"

	log.Printf("Trying load configuration file at %s", cfgPath)
	f, err := os.Open("config.json")
	if err != nil {
		log.Printf("[%s] %s", fn, err)
		log.Fatalf("[%s] Cannot open configuration file", fn)
	}
	defer f.Close()

	var cfgFile []byte
	_, err = f.Read(cfgFile)

	var cfg Config
	err = json.Unmarshal(cfgFile, &cfg)
	if err != nil {
		log.Fatal("Configuration parsing failed")
	}

	return &cfg
}
