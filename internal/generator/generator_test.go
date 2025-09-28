package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/madstone-tech/ason/internal/engine"
	"github.com/madstone-tech/ason/internal/template"
)

// MockEngine for testing
type MockEngine struct {
	renderFunc     func(string, map[string]interface{}) (string, error)
	renderFileFunc func(string, map[string]interface{}) (string, error)
}

func (m *MockEngine) Render(tmpl string, context map[string]interface{}) (string, error) {
	if m.renderFunc != nil {
		return m.renderFunc(tmpl, context)
	}
	// Simple variable substitution for testing
	result := tmpl
	for key, value := range context {
		placeholder := "{{ " + key + " }}"
		if str, ok := value.(string); ok {
			result = strings.ReplaceAll(result, placeholder, str)
		}
	}
	return result, nil
}

func (m *MockEngine) RenderFile(filepath string, context map[string]interface{}) (string, error) {
	if m.renderFileFunc != nil {
		return m.renderFileFunc(filepath, context)
	}
	return "", nil
}

func TestNew(t *testing.T) {
	tmpl := &Template{
		Path: "/test/path",
		Config: &template.Config{
			Name: "test",
		},
	}
	mockEngine := &MockEngine{}

	generator := New(tmpl, mockEngine)

	if generator == nil {
		t.Fatal("New() returned nil")
	}

	if generator.template != tmpl {
		t.Error("Generator template not set correctly")
	}

	if generator.engine != mockEngine {
		t.Error("Generator engine not set correctly")
	}
}

func TestTemplate(t *testing.T) {
	// Test Template struct
	config := &template.Config{
		Name:        "test-template",
		Description: "Test description",
		Version:     "1.0.0",
	}

	tmpl := &Template{
		Path:   "/test/path",
		Config: config,
	}

	if tmpl.Path != "/test/path" {
		t.Errorf("Template.Path = %v, want %v", tmpl.Path, "/test/path")
	}

	if tmpl.Config != config {
		t.Error("Template.Config not set correctly")
	}
}

func TestOptions(t *testing.T) {
	// Test Options struct
	opts := Options{
		SkipHooks: true,
		DryRun:    true,
		Verbose:   false,
	}

	if !opts.SkipHooks {
		t.Error("Options.SkipHooks should be true")
	}

	if !opts.DryRun {
		t.Error("Options.DryRun should be true")
	}

	if opts.Verbose {
		t.Error("Options.Verbose should be false")
	}
}

func TestGenerator_Generate_DryRun(t *testing.T) {
	// Create temporary template directory for dry run test
	tmpTemplateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp template dir: %v", err)
	}
	defer os.RemoveAll(tmpTemplateDir)

	// Create test template files
	err = os.WriteFile(filepath.Join(tmpTemplateDir, "README.md"), []byte("# {{ name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	tmpl := &Template{
		Path: tmpTemplateDir,
	}
	mockEngine := &MockEngine{}
	generator := New(tmpl, mockEngine)

	// Test dry run
	context := map[string]interface{}{
		"name": "test-project",
	}

	opts := Options{
		DryRun: true,
	}

	err = generator.Generate("/tmp/test-output", context, opts)
	if err != nil {
		t.Errorf("Generate() with dry run failed: %v", err)
	}

	// Verify no directory was created (dry run)
	if _, err := os.Stat("/tmp/test-output"); !os.IsNotExist(err) {
		t.Error("Directory should not exist after dry run")
		os.RemoveAll("/tmp/test-output") // Clean up if it was created
	}
}

