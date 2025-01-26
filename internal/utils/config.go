package utils

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlStruct struct {
	Config Config `yaml:"config"`
}

type Config struct {
	SecretAccess      string `yaml:"secret_access"`
	KeyAccessMetadata string `yaml:"key_access_metadata"`
	Port              int    `yaml:"port"`
	StoragePath       string `yaml:"storage_path"`
	UI                bool   `yaml:"ui"`
	DBName            string `yaml:"db_name"`
	MigrationPath     string `yaml:"migration_path"`
}

func OpenYAML() (*YamlStruct, error) {
	var pathFile string
	if _, err := os.Stat("config.yaml"); err == nil {
		pathFile = "config.yaml"
	} else if _, err := os.Stat("config.yml"); err == nil {
		pathFile = "config.yml"
	} else {
		return nil, errors.New("config yaml file not found")
	}

	data, err := os.ReadFile(pathFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var config YamlStruct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &config, nil
}
