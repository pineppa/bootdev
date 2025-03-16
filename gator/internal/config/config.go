package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFilePath = "/.gatorconfig.json"

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homePath + configFilePath, nil
}

func Read() Config {
	// Retrieve the correct path location
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Error compiling the file path:", err)
		return Config{}
	}

	// Open the JSON jsonFile
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return Config{}
	}
	defer jsonFile.Close()

	// Read JSON into bytes
	body, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return Config{}
	}

	// Unmarshal content from bytes to a config struct
	var config Config
	json.Unmarshal(body, &config)
	return config
}

func (c Config) SetUser() error {
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Error compiling the file path:", err)
		return err
	}

	// Open the JSON jsonFile
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer jsonFile.Close()

	// Unmarshal content from bytes to a config struct
	body, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Read JSON into bytes
	return os.WriteFile(path, body, os.FileMode(0644))
}
