package varfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Load reads variables from a file and returns them as a map.
// Supports TOML, YAML, and JSON formats based on file extension.
// For TOML files, it supports both simple key-value format and the template format with [variables] section.
func Load(filePath string) (map[string]string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("variable file not found: %s", filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read variable file: %w", err)
	}

	// Determine format by extension
	ext := strings.ToLower(filepath.Ext(filePath))

	var variables map[string]string
	switch ext {
	case ".toml":
		variables, err = loadTOML(content)
	case ".yaml", ".yml":
		variables, err = loadYAML(content)
	case ".json":
		variables, err = loadJSON(content)
	default:
		return nil, fmt.Errorf("unsupported file format: %s (supported: .toml, .yaml, .yml, .json)", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse %s file: %w", ext, err)
	}

	return variables, nil
}

// loadTOML parses a TOML file and extracts variables.
// Supports both simple key-value format and template format with [variables] section.
func loadTOML(content []byte) (map[string]string, error) {
	// First try to parse as a template-style TOML with [variables] section
	var templateFormat struct {
		Variables map[string]interface{} `toml:"variables"`
	}

	if err := toml.Unmarshal(content, &templateFormat); err == nil && len(templateFormat.Variables) > 0 {
		// Extract default values or direct string values from variables
		variables := make(map[string]string)
		for key, value := range templateFormat.Variables {
			switch v := value.(type) {
			case string:
				// Direct string value
				variables[key] = v
			case map[string]interface{}:
				// Variable definition with default value
				if defaultVal, ok := v["default"]; ok {
					variables[key] = fmt.Sprintf("%v", defaultVal)
				}
			}
		}
		if len(variables) > 0 {
			return variables, nil
		}
	}

	// Try simple key-value format
	var simpleFormat map[string]interface{}
	if err := toml.Unmarshal(content, &simpleFormat); err != nil {
		return nil, err
	}

	// Convert all values to strings
	variables := make(map[string]string)
	for key, value := range simpleFormat {
		// Skip special sections like [template] or [variables]
		if key == "template" || key == "variables" {
			continue
		}
		variables[key] = fmt.Sprintf("%v", value)
	}

	return variables, nil
}

// loadYAML parses a YAML file and extracts variables.
func loadYAML(content []byte) (map[string]string, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	// Check if there's a variables section
	if vars, ok := data["variables"].(map[string]interface{}); ok {
		return convertToStringMap(vars), nil
	}

	// Otherwise use the entire document
	return convertToStringMap(data), nil
}

// loadJSON parses a JSON file and extracts variables.
func loadJSON(content []byte) (map[string]string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	// Check if there's a variables section
	if vars, ok := data["variables"].(map[string]interface{}); ok {
		return convertToStringMap(vars), nil
	}

	// Otherwise use the entire document
	return convertToStringMap(data), nil
}

// convertToStringMap converts a map[string]interface{} to map[string]string.
func convertToStringMap(data map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range data {
		// Handle nested maps (variable definitions with default values)
		if m, ok := value.(map[string]interface{}); ok {
			if defaultVal, exists := m["default"]; exists {
				result[key] = fmt.Sprintf("%v", defaultVal)
				continue
			}
		}
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

// Merge combines variables from a file with command-line variables.
// Command-line variables take precedence over file variables.
func Merge(fileVars, cliVars map[string]string) map[string]string {
	result := make(map[string]string)

	// Start with file variables
	for key, value := range fileVars {
		result[key] = value
	}

	// Override with CLI variables
	for key, value := range cliVars {
		result[key] = value
	}

	return result
}