func TestGenerator_Generate_RealRun(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ason_generate_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create temporary template directory
	tmpTemplateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp template dir: %v", err)
	}
	defer os.RemoveAll(tmpTemplateDir)

	// Create test template files
	err = os.WriteFile(filepath.Join(tmpTemplateDir, "README.md"), []byte("# {{ name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpTemplateDir, "config.txt"), []byte("project: {{ name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	tmpl := &Template{
		Path: tmpTemplateDir,
	}
	mockEngine := &MockEngine{}
	generator := New(tmpl, mockEngine)

	outputPath := filepath.Join(tmpDir, "test-output")
	context := map[string]interface{}{
		"name": "test-project",
	}

	opts := Options{
		DryRun: false,
	}

	err = generator.Generate(outputPath, context, opts)
	if err != nil {
		t.Errorf("Generate() failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output directory was not created")
	}

	// Verify files were created and processed
	readmeContent, err := os.ReadFile(filepath.Join(outputPath, "README.md"))
	if err != nil {
		t.Errorf("README.md was not created: %v", err)
	} else if string(readmeContent) != "# test-project" {
		t.Errorf("README.md content = %q, want %q", string(readmeContent), "# test-project")
	}

	configContent, err := os.ReadFile(filepath.Join(outputPath, "config.txt"))
	if err != nil {
		t.Errorf("config.txt was not created: %v", err)
	} else if string(configContent) != "project: test-project" {
		t.Errorf("config.txt content = %q, want %q", string(configContent), "project: test-project")
	}
}

func TestGenerator_Generate_DirectoryCreationError(t *testing.T) {
	// Create temporary template directory
	tmpTemplateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp template dir: %v", err)
	}
	defer os.RemoveAll(tmpTemplateDir)

	// Create test template file
	err = os.WriteFile(filepath.Join(tmpTemplateDir, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	tmpl := &Template{
		Path: tmpTemplateDir,
	}
	mockEngine := &MockEngine{}
	generator := New(tmpl, mockEngine)

	// Try to create directory in location that should fail on most systems
	outputPath := "/proc/invalid/directory"
	context := map[string]interface{}{}

	opts := Options{
		DryRun: false,
	}

	err = generator.Generate(outputPath, context, opts)
	if err == nil {
		t.Error("Expected error when creating directory in invalid location, got nil")
	}
}

func TestGenerator_WithRealEngine(t *testing.T) {
	// Create temporary template directory
	tmpTemplateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp template dir: %v", err)
	}
	defer os.RemoveAll(tmpTemplateDir)

	// Create test template file with Pongo2 syntax
	err = os.WriteFile(filepath.Join(tmpTemplateDir, "test.md"), []byte("# {{ name }}\n\nAuthor: {{ author | default:\"Unknown\" }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Test with real Pongo2 engine
	tmpl := &Template{
		Path: tmpTemplateDir,
	}
	realEngine := engine.NewPongo2Engine()
	generator := New(tmpl, realEngine)

	if generator.engine == nil {
		t.Error("Generator should have real engine")
	}

	// Create temporary output directory
	tmpOutputDir, err := os.MkdirTemp("", "ason_output_test")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(tmpOutputDir)

	// Test real generation with real engine
	context := map[string]interface{}{
		"name":   "Real Test Project",
		"author": "Test Author",
	}

	opts := Options{
		DryRun: false,
	}

	err = generator.Generate(tmpOutputDir, context, opts)
	if err != nil {
		t.Errorf("Generate() with real engine failed: %v", err)
	}

	// Verify the file was created and processed correctly
	testContent, err := os.ReadFile(filepath.Join(tmpOutputDir, "test.md"))
	if err != nil {
		t.Errorf("test.md was not created: %v", err)
	} else {
		expected := "# Real Test Project\n\nAuthor: Test Author"
		if string(testContent) != expected {
			t.Errorf("test.md content = %q, want %q", string(testContent), expected)
		}
	}
}

func TestGenerator_BinaryFileHandling(t *testing.T) {
	// Create temporary template directory
	tmpTemplateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp template dir: %v", err)
	}
	defer os.RemoveAll(tmpTemplateDir)

	// Create a binary file (simulate by using non-UTF8 content)
	binaryContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	err = os.WriteFile(filepath.Join(tmpTemplateDir, "image.png"), binaryContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	// Create a text file
	err = os.WriteFile(filepath.Join(tmpTemplateDir, "README.md"), []byte("# {{ name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create text file: %v", err)
	}

	tmpl := &Template{
		Path: tmpTemplateDir,
	}
	mockEngine := &MockEngine{}
	generator := New(tmpl, mockEngine)

	// Create temporary output directory
	tmpOutputDir, err := os.MkdirTemp("", "ason_output_test")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(tmpOutputDir)

	context := map[string]interface{}{
		"name": "Test Project",
	}

	opts := Options{
		DryRun: false,
	}

	err = generator.Generate(tmpOutputDir, context, opts)
	if err != nil {
		t.Errorf("Generate() failed: %v", err)
	}

	// Verify binary file was copied as-is
	copiedBinary, err := os.ReadFile(filepath.Join(tmpOutputDir, "image.png"))
	if err != nil {
		t.Errorf("Binary file was not copied: %v", err)
	} else if string(copiedBinary) != string(binaryContent) {
		t.Error("Binary file content was modified")
	}

	// Verify text file was processed
	processedText, err := os.ReadFile(filepath.Join(tmpOutputDir, "README.md"))
	if err != nil {
		t.Errorf("Text file was not processed: %v", err)
	} else if string(processedText) != "# Test Project" {
		t.Errorf("Text file was not processed correctly: got %q", string(processedText))
	}
}

func TestGenerator_NestedDirectories(t *testing.T) {
	// Create temporary template directory with nested structure
	tmpTemplateDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp template dir: %v", err)
	}
	defer os.RemoveAll(tmpTemplateDir)

	// Create nested directory structure
	srcDir := filepath.Join(tmpTemplateDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create src directory: %v", err)
	}

	// Create files in nested directories
	err = os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package {{ package_name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpTemplateDir, "README.md"), []byte("# {{ project_name }}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	tmpl := &Template{
		Path: tmpTemplateDir,
	}
	mockEngine := &MockEngine{}
	generator := New(tmpl, mockEngine)

	// Create temporary output directory
	tmpOutputDir, err := os.MkdirTemp("", "ason_output_test")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(tmpOutputDir)

	context := map[string]interface{}{
		"project_name": "MyProject",
		"package_name": "main",
	}

	opts := Options{
		DryRun: false,
	}

	err = generator.Generate(tmpOutputDir, context, opts)
	if err != nil {
		t.Errorf("Generate() failed: %v", err)
	}

	// Verify nested directory was created
	outputSrcDir := filepath.Join(tmpOutputDir, "src")
	if _, err := os.Stat(outputSrcDir); os.IsNotExist(err) {
		t.Error("Nested src directory was not created")
	}

	// Verify nested file was processed
	mainGoContent, err := os.ReadFile(filepath.Join(outputSrcDir, "main.go"))
	if err != nil {
		t.Errorf("main.go was not created: %v", err)
	} else if string(mainGoContent) != "package main" {
		t.Errorf("main.go content = %q, want %q", string(mainGoContent), "package main")
	}

	// Verify root file was processed
	readmeContent, err := os.ReadFile(filepath.Join(tmpOutputDir, "README.md"))
	if err != nil {
		t.Errorf("README.md was not created: %v", err)
	} else if string(readmeContent) != "# MyProject" {
		t.Errorf("README.md content = %q, want %q", string(readmeContent), "# MyProject")
	}
}

func TestGenerator_shouldProcessAsTemplate(t *testing.T) {
	generator := &Generator{}

	// Test text files (should be processed)
	textFiles := []string{
		"README.md",
		"config.yaml",
		"script.sh",
		"source.go",
		"package.json",
		"Dockerfile",
	}

	for _, file := range textFiles {
		if !generator.shouldProcessAsTemplate(file) {
			t.Errorf("shouldProcessAsTemplate(%q) = false, want true", file)
		}
	}

	// Test binary files (should not be processed)
	binaryFiles := []string{
		"image.png",
		"photo.jpg",
		"document.pdf",
		"archive.zip",
		"program.exe",
		"library.so",
		"font.woff",
	}

	for _, file := range binaryFiles {
		if generator.shouldProcessAsTemplate(file) {
			t.Errorf("shouldProcessAsTemplate(%q) = true, want false", file)
		}
	}
}
