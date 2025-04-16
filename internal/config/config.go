package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonData, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err = write(jsonData); err != nil {
		return err
	}
	return nil
}

func write(jsonData []byte) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(homeDir, configFileName)
	return fullPath, nil
}
