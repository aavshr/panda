package config

import (
	"encoding/json"
	"errors"
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
)

const (
	appConfigDir   = "aavshr-panda" // to not have name conflicts with other apps
	configFileName = "config.json"
	defaultModel   = "gpt-4o-mini"
)

var (
	ErrConfigNotFound = errors.New("config file not found")
)

type Config struct {
	LLMAPIKey string `json:"llm_api_key"`
	LLMModel  string `json:"llm_model"`
}

func GetDir() string {
	configDir := xdg.ConfigHome
	if configDir == "" {
		configDir = filepath.Join(xdg.Home, ".config")
	}
	return filepath.Join(configDir, appConfigDir)
}

func GetDataDir() string {
	dataDir := xdg.DataHome
	if dataDir == "" {
		dataDir = filepath.Join(xdg.Home, ".local", "share")
	}
	return filepath.Join(dataDir, appConfigDir)
}

func GetFilePath() string {
	configDir := GetDir()
	return filepath.Join(configDir, configFileName)
}

func Load() (*Config, error) {
	configFilePath := GetFilePath()
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, ErrConfigNotFound
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	config := &Config{}
	if err := json.NewDecoder(configFile).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func Save(config Config) (*Config, error) {
	if config.LLMModel == "" {
		config.LLMModel = defaultModel
	}
	configDir := GetDir()
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return &config, err
	}

	configFilePath := GetFilePath()
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return &config, err
	}
	defer configFile.Close()

	return &config, json.NewEncoder(configFile).Encode(config)
}
