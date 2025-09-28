package registry

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewRegistry(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_home_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	registry, err := NewRegistry()
	if err != nil {
		t.Fatalf("NewRegistry() failed: %v", err)
	}

	if registry == nil {
		t.Fatal("NewRegistry() returned nil registry")
	}

	// Check that registry directory was created (XDG-compliant)
	expectedPath := filepath.Join(tmpHome, ".local", "share", "ason")
	if registry.path != expectedPath {
		t.Errorf("Registry path = %v, want %v", registry.path, expectedPath)
	}

	// Check that directories exist
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("Registry directory was not created")
	}

	templatesPath := filepath.Join(expectedPath, "templates")
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		t.Error("Templates directory was not created")
	}
}

func TestRegistry_List_Empty(t *testing.T) {
	// Create temporary registry
	tmpDir, err := os.MkdirTemp("", "ason_registry_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry := &Registry{path: tmpDir}

	templates, err := registry.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	// Should return empty list for new registry
	if len(templates) != 0 {
		t.Errorf("Expected empty template list, got %d templates", len(templates))
	}
}

func TestRegistry_AddAndList(t *testing.T) {
	// Create temporary registry
	tmpDir, err := os.MkdirTemp("", "ason_registry_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry := &Registry{path: tmpDir}

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

	err = os.WriteFile(filepath.Join(testTemplateDir, "ason.toml"), []byte(`
name = "Test Template"
description = "A test template"

[[variables]]
name = "project_name"
required = true
`), 0644)
	if err != nil {
		t.Fatalf("Failed to create ason.toml: %v", err)
	}

	// Add template to registry
	err = registry.Add("test-template", testTemplateDir, "Test description", "test")
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// List templates
	templates, err := registry.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(templates))
	}

	tmpl := templates[0]
	if tmpl.Name != "test-template" {
		t.Errorf("Template name = %v, want %v", tmpl.Name, "test-template")
	}
	if tmpl.Description != "Test description" {
		t.Errorf("Template description = %v, want %v", tmpl.Description, "Test description")
	}
	if tmpl.Type != "test" {
		t.Errorf("Template type = %v, want %v", tmpl.Type, "test")
	}
	if tmpl.Size <= 0 {
		t.Errorf("Template size should be > 0, got %d", tmpl.Size)
	}
	if tmpl.Files != 2 {
		t.Errorf("Template files = %d, want 2", tmpl.Files)
	}
	if len(tmpl.Variables) != 1 || tmpl.Variables[0] != "project_name" {
		t.Errorf("Template variables = %v, want [project_name]", tmpl.Variables)
	}
}

