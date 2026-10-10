package config

import (
	"os"
	"path/filepath"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configContent, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer configContent.Close()

	var configData Config
	decoder := json.NewDecoder(configContent)
	err = decoder.Decode(&configData)
	if err != nil {
		return Config{}, err
	}

	return configData, nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	err := write(*c)
	return err
}

func write(cfg Config) error {
    configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(home, configFileName)
	return configPath, nil
}