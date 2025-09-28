package template

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the template configuration
type Config struct {
	Name        string     `toml:"name" json:"name"`
	Description string     `toml:"description" json:"description"`
	Version     string     `toml:"version" json:"version"`
	Author      string     `toml:"author" json:"author"`
	Engine      string     `toml:"engine" json:"engine"`
	Variables   []Variable `toml:"variables" json:"variables"`
}

// Variable represents a template variable
type Variable struct {
	Name     string      `toml:"name" json:"name"`
	Type     string      `toml:"type" json:"type"`
	Prompt   string      `toml:"prompt" json:"prompt"`
	Default  interface{} `toml:"default,omitempty" json:"default,omitempty"`
	Required bool        `toml:"required,omitempty" json:"required,omitempty"`
	Choices  []string    `toml:"choices,omitempty" json:"choices,omitempty"`
}

// LoadConfig loads template configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config

	// Try TOML first
	if err := toml.Unmarshal(data, &config); err == nil {
		return &config, nil
	}

	// Try JSON as fallback
	if err := json.Unmarshal(data, &config); err == nil {
		return &config, nil
	}

	return nil, fmt.Errorf("failed to parse config file")
}
