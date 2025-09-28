package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListCmd(t *testing.T) {
	// Test list command properties
	if listCmd == nil {
		t.Fatal("listCmd should not be nil")
	}

	if listCmd.Use != "list" {
		t.Errorf("listCmd.Use = %v, want %v", listCmd.Use, "list")
	}

	if listCmd.Short != "List available templates" {
		t.Errorf("listCmd.Short = %v, want %v", listCmd.Short, "List available templates")
	}
}

func TestListCmdExecution(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_list_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Capture output
	var buf bytes.Buffer
	listCmd.SetOut(&buf)

	// Execute list command
	err = listCmd.RunE(listCmd, []string{})
	if err != nil {
		t.Fatalf("listCmd execution failed: %v", err)
	}

	// Test passed if no error occurred
	// Reset
	listCmd.SetOut(nil)
}

func TestAddCmd(t *testing.T) {
	// Test add command properties
	if addCmd == nil {
		t.Fatal("addCmd should not be nil")
	}

	if addCmd.Use != "add [name] [path]" {
		t.Errorf("addCmd.Use = %v, want %v", addCmd.Use, "add [name] [path]")
	}

	if addCmd.Short != "Add a template to the registry" {
		t.Errorf("addCmd.Short = %v, want %v", addCmd.Short, "Add a template to the registry")
	}
}

func TestAddCmdExecution(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_add_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Create a test template directory
	testTemplateDir, err := os.MkdirTemp("", "test_template")
	if err != nil {
		t.Fatalf("Failed to create test template dir: %v", err)
	}
	defer os.RemoveAll(testTemplateDir)

	// Add some files to the template
	err = os.WriteFile(filepath.Join(testTemplateDir, "README.md"), []byte("# {{ project_name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Capture output
	var buf bytes.Buffer
	addCmd.SetOut(&buf)

	// Execute add command with valid template path
	err = addCmd.RunE(addCmd, []string{"test-template", testTemplateDir})
	if err != nil {
		t.Fatalf("addCmd execution failed: %v", err)
	}

	// Test passed if no error occurred
	// Reset
	addCmd.SetOut(nil)
}

func TestRemoveCmd(t *testing.T) {
	// Test remove command properties
	if removeCmd == nil {
		t.Fatal("removeCmd should not be nil")
	}

	if removeCmd.Use != "remove [name]" {
		t.Errorf("removeCmd.Use = %v, want %v", removeCmd.Use, "remove [name]")
	}

	if removeCmd.Short != "Remove a template from the registry" {
		t.Errorf("removeCmd.Short = %v, want %v", removeCmd.Short, "Remove a template from the registry")
	}

	// Check aliases
	expectedAliases := []string{"rm", "delete"}
	if len(removeCmd.Aliases) != len(expectedAliases) {
		t.Errorf("removeCmd should have %d aliases, got %d", len(expectedAliases), len(removeCmd.Aliases))
	}

	for i, alias := range expectedAliases {
		if i < len(removeCmd.Aliases) && removeCmd.Aliases[i] != alias {
			t.Errorf("removeCmd.Aliases[%d] = %v, want %v", i, removeCmd.Aliases[i], alias)
		}
	}
}

func TestRemoveCmdExecution(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_remove_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Capture output
	var buf bytes.Buffer
	removeCmd.SetOut(&buf)
	removeCmd.SetErr(&buf)

	// Execute remove command with non-existent template (should fail)
	err = removeCmd.RunE(removeCmd, []string{"test-template"})
	if err == nil {
		t.Error("removeCmd should return error for non-existent template")
	}

	// Reset
	removeCmd.SetOut(nil)
	removeCmd.SetErr(nil)
}

func TestValidateCmd(t *testing.T) {
	// Test validate command properties
	if validateCmd == nil {
		t.Fatal("validateCmd should not be nil")
	}

	if validateCmd.Use != "validate [path]" {
		t.Errorf("validateCmd.Use = %v, want %v", validateCmd.Use, "validate [path]")
	}

	if validateCmd.Short != "Validate a template" {
		t.Errorf("validateCmd.Short = %v, want %v", validateCmd.Short, "Validate a template")
	}
}

func TestValidateCmdExecution_ValidPath(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "ason_validate_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a valid ason.yaml file
	testFile := filepath.Join(tmpDir, "ason.yaml")
	yamlContent := `name: "Test Template"
description: "A test template"
version: "1.0.0"
variables:
  - name: project_name
    required: true`
	err = os.WriteFile(testFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Capture output
	var buf bytes.Buffer
	validateCmd.SetOut(&buf)

	// Execute validate command with valid directory
	err = validateCmd.RunE(validateCmd, []string{tmpDir})
	if err != nil {
		t.Fatalf("validateCmd execution failed: %v", err)
	}

	// Test passed if no error occurred
	// Reset
	validateCmd.SetOut(nil)
}

func TestValidateCmdExecution_InvalidPath(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	validateCmd.SetOut(&buf)
	validateCmd.SetErr(&buf)

	// Execute validate command with invalid path
	err := validateCmd.RunE(validateCmd, []string{"/non/existent/path"})
	if err == nil {
		t.Error("validateCmd should return error for non-existent path")
	}

	// Reset
	validateCmd.SetOut(nil)
	validateCmd.SetErr(nil)
}

func TestCommandsAreRegistered(t *testing.T) {
	// Test that all commands are properly registered with root
	commands := rootCmd.Commands()

	expectedCommands := map[string]bool{
		"list":     false,
		"add":      false,
		"remove":   false,
		"validate": false,
		"new":      false,
	}

	for _, cmd := range commands {
		for cmdName := range expectedCommands {
			if strings.HasPrefix(cmd.Use, cmdName) {
				expectedCommands[cmdName] = true
				break
			}
		}
	}

	for cmdName, found := range expectedCommands {
		if !found {
			t.Errorf("Command %v not found in root command", cmdName)
		}
	}
}
