package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	// Test Config struct
	config := Config{
		Name:        "test-template",
		Description: "A test template",
		Version:     "1.0.0",
		Author:      "Test Author",
		Engine:      "pongo2",
		Variables: []Variable{
			{
				Name:     "project_name",
				Type:     "string",
				Prompt:   "Enter project name:",
				Default:  "my-project",
				Required: true,
			},
		},
	}

	if config.Name != "test-template" {
		t.Errorf("Config.Name = %v, want %v", config.Name, "test-template")
	}

	if config.Description != "A test template" {
		t.Errorf("Config.Description = %v, want %v", config.Description, "A test template")
	}

	if config.Version != "1.0.0" {
		t.Errorf("Config.Version = %v, want %v", config.Version, "1.0.0")
	}

	if config.Author != "Test Author" {
		t.Errorf("Config.Author = %v, want %v", config.Author, "Test Author")
	}

	if config.Engine != "pongo2" {
		t.Errorf("Config.Engine = %v, want %v", config.Engine, "pongo2")
	}

	if len(config.Variables) != 1 {
		t.Errorf("Config.Variables length = %v, want %v", len(config.Variables), 1)
	}
}

func TestVariable(t *testing.T) {
	// Test Variable struct
	variable := Variable{
		Name:     "project_name",
		Type:     "string",
		Prompt:   "Enter project name:",
		Default:  "my-project",
		Required: true,
		Choices:  []string{"option1", "option2"},
	}

	if variable.Name != "project_name" {
		t.Errorf("Variable.Name = %v, want %v", variable.Name, "project_name")
	}

	if variable.Type != "string" {
		t.Errorf("Variable.Type = %v, want %v", variable.Type, "string")
	}

	if variable.Prompt != "Enter project name:" {
		t.Errorf("Variable.Prompt = %v, want %v", variable.Prompt, "Enter project name:")
	}

	if variable.Default != "my-project" {
		t.Errorf("Variable.Default = %v, want %v", variable.Default, "my-project")
	}

	if !variable.Required {
		t.Error("Variable.Required should be true")
	}

	if len(variable.Choices) != 2 {
		t.Errorf("Variable.Choices length = %v, want %v", len(variable.Choices), 2)
	}
}

func TestLoadConfig_TOML(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "ason_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test TOML config
	tomlContent := `name = "test-template"
description = "A test template"
version = "1.0.0"
author = "Test Author"
engine = "pongo2"

[[variables]]
name = "project_name"
type = "string"
prompt = "Enter project name:"
default = "my-project"
required = true

[[variables]]
name = "use_docker"
type = "boolean"
prompt = "Use Docker?"
default = false
required = false
`

	tomlPath := filepath.Join(tmpDir, "template.toml")
	err = os.WriteFile(tomlPath, []byte(tomlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write TOML file: %v", err)
	}

	config, err := LoadConfig(tomlPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.Name != "test-template" {
		t.Errorf("Config.Name = %v, want %v", config.Name, "test-template")
	}

	if config.Description != "A test template" {
		t.Errorf("Config.Description = %v, want %v", config.Description, "A test template")
	}

	if config.Version != "1.0.0" {
		t.Errorf("Config.Version = %v, want %v", config.Version, "1.0.0")
	}

	if config.Author != "Test Author" {
		t.Errorf("Config.Author = %v, want %v", config.Author, "Test Author")
	}

	if config.Engine != "pongo2" {
		t.Errorf("Config.Engine = %v, want %v", config.Engine, "pongo2")
	}

	if len(config.Variables) != 2 {
		t.Errorf("Config.Variables length = %v, want %v", len(config.Variables), 2)
	}

	// Check first variable
	if config.Variables[0].Name != "project_name" {
		t.Errorf("Variables[0].Name = %v, want %v", config.Variables[0].Name, "project_name")
	}

	if config.Variables[0].Type != "string" {
		t.Errorf("Variables[0].Type = %v, want %v", config.Variables[0].Type, "string")
	}

	if !config.Variables[0].Required {
		t.Error("Variables[0].Required should be true")
	}

	// Check second variable
	if config.Variables[1].Name != "use_docker" {
		t.Errorf("Variables[1].Name = %v, want %v", config.Variables[1].Name, "use_docker")
	}

	if config.Variables[1].Type != "boolean" {
		t.Errorf("Variables[1].Type = %v, want %v", config.Variables[1].Type, "boolean")
	}

	if config.Variables[1].Required {
		t.Error("Variables[1].Required should be false")
	}
}

func TestLoadConfig_JSON(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "ason_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test JSON config
	jsonContent := `{
  "name": "json-template",
  "description": "A JSON test template",
  "version": "2.0.0",
  "author": "JSON Author",
  "engine": "pongo2",
  "variables": [
    {
      "name": "app_name",
      "type": "string",
      "prompt": "Enter app name:",
      "default": "my-app",
      "required": true,
      "choices": ["web", "api", "cli"]
    }
  ]
}`

	jsonPath := filepath.Join(tmpDir, "template.json")
	err = os.WriteFile(jsonPath, []byte(jsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	config, err := LoadConfig(jsonPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if config.Name != "json-template" {
		t.Errorf("Config.Name = %v, want %v", config.Name, "json-template")
	}

	if config.Description != "A JSON test template" {
		t.Errorf("Config.Description = %v, want %v", config.Description, "A JSON test template")
	}

	if config.Version != "2.0.0" {
		t.Errorf("Config.Version = %v, want %v", config.Version, "2.0.0")
	}

	if len(config.Variables) != 1 {
		t.Errorf("Config.Variables length = %v, want %v", len(config.Variables), 1)
	}

	// Check variable with choices
	if len(config.Variables[0].Choices) != 3 {
		t.Errorf("Variables[0].Choices length = %v, want %v", len(config.Variables[0].Choices), 3)
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	_, err := LoadConfig("/non/existent/file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoadConfig_InvalidFormat(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "ason_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test invalid format
	invalidContent := `this is not valid yaml or json`

	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	err = os.WriteFile(invalidPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid file: %v", err)
	}

	_, err = LoadConfig(invalidPath)
	if err == nil {
		t.Error("Expected error for invalid config format, got nil")
	}
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "ason_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test empty file
	emptyPath := filepath.Join(tmpDir, "empty.yaml")
	err = os.WriteFile(emptyPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}

	config, err := LoadConfig(emptyPath)
	if err != nil {
		t.Fatalf("LoadConfig() failed for empty file: %v", err)
	}

	// Should return zero-value config
	if config.Name != "" {
		t.Errorf("Expected empty name for empty config, got %v", config.Name)
	}
}
