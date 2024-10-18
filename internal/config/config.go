package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetConfigPath returns the path to the configuration file.
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	telConfigDir := filepath.Join(configDir, "tel")
	if _, err := os.Stat(telConfigDir); os.IsNotExist(err) {
		err := os.MkdirAll(telConfigDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(telConfigDir, "config.json"), nil
}

// ReadConfig reads the configuration file and returns it as a Config struct.
func ReadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// WriteConfig writes the Config struct to the configuration file.
func WriteConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("bot_token is empty")
	}
	if c.ChatID == "" {
		return fmt.Errorf("chat_id is empty")
	}
	return nil
}
