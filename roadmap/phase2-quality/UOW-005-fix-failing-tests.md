# UOW-005: Fix Failing Template Registry Test

## Overview
**Phase**: 2 - Code Quality
**Priority**: High
**Estimated Effort**: 4-6 hours
**Dependencies**: UOW-004 (global state removal)

## Problem Description
The failing test `TestNewCmdWithExistingTemplate` indicates registry state management issues:

```
FAIL: TestNewCmdWithExistingTemplate (0.00s)
new_test.go:170: newCmd execution failed: template not found: test-template
```

This suggests:
- Template registry is not properly initialized in tests
- Test isolation issues between test cases
- Registry state leaking between tests
- Incorrect test setup or teardown

## Acceptance Criteria
- [ ] All failing tests pass consistently
- [ ] Tests are properly isolated from each other
- [ ] Registry state is correctly managed in tests
- [ ] Test setup and teardown procedures are robust
- [ ] No flaky test behavior
- [ ] Clear test documentation for maintenance

## Technical Approach

### Root Cause Analysis Steps
1. Analyze the failing test in detail
2. Check registry initialization in test environment
3. Verify test isolation and cleanup
4. Review registry state management
5. Identify race conditions or timing issues

### Investigation Areas

**Test Setup Issues**:
```go
// Check current test setup in new_test.go
func TestNewCmdWithExistingTemplate(t *testing.T) {
    // Verify registry initialization
    // Check template existence
    // Validate test environment
}
```

**Registry State Management**:
```go
// Review registry operations in tests
func setupTestRegistry(t *testing.T) *registry.Registry {
    // Ensure clean test environment
    // Create temporary registry
    // Add test templates
}

func teardownTestRegistry(t *testing.T, reg *registry.Registry) {
    // Clean up test artifacts
    // Reset registry state
    // Remove temporary files
}
```

### Implementation Steps

1. **Diagnose the Current Failure**
   - Run the failing test in isolation
   - Add debug logging to understand the failure
   - Check registry state before and after test operations

2. **Fix Registry Test Isolation**
   - Implement proper test setup/teardown
   - Use temporary directories for test registries
   - Ensure no cross-test contamination

3. **Improve Test Reliability**
   - Add retry logic for flaky operations
   - Implement proper error handling in tests
   - Add comprehensive assertions

### Code Changes

**Update**: `cmd/new_test.go`
```go
package cmd

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "your-module/internal/registry"
    "your-module/internal/xdg"
)

// TestRegistry manages registry state for tests
type TestRegistry struct {
    tempDir    string
    registry   *registry.Registry
    originalXDG string
}

// setupTestRegistry creates an isolated test registry
func setupTestRegistry(t *testing.T) *TestRegistry {
    t.Helper()

    // Create temporary directory for test registry
    tempDir, err := os.MkdirTemp("", "ason-test-registry-*")
    require.NoError(t, err)

    // Override XDG_DATA_HOME for test isolation
    originalXDG := os.Getenv("XDG_DATA_HOME")
    os.Setenv("XDG_DATA_HOME", tempDir)

    // Create test registry
    reg, err := registry.New()
    require.NoError(t, err)

    return &TestRegistry{
        tempDir:     tempDir,
        registry:    reg,
        originalXDG: originalXDG,
    }
}

// teardownTestRegistry cleans up test registry
func (tr *TestRegistry) teardown(t *testing.T) {
    t.Helper()

    // Restore original XDG_DATA_HOME
    if tr.originalXDG != "" {
        os.Setenv("XDG_DATA_HOME", tr.originalXDG)
    } else {
        os.Unsetenv("XDG_DATA_HOME")
    }

    // Remove temporary directory
    if err := os.RemoveAll(tr.tempDir); err != nil {
        t.Logf("Warning: failed to clean up test directory %s: %v", tr.tempDir, err)
    }
}

// addTestTemplate adds a template to the test registry
func (tr *TestRegistry) addTestTemplate(t *testing.T, name, sourcePath string) {
    t.Helper()

    err := tr.registry.Add(name, sourcePath)
    require.NoError(t, err)

    // Verify template was added
    templates, err := tr.registry.List()
    require.NoError(t, err)

    found := false
    for _, template := range templates {
        if template.Name == name {
            found = true
            break
        }
    }
    require.True(t, found, "Template %s was not added to registry", name)
}

func TestNewCmdWithExistingTemplate(t *testing.T) {
    // Setup isolated test registry
    testReg := setupTestRegistry(t)
    defer testReg.teardown(t)

    // Create test template directory
    templateDir, err := os.MkdirTemp("", "test-template-*")
    require.NoError(t, err)
    defer os.RemoveAll(templateDir)

    // Create a simple test template
    templateContent := `# {{name}}
