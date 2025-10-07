package varfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_TOML_Simple(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create a simple TOML file
	tomlFile := filepath.Join(tempDir, "vars.toml")
	content := `
environment = "prod"
aws_region = "us-west-2"
organization = "acme"
`
	if err := os.WriteFile(tomlFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load variables
	vars, err := Load(tomlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify variables
	expected := map[string]string{
		"environment":  "prod",
		"aws_region":   "us-west-2",
		"organization": "acme",
	}

	if len(vars) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(vars))
	}

	for key, expectedValue := range expected {
		if actualValue, ok := vars[key]; !ok {
			t.Errorf("Missing variable: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Variable %s: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestLoad_TOML_TemplateFormat(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create a template-style TOML file with [variables] section
	tomlFile := filepath.Join(tempDir, "template.toml")
	content := `
[template]
name = "lambda-waf-ipset"
description = "AWS Lambda function"

[variables]
organization = { type = "string", description = "Organization name", default = "acme" }
aws_region = { type = "string", description = "AWS region", default = "us-east-1" }
environment = { type = "string", description = "Environment", default = "dev" }
`
	if err := os.WriteFile(tomlFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load variables
	vars, err := Load(tomlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify variables
	expected := map[string]string{
		"organization": "acme",
		"aws_region":   "us-east-1",
		"environment":  "dev",
	}

	if len(vars) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(vars))
	}

	for key, expectedValue := range expected {
		if actualValue, ok := vars[key]; !ok {
			t.Errorf("Missing variable: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Variable %s: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestLoad_YAML(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create a YAML file
	yamlFile := filepath.Join(tempDir, "vars.yaml")
	content := `
environment: prod
aws_region: us-west-2
organization: acme
`
	if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load variables
	vars, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify variables
	expected := map[string]string{
		"environment":  "prod",
		"aws_region":   "us-west-2",
		"organization": "acme",
	}

	if len(vars) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(vars))
	}

	for key, expectedValue := range expected {
		if actualValue, ok := vars[key]; !ok {
			t.Errorf("Missing variable: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Variable %s: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestLoad_JSON(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	// Create a JSON file
	jsonFile := filepath.Join(tempDir, "vars.json")
	content := `{
  "environment": "prod",
  "aws_region": "us-west-2",
  "organization": "acme"
}`
	if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load variables
	vars, err := Load(jsonFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify variables
	expected := map[string]string{
		"environment":  "prod",
		"aws_region":   "us-west-2",
		"organization": "acme",
	}

	if len(vars) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(vars))
	}

	for key, expectedValue := range expected {
		if actualValue, ok := vars[key]; !ok {
			t.Errorf("Missing variable: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Variable %s: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/file.toml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestLoad_UnsupportedFormat(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create a file with unsupported extension
	txtFile := filepath.Join(tempDir, "vars.txt")
	if err := os.WriteFile(txtFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := Load(txtFile)
	if err == nil {
		t.Error("Expected error for unsupported file format, got nil")
	}
}

func TestMerge(t *testing.T) {
	fileVars := map[string]string{
		"environment":  "dev",
		"aws_region":   "us-east-1",
		"organization": "acme",
	}

	cliVars := map[string]string{
		"environment": "prod",
		"new_var":     "value",
	}

	merged := Merge(fileVars, cliVars)

	// Verify merged results
	expected := map[string]string{
		"environment":  "prod",      // CLI override
		"aws_region":   "us-east-1", // From file
		"organization": "acme",      // From file
		"new_var":      "value",     // From CLI
	}

	if len(merged) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(merged))
	}

	for key, expectedValue := range expected {
		if actualValue, ok := merged[key]; !ok {
			t.Errorf("Missing variable: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Variable %s: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	// Test with empty file vars
	result := Merge(nil, map[string]string{"key": "value"})
	if len(result) != 1 || result["key"] != "value" {
		t.Error("Merge with empty file vars failed")
	}

	// Test with empty CLI vars
	result = Merge(map[string]string{"key": "value"}, nil)
	if len(result) != 1 || result["key"] != "value" {
		t.Error("Merge with empty CLI vars failed")
	}

	// Test with both empty
	result = Merge(nil, nil)
	if len(result) != 0 {
		t.Error("Merge with both empty should return empty map")
	}
}

func TestLoad_YAML_WithVariablesSection(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create a YAML file with variables section
	yamlFile := filepath.Join(tempDir, "vars.yaml")
	content := `
variables:
  environment: prod
  aws_region: us-west-2
  organization: acme
`
	if err := os.WriteFile(yamlFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load variables
	vars, err := Load(yamlFile)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify variables
	expected := map[string]string{
		"environment":  "prod",
		"aws_region":   "us-west-2",
		"organization": "acme",
	}

	if len(vars) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(vars))
	}

	for key, expectedValue := range expected {
		if actualValue, ok := vars[key]; !ok {
			t.Errorf("Missing variable: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Variable %s: expected %q, got %q", key, expectedValue, actualValue)
		}
	}
}
