package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads and parses a YAML configuration file into a struct of type T
func LoadConfig[T any](filePath string) (T, error) {
	var config T

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}
