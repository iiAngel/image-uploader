package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type BackendConfig struct {
	MaxUploadSize int    `toml:"max_upload_size" comment:"The maximum size of a file that can be uploaded in Megabytes"`
	FrontendPath  string `toml:"frontend_path"`
	FilesPath     string `toml:"files_path" comment:"Absolute path where the files uploaded need to be stored"`
	Port          uint16 `toml:"port" comment:"Port to listen to"`
}

var (
	DefaultConfig = BackendConfig{
		MaxUploadSize: 2,
		FrontendPath:  "./frontend",
		FilesPath:     "",
		Port:          7445,
	}
	LoadedConfig BackendConfig
)

func CreateConfig(filename string, info *BackendConfig) error {
	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	marshalledToml, err := toml.Marshal(info)
	if err != nil {
		return err
	}

	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(marshalledToml)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig(filename string) (BackendConfig, error) {
	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return BackendConfig{}, err
	}

	file, err := os.Open(fullpath)
	if err != nil {
		return BackendConfig{}, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	infoConfig := BackendConfig{}

	if err := toml.Unmarshal(data, &infoConfig); err != nil {
		return BackendConfig{}, err
	}

	return infoConfig, nil
}

func TryLoadConfig() {
	loadedConfig, err := LoadConfig("./config.toml")

	if err != nil {
		log.Println("Config file not found, creating default config file...")
		err = CreateConfig("./config.toml", &DefaultConfig)

		if err != nil {
			log.Panic(err)
		}

		loadedConfig = DefaultConfig
	}

	log.Println("Loaded config file!")

	LoadedConfig = loadedConfig
}