This is a test template for {{name}}.
`
    templateFile := filepath.Join(templateDir, "README.md")
    err = os.WriteFile(templateFile, []byte(templateContent), 0644)
    require.NoError(t, err)

    // Add template to registry
    testReg.addTestTemplate(t, "test-template", templateDir)

    // Create output directory
    outputDir, err := os.MkdirTemp("", "test-output-*")
    require.NoError(t, err)
    defer os.RemoveAll(outputDir)

    // Test the new command
    cmd := NewRootCmd()
    cmd.SetArgs([]string{
        "new",
        "test-template",
        outputDir,
        "--var", "name=MyProject",
    })

    err = cmd.Execute()
    assert.NoError(t, err)

    // Verify output was generated
    outputFile := filepath.Join(outputDir, "README.md")
    assert.FileExists(t, outputFile)

    // Verify template was processed correctly
    content, err := os.ReadFile(outputFile)
    require.NoError(t, err)
    assert.Contains(t, string(content), "MyProject")
    assert.NotContains(t, string(content), "{{name}}")
}

func TestNewCmdWithNonExistentTemplate(t *testing.T) {
    // Setup isolated test registry
    testReg := setupTestRegistry(t)
    defer testReg.teardown(t)

    // Create output directory
    outputDir, err := os.MkdirTemp("", "test-output-*")
    require.NoError(t, err)
    defer os.RemoveAll(outputDir)

    // Test with non-existent template
    cmd := NewRootCmd()
    cmd.SetArgs([]string{
        "new",
        "non-existent-template",
        outputDir,
    })

    err = cmd.Execute()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "template not found")
}

// Parallel test to verify no race conditions
func TestNewCmdConcurrent(t *testing.T) {
    t.Parallel()

    testReg := setupTestRegistry(t)
    defer testReg.teardown(t)

    // Create test template
    templateDir, err := os.MkdirTemp("", "concurrent-template-*")
    require.NoError(t, err)
    defer os.RemoveAll(templateDir)

    templateContent := `# {{name}}`
    templateFile := filepath.Join(templateDir, "README.md")
    err = os.WriteFile(templateFile, []byte(templateContent), 0644)
    require.NoError(t, err)

    testReg.addTestTemplate(t, "concurrent-template", templateDir)

    // Run multiple concurrent operations
    const numConcurrent = 5
    done := make(chan bool, numConcurrent)

    for i := 0; i < numConcurrent; i++ {
        go func(index int) {
            defer func() { done <- true }()

            outputDir, err := os.MkdirTemp("", fmt.Sprintf("concurrent-output-%d-*", index))
            assert.NoError(t, err)
            defer os.RemoveAll(outputDir)

            cmd := NewRootCmd()
            cmd.SetArgs([]string{
                "new",
                "concurrent-template",
                outputDir,
                "--var", fmt.Sprintf("name=Project%d", index),
            })

            err = cmd.Execute()
            assert.NoError(t, err)
        }(i)
    }

    // Wait for all goroutines to complete
    for i := 0; i < numConcurrent; i++ {
        <-done
    }
}
```

**Add**: Helper utilities in `internal/testutil/testutil.go`
```go
package testutil

import (
    "os"
    "path/filepath"
    "testing"
)

// CreateTempTemplate creates a temporary template directory for testing
func CreateTempTemplate(t *testing.T, files map[string]string) string {
    t.Helper()

    tempDir, err := os.MkdirTemp("", "ason-test-template-*")
    if err != nil {
        t.Fatalf("Failed to create temp template directory: %v", err)
    }

    for path, content := range files {
        fullPath := filepath.Join(tempDir, path)
        dir := filepath.Dir(fullPath)

        if err := os.MkdirAll(dir, 0755); err != nil {
            t.Fatalf("Failed to create directory %s: %v", dir, err)
        }

        if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
            t.Fatalf("Failed to write file %s: %v", fullPath, err)
        }
    }

    return tempDir
}

// AssertFileContent checks that a file exists and contains expected content
func AssertFileContent(t *testing.T, filePath, expectedContent string) {
    t.Helper()

    content, err := os.ReadFile(filePath)
    if err != nil {
        t.Fatalf("Failed to read file %s: %v", filePath, err)
    }

    if string(content) != expectedContent {
        t.Errorf("File content mismatch.\nExpected:\n%s\nActual:\n%s", expectedContent, string(content))
    }
}
```

## Files to Modify
- `cmd/new_test.go` (major refactoring)
- `internal/testutil/testutil.go` (new file)
- `internal/registry/registry_test.go` (review and fix)
- `cmd/commands_test.go` (ensure consistency)

## Testing Strategy
- Run tests in isolation and in parallel
- Test with different registry states
- Verify cleanup procedures work correctly
- Add stress testing for concurrent operations
- Test edge cases (empty registry, corrupted state)

## Debug Procedures
1. Add temporary debug logging to failing tests
2. Run tests with `-v` flag for verbose output
3. Use `t.TempDir()` for better temporary directory management
4. Add registry state validation between test steps

## Root Cause Prevention
- Implement test utilities for common patterns
- Add linting rules for test isolation
- Document test patterns and best practices
- Create test templates for consistent setup

## Definition of Done
- All tests pass consistently (10+ runs)
- Test execution time is reasonable (<30s for full suite)
- No flaky test behavior observed
- Test isolation verified through parallel execution
- Registry state management is robust
- Clear test documentation for future maintenance