func TestRegistry_Get(t *testing.T) {
	// Create temporary registry
	tmpDir, err := os.MkdirTemp("", "ason_registry_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry := &Registry{path: tmpDir}

	// Create a test template directory
	testTemplateDir, err := os.MkdirTemp("", "test_template")
	if err != nil {
		t.Fatalf("Failed to create test template dir: %v", err)
	}
	defer os.RemoveAll(testTemplateDir)

	// Add some files to the template
	err = os.WriteFile(filepath.Join(testTemplateDir, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Add template to registry
	err = registry.Add("test-template", testTemplateDir, "Test description", "test")
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Test case 1: Template exists
	gotPath, err := registry.Get("test-template")
	if err != nil {
		t.Errorf("Get() failed for existing template: %v", err)
	}
	expectedPath := filepath.Join(tmpDir, "templates", "test-template")
	if gotPath != expectedPath {
		t.Errorf("Get() = %v, want %v", gotPath, expectedPath)
	}

	// Test case 2: Template doesn't exist
	_, err = registry.Get("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

func TestRegistry_Remove(t *testing.T) {
	// Create temporary registry
	tmpDir, err := os.MkdirTemp("", "ason_registry_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry := &Registry{path: tmpDir}

	// Create a test template directory
	testTemplateDir, err := os.MkdirTemp("", "test_template")
	if err != nil {
		t.Fatalf("Failed to create test template dir: %v", err)
	}
	defer os.RemoveAll(testTemplateDir)

	// Add some files to the template
	err = os.WriteFile(filepath.Join(testTemplateDir, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Add template to registry
	err = registry.Add("test-template", testTemplateDir, "Test description", "test")
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Verify template exists
	templates, err := registry.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(templates) != 1 {
		t.Fatalf("Expected 1 template before removal, got %d", len(templates))
	}

	// Remove template
	err = registry.Remove("test-template", false, "")
	if err != nil {
		t.Fatalf("Remove() failed: %v", err)
	}

	// Verify template is gone
	templates, err = registry.List()
	if err != nil {
		t.Fatalf("List() failed after removal: %v", err)
	}
	if len(templates) != 0 {
		t.Errorf("Expected 0 templates after removal, got %d", len(templates))
	}

	// Verify template directory is gone
	templatePath := filepath.Join(tmpDir, "templates", "test-template")
	if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
		t.Error("Template directory should be removed")
	}
}

func TestRegistry_RemoveWithBackup(t *testing.T) {
	// Create temporary registry
	tmpDir, err := os.MkdirTemp("", "ason_registry_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry := &Registry{path: tmpDir}

	// Create a test template directory
	testTemplateDir, err := os.MkdirTemp("", "test_template")
	if err != nil {
		t.Fatalf("Failed to create test template dir: %v", err)
	}
	defer os.RemoveAll(testTemplateDir)

	// Add some files to the template
	err = os.WriteFile(filepath.Join(testTemplateDir, "test.txt"), []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Add template to registry
	err = registry.Add("test-template", testTemplateDir, "Test description", "test")
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Remove template with backup
	backupDir := filepath.Join(tmpDir, "test-backups")
	err = registry.Remove("test-template", true, backupDir)
	if err != nil {
		t.Fatalf("Remove() with backup failed: %v", err)
	}

	// Check that backup was created
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		t.Error("Backup directory was not created")
	}

	// Check that backup contains files
	backupFiles, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("Failed to read backup directory: %v", err)
	}
	if len(backupFiles) == 0 {
		t.Error("Backup directory is empty")
	}
}

func TestRegistry_RemoveNonExistent(t *testing.T) {
	// Create temporary registry
	tmpDir, err := os.MkdirTemp("", "ason_registry_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	registry := &Registry{path: tmpDir}

	// Try to remove non-existent template
	err = registry.Remove("non-existent", false, "")
	if err == nil {
		t.Error("Expected error when removing non-existent template, got nil")
	}
}

func TestTemplateConfig(t *testing.T) {
	// Test TemplateConfig struct
	config := TemplateConfig{
		Name:        "Test Template",
		Description: "A test template",
		Version:     "1.0.0",
		Type:        "test",
		Variables: []TemplateVariable{
			{
				Name:        "project_name",
				Description: "Name of the project",
				Required:    true,
				Default:     "my-project",
			},
		},
	}

	if config.Name != "Test Template" {
		t.Errorf("TemplateConfig.Name = %v, want %v", config.Name, "Test Template")
	}
	if len(config.Variables) != 1 {
		t.Errorf("TemplateConfig.Variables length = %d, want 1", len(config.Variables))
	}
	if config.Variables[0].Name != "project_name" {
		t.Errorf("Variable name = %v, want %v", config.Variables[0].Name, "project_name")
	}
}

func TestTemplateEntry(t *testing.T) {
	// Test TemplateEntry struct
	now := time.Now()
	entry := TemplateEntry{
		Name:        "test",
		Path:        "/path/to/test",
		Description: "Test template",
		Source:      "/source/path",
		Type:        "example",
		Size:        1024,
		Files:       5,
		Added:       now,
		Variables:   []string{"name", "version"},
	}

	if entry.Name != "test" {
		t.Errorf("TemplateEntry.Name = %v, want %v", entry.Name, "test")
	}
	if entry.Path != "/path/to/test" {
		t.Errorf("TemplateEntry.Path = %v, want %v", entry.Path, "/path/to/test")
	}
	if entry.Description != "Test template" {
		t.Errorf("TemplateEntry.Description = %v, want %v", entry.Description, "Test template")
	}
	if entry.Type != "example" {
		t.Errorf("TemplateEntry.Type = %v, want %v", entry.Type, "example")
	}
	if entry.Size != 1024 {
		t.Errorf("TemplateEntry.Size = %v, want %v", entry.Size, 1024)
	}
	if entry.Files != 5 {
		t.Errorf("TemplateEntry.Files = %v, want %v", entry.Files, 5)
	}
	if !entry.Added.Equal(now) {
		t.Errorf("TemplateEntry.Added = %v, want %v", entry.Added, now)
	}
	if len(entry.Variables) != 2 || entry.Variables[0] != "name" || entry.Variables[1] != "version" {
		t.Errorf("TemplateEntry.Variables = %v, want [name version]", entry.Variables)
	}
}
