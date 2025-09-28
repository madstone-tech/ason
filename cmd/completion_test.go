package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompleteTemplateNames(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_completion_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Test with empty registry (should return no completions)
	completions, directive := completeTemplateNames(nil, []string{}, "test")
	if len(completions) != 0 {
		t.Errorf("Expected no completions for empty registry, got %d", len(completions))
	}

	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("Expected NoFileComp directive, got %v", directive)
	}
}

func TestCompleteTemplateNamesOrPaths(t *testing.T) {
	// Save original home directory
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary home directory
	tmpHome, err := os.MkdirTemp("", "ason_completion_test")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	os.Setenv("HOME", tmpHome)

	// Test with empty registry (should fall back to directory completion)
	completions, directive := completeTemplateNamesOrPaths(nil, []string{}, "test")
	if len(completions) != 0 {
		t.Errorf("Expected no completions for empty registry, got %d", len(completions))
	}

	if directive != cobra.ShellCompDirectiveFilterDirs {
		t.Errorf("Expected FilterDirs directive, got %v", directive)
	}
}

func TestCompleteOutputPaths(t *testing.T) {
	completions, directive := completeOutputPaths(nil, []string{}, "test")
	if len(completions) != 0 {
		t.Errorf("Expected no completions from completeOutputPaths, got %d", len(completions))
	}

	if directive != cobra.ShellCompDirectiveFilterDirs {
		t.Errorf("Expected FilterDirs directive, got %v", directive)
	}
}

func TestCompleteTemplatePaths(t *testing.T) {
	// Create temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "ason_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create test files
	testFiles := []string{
		"template.yaml",
		"config.json",
		"readme.md",
		"script.sh",
	}

	for _, file := range testFiles {
		err := os.WriteFile(file, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create test directory
	err = os.Mkdir("templates", 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test completion for template files
	completions, directive := completeTemplatePaths(nil, []string{}, "temp")

	// Should include template.yaml and templates/ directory
	expectedCount := 2 // template.yaml and templates/
	if len(completions) != expectedCount {
		t.Errorf("Expected %d completions, got %d: %v", expectedCount, len(completions), completions)
	}

	if directive != cobra.ShellCompDirectiveDefault {
		t.Errorf("Expected Default directive, got %v", directive)
	}

	// Check that template.yaml is included
	found := false
	for _, completion := range completions {
		if completion == "template.yaml" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected template.yaml in completions")
	}

	// Check that templates/ directory is included with trailing slash
	found = false
	for _, completion := range completions {
		if completion == "templates"+string(filepath.Separator) {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected templates/ directory in completions")
	}
}

func TestIsTemplateFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"template.yaml", true},
		{"config.yml", true},
		{"data.json", true},
		{"file.tmpl", true},
		{"example.template", true},
		{"readme.md", false},
		{"script.sh", false},
		{"binary.exe", false},
		{"no-extension", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := isTemplateFile(tt.filename)
			if result != tt.expected {
				t.Errorf("isTemplateFile(%s) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestCompleteVariableKeys(t *testing.T) {
	tests := []struct {
		toComplete string
		minCount   int
	}{
		{"", 13},    // Should return all common variables
		{"name", 2}, // Should match name=, project_name=
		{"use_", 3}, // Should match use_docker=, use_tests=, use_ci=
		{"xyz", 0},  // Should match nothing
	}

	for _, tt := range tests {
		t.Run(tt.toComplete, func(t *testing.T) {
			completions, directive := completeVariableKeys(nil, []string{}, tt.toComplete)

			if len(completions) < tt.minCount {
				t.Errorf("Expected at least %d completions for '%s', got %d", tt.minCount, tt.toComplete, len(completions))
			}

			expectedDirective := cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
			if directive != expectedDirective {
				t.Errorf("Expected NoSpace|NoFileComp directive, got %v", directive)
			}

			// Verify all completions end with '='
			for _, completion := range completions {
				if completion[len(completion)-1] != '=' {
					t.Errorf("Expected completion '%s' to end with '='", completion)
				}
			}
		})
	}
}

func TestCompleteAddCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected cobra.ShellCompDirective
	}{
		{
			name:     "first argument (template name)",
			args:     []string{},
			expected: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name:     "second argument (template path)",
			args:     []string{"my-template"},
			expected: cobra.ShellCompDirectiveFilterDirs,
		},
		{
			name:     "third argument (no more args)",
			args:     []string{"my-template", "/path/to/template"},
			expected: cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completions, directive := completeAddCommand(nil, tt.args, "")

			// Should always return no completions for add command
			if len(completions) != 0 {
				t.Errorf("Expected no completions, got %d", len(completions))
			}

			if directive != tt.expected {
				t.Errorf("Expected directive %v, got %v", tt.expected, directive)
			}
		})
	}
}

func TestSetupCompletions(t *testing.T) {
	// Test that setupCompletions doesn't panic
	setupCompletions()

	// Test that completion functions are set
	if newCmd.ValidArgsFunction == nil {
		t.Error("newCmd should have ValidArgsFunction set")
	}

	if addCmd.ValidArgsFunction == nil {
		t.Error("addCmd should have ValidArgsFunction set")
	}

	if removeCmd.ValidArgsFunction == nil {
		t.Error("removeCmd should have ValidArgsFunction set")
	}

	if validateCmd.ValidArgsFunction == nil {
		t.Error("validateCmd should have ValidArgsFunction set")
	}
}
