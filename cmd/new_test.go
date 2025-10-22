package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCmd(t *testing.T) {
	// Test new command properties
	if newCmd == nil {
		t.Fatal("newCmd should not be nil")
	}

	if newCmd.Use != "new [template] [output]" {
		t.Errorf("newCmd.Use = %v, want %v", newCmd.Use, "new [template] [output]")
	}

	if newCmd.Short != "Create a new project from a template" {
		t.Errorf("newCmd.Short = %v, want %v", newCmd.Short, "Create a new project from a template")
	}

	if !strings.Contains(newCmd.Long, "Create a new project from a template") {
		t.Error("newCmd.Long should contain description")
	}
}

func TestNewCmdFlags(t *testing.T) {
	// Test that flags are properly defined
	flags := newCmd.Flags()

	// Test output flag
	outputFlag := flags.Lookup("output")
	if outputFlag == nil {
		t.Error("--output flag should be defined")
	}

	// Test no-input flag
	noInputFlag := flags.Lookup("no-input")
	if noInputFlag == nil {
		t.Error("--no-input flag should be defined")
	}

	// Test var flag
	varFlag := flags.Lookup("var")
	if varFlag == nil {
		t.Error("--var flag should be defined")
	}

	// Test dry-run flag
	dryRunFlag := flags.Lookup("dry-run")
	if dryRunFlag == nil {
		t.Error("--dry-run flag should be defined")
	}
}

func TestNewCmdDryRun(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_new_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Create temporary template directory
	tmpTemplate, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp template: %v", err)
	}
	defer os.RemoveAll(tmpTemplate)

	// Add a test file to the template
	err = os.WriteFile(filepath.Join(tmpTemplate, "README.md"), []byte("# {{ name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Capture output
	var buf bytes.Buffer
	newCmd.SetOut(&buf)

	// Set dry-run flag
	dryRun = true
	defer func() { dryRun = false }()

	// Execute new command with dry run
	err = newCmd.RunE(newCmd, []string{tmpTemplate})
	if err != nil {
		t.Fatalf("newCmd dry run execution failed: %v", err)
	}

	// Test passed if no error occurred
	// Reset
	newCmd.SetOut(nil)
}

func TestNewCmdWithExistingTemplate(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_new_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Create template directory structure (XDG compliant)
	registryDir := filepath.Join(tmpHome, ".local", "share", "ason", "templates")
	err = os.MkdirAll(registryDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create registry dir: %v", err)
	}

	templateDir := filepath.Join(registryDir, "test-template")
	err = os.Mkdir(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Add a test file to the template
	err = os.WriteFile(filepath.Join(templateDir, "README.md"), []byte("# {{ name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Create registry metadata file (TOML format)
	registryFile := filepath.Join(tmpHome, ".local", "share", "ason", "registry.toml")
	registryContent := `[templates.test-template]
name = "test-template"
path = "` + templateDir + `"
description = "Test template"
type = "test"
size = 100
files = 1
added = 2023-01-01T00:00:00Z
variables = []
`
	err = os.WriteFile(registryFile, []byte(registryContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create registry file: %v", err)
	}

	// Create output directory
	outputDir, err := os.MkdirTemp("", "ason_output_test")
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Capture output
	var buf bytes.Buffer
	newCmd.SetOut(&buf)

	// Execute new command
	err = newCmd.RunE(newCmd, []string{"test-template", outputDir})
	if err != nil {
		t.Fatalf("newCmd execution failed: %v", err)
	}

	// Test passed if no error occurred
	// Reset
	newCmd.SetOut(nil)
}

func TestNewCmdWithDirectPath(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_new_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Create template directory
	templateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}
	defer os.RemoveAll(templateDir)

	// Add a test file to the template
	err = os.WriteFile(filepath.Join(templateDir, "README.md"), []byte("# {{ name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Create output directory
	outputDir, err := os.MkdirTemp("", "ason_output_test")
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Capture output
	var buf bytes.Buffer
	newCmd.SetOut(&buf)

	// Execute new command with direct path
	err = newCmd.RunE(newCmd, []string{templateDir, outputDir})
	if err != nil {
		t.Fatalf("newCmd execution with direct path failed: %v", err)
	}

	// Test passed if no error occurred
	// Reset
	newCmd.SetOut(nil)
}

func TestNewCmdWithNonExistentTemplate(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_new_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Capture output
	var buf bytes.Buffer
	newCmd.SetOut(&buf)
	newCmd.SetErr(&buf)

	// Execute new command with non-existent template
	err = newCmd.RunE(newCmd, []string{"non-existent-template"})
	if err == nil {
		t.Error("newCmd should return error for non-existent template")
	}

	// Reset
	newCmd.SetOut(nil)
	newCmd.SetErr(nil)
}

func TestNewCmdVariables(t *testing.T) {
	// Test that global variables exist and have correct types
	tests := []struct {
		name     string
		variable interface{}
	}{
		{"outputDir", &outputDir},
		{"noInput", &noInput},
		{"extraVars", &extraVars},
		{"configFile", &configFile},
		{"skipHooks", &skipHooks},
		{"dryRun", &dryRun},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.variable == nil {
				t.Errorf("Variable %s should not be nil", tt.name)
			}
		})
	}
}

func TestNewCmdWithExtraVars(t *testing.T) {
	// Save original values
	originalExtraVars := extraVars
	defer func() { extraVars = originalExtraVars }()

	// Set extra vars
	extraVars = map[string]string{
		"name":    "test-project",
		"version": "1.0.0",
	}

	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_new_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Create template directory
	templateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}
	defer os.RemoveAll(templateDir)

	// Set dry run to avoid actual generation
	dryRun = true
	defer func() { dryRun = false }()

	// Capture output
	var buf bytes.Buffer
	newCmd.SetOut(&buf)

	// Execute new command with extra vars
	err = newCmd.RunE(newCmd, []string{templateDir})
	if err != nil {
		t.Fatalf("newCmd execution with extra vars failed: %v", err)
	}

	// Reset
	newCmd.SetOut(nil)
}
