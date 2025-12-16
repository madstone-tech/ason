package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/madstone-tech/ason/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRegistry tests Registry creation with default path
func TestNewRegistry(t *testing.T) {
	registry, err := pkg.NewRegistry()
	require.NoError(t, err)
	assert.NotNil(t, registry)
}

// TestNewRegistryAt tests Registry creation with custom path
func TestNewRegistryAt(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "custom-registry.toml")

	registry, err := pkg.NewRegistryAt(regPath)
	require.NoError(t, err)
	assert.NotNil(t, registry)

	// Verify empty registry
	templates, err := registry.List()
	require.NoError(t, err)
	assert.Len(t, templates, 0)
}

// TestNewRegistryAtEmptyPath tests Registry creation with empty path
func TestNewRegistryAtEmptyPath(t *testing.T) {
	registry, err := pkg.NewRegistryAt("")
	assert.Error(t, err)
	assert.Nil(t, registry)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestRegisterTemplate tests registering a new template
func TestRegisterTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	// Create a template directory
	templateDir := filepath.Join(tmpDir, "my-template")
	err := os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	// Register the template
	err = registry.Register("my-template", templateDir, "Test template")
	require.NoError(t, err)

	// Verify it appears in the list
	templates, err := registry.List()
	require.NoError(t, err)
	assert.Len(t, templates, 1)
	assert.Equal(t, "my-template", templates[0].Name)
	assert.Equal(t, templateDir, templates[0].Path)
	assert.Equal(t, "Test template", templates[0].Description)
}

// TestRegisterMultipleTemplates tests registering multiple templates
func TestRegisterMultipleTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	// Create multiple template directories
	for i := 1; i <= 3; i++ {
		templateDir := filepath.Join(tmpDir, "template-"+string(rune('0'+i)))
		_ = os.MkdirAll(templateDir, 0755)
		_ = registry.Register("template-"+string(rune('0'+i)), templateDir, "Template "+string(rune('0'+i)))
	}

	// Verify all appear in the list
	templates, err := registry.List()
	require.NoError(t, err)
	assert.Len(t, templates, 3)

	// Verify alphabetical ordering
	assert.Equal(t, "template-1", templates[0].Name)
	assert.Equal(t, "template-2", templates[1].Name)
	assert.Equal(t, "template-3", templates[2].Name)
}

// TestRegisterEmptyName tests registering with empty name
func TestRegisterEmptyName(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	templateDir := filepath.Join(tmpDir, "template")
	_ = os.MkdirAll(templateDir, 0755)

	err := registry.Register("", templateDir, "Test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestRegisterEmptyPath tests registering with empty path
func TestRegisterEmptyPath(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	err := registry.Register("test", "", "Test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestRegisterNonexistentPath tests registering non-existent path
func TestRegisterNonexistentPath(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	err := registry.Register("test", "/nonexistent/path", "Test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

// TestRegisterInvalidName tests registering with invalid name format
func TestRegisterInvalidName(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	templateDir := filepath.Join(tmpDir, "template")
	_ = os.MkdirAll(templateDir, 0755)

	// Invalid name starting with number
	err := registry.Register("123invalid", templateDir, "Test")
	assert.Error(t, err)
}

// TestRegisterOverwrite tests overwriting an existing template
func TestRegisterOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	// Create two template directories
	templateDir1 := filepath.Join(tmpDir, "template-1")
	_ = os.MkdirAll(templateDir1, 0755)
	templateDir2 := filepath.Join(tmpDir, "template-2")
	_ = os.MkdirAll(templateDir2, 0755)

	// Register first template
	_ = registry.Register("test", templateDir1, "First")

	// Register with same name (should overwrite)
	_ = registry.Register("test", templateDir2, "Second")

	// Verify only one exists and it's the second
	templates, _ := registry.List()
	assert.Len(t, templates, 1)
	assert.Equal(t, templateDir2, templates[0].Path)
	assert.Equal(t, "Second", templates[0].Description)
}

// TestRemoveTemplate tests removing a template
func TestRemoveTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	// Register a template
	templateDir := filepath.Join(tmpDir, "template")
	_ = os.MkdirAll(templateDir, 0755)
	_ = registry.Register("test", templateDir, "Test")

	// Verify it exists
	templates, _ := registry.List()
	assert.Len(t, templates, 1)

	// Remove it
	err := registry.Remove("test")
	require.NoError(t, err)

	// Verify it's gone
	templates, _ = registry.List()
	assert.Len(t, templates, 0)
}

// TestRemoveNonexistent tests removing a non-existent template (should be idempotent)
func TestRemoveNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	// Should not error when removing non-existent template
	err := registry.Remove("nonexistent")
	assert.NoError(t, err)
}

// TestListEmpty tests listing empty registry
func TestListEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	templates, err := registry.List()
	require.NoError(t, err)
	assert.Empty(t, templates)
}

// TestRegistryPersistence tests that registry survives close/reopen
func TestRegistryPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")

	// Create and populate registry
	templateDir := filepath.Join(tmpDir, "template")
	_ = os.MkdirAll(templateDir, 0755)

	registry1, _ := pkg.NewRegistryAt(regPath)
	_ = registry1.Register("test", templateDir, "Test template")

	// Create new registry from same path
	registry2, _ := pkg.NewRegistryAt(regPath)
	templates, _ := registry2.List()

	// Verify data persisted
	assert.Len(t, templates, 1)
	assert.Equal(t, "test", templates[0].Name)
}

// TestRegistryConcurrentReads tests concurrent read operations
func TestRegistryConcurrentReads(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	// Register a template
	templateDir := filepath.Join(tmpDir, "template")
	_ = os.MkdirAll(templateDir, 0755)
	_ = registry.Register("test", templateDir, "Test")

	// Run concurrent reads
	results := RunConcurrent(t, 5, func(index int) error {
		templates, err := registry.List()
		if err != nil {
			return err
		}
		if len(templates) != 1 {
			t.Errorf("expected 1 template, got %d", len(templates))
		}
		return nil
	})

	assert.Len(t, results, 5)
	for _, result := range results {
		assert.True(t, result.Success, "concurrent read failed: %v", result.Error)
	}
}

// TestRegistryConcurrentWrites tests concurrent write operations (should serialize)
func TestRegistryConcurrentWrites(t *testing.T) {
	tmpDir := t.TempDir()
	regPath := filepath.Join(tmpDir, "registry.toml")
	registry, _ := pkg.NewRegistryAt(regPath)

	// Create template directories
	templateDirs := make([]string, 5)
	for i := 0; i < 5; i++ {
		templateDirs[i] = filepath.Join(tmpDir, "template-"+string(rune('0'+i)))
		_ = os.MkdirAll(templateDirs[i], 0755)
	}

	// Run concurrent writes
	results := RunConcurrent(t, 5, func(index int) error {
		return registry.Register("template-"+string(rune('0'+index)), templateDirs[index], "Test "+string(rune('0'+index)))
	})

	// All writes should succeed
	for _, result := range results {
		assert.True(t, result.Success, "concurrent write failed: %v", result.Error)
	}

	// Final registry should have all templates
	templates, _ := registry.List()
	assert.Len(t, templates, 5)
}
