# UOW-002: Fix File Permission Preservation

## Overview
**Phase**: 1 - Security
**Priority**: Critical
**Estimated Effort**: 3-4 hours
**Dependencies**: None

## Problem Description
Currently, all generated files are written with fixed permissions (0644), which:
- Doesn't preserve executable permissions from template files
- May create security issues with overly permissive defaults
- Ignores the original file's intended permission model

```go
// Current problematic code in generator.go:154
if err := os.WriteFile(destPath, []byte(processedContent), 0644); err != nil {
```

## Acceptance Criteria
- [ ] Original file permissions are preserved during template processing
- [ ] Executable files remain executable in generated projects
- [ ] Configurable permission settings for different file types
- [ ] Secure defaults that don't expose sensitive information
- [ ] Cross-platform compatibility (Windows/Unix)
- [ ] Unit tests validate permission preservation

## Technical Approach

### Implementation Steps
1. Extract file permission information during template walking
2. Create permission mapping and validation logic
3. Update file writing to preserve permissions
4. Add configuration options for permission handling
5. Implement fallback for edge cases

### Code Changes

**Update**: `internal/generator/generator.go`
```go
// Add to Generator struct
type Generator struct {
    engine Engine
    preservePermissions bool
    defaultFileMode     os.FileMode
    defaultDirMode      os.FileMode
}

// Update walkTemplateFiles to capture original permissions
func (g *Generator) walkTemplateFiles(templatePath, outputPath string, context map[string]interface{}, dryRun bool) error {
    return filepath.Walk(templatePath, func(srcPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Capture original permissions
        originalMode := info.Mode()

        // ... existing logic ...

        if !info.IsDir() {
            return g.processFile(srcPath, destPath, context, originalMode, dryRun)
        }

        return g.createDirectory(destPath, originalMode, dryRun)
    })
}

// New method for file processing with permissions
func (g *Generator) processFile(srcPath, destPath string, context map[string]interface{}, originalMode os.FileMode, dryRun bool) error {
    // ... existing content processing ...

    // Determine final permissions
    fileMode := g.determineFileMode(originalMode, destPath)

    if !dryRun {
        if err := os.WriteFile(destPath, []byte(processedContent), fileMode); err != nil {
            return fmt.Errorf("failed to write file %s: %w", destPath, err)
        }
    }

    return nil
}

// Permission determination logic
func (g *Generator) determineFileMode(originalMode os.FileMode, destPath string) os.FileMode {
    if !g.preservePermissions {
        return g.defaultFileMode
    }

    // Preserve execute permissions for owner, group, other
    if originalMode&0111 != 0 {
        return originalMode & 0777
    }

    // Regular file permissions
    return originalMode & 0666
}
```

**New File**: `internal/generator/permissions.go`
```go
package generator

import (
    "os"
    "path/filepath"
    "runtime"
)

// PermissionConfig holds permission-related configuration
type PermissionConfig struct {
    PreservePermissions bool
    DefaultFileMode     os.FileMode
    DefaultDirMode      os.FileMode
    ExecutableExtensions []string
}

// DefaultPermissionConfig returns sensible defaults
func DefaultPermissionConfig() *PermissionConfig {
    return &PermissionConfig{
        PreservePermissions:  true,
        DefaultFileMode:      0644,
        DefaultDirMode:       0755,
        ExecutableExtensions: []string{".sh", ".bat", ".exe", ".py"},
    }
}

// IsExecutableFile determines if a file should be executable
func (pc *PermissionConfig) IsExecutableFile(path string) bool {
    if runtime.GOOS == "windows" {
        // On Windows, check file extension
        ext := filepath.Ext(path)
        for _, execExt := range pc.ExecutableExtensions {
            if ext == execExt {
                return true
            }
        }
        return false
    }

    // On Unix-like systems, this will be determined by original permissions
    return false
}

// SanitizePermissions ensures permissions are secure
func (pc *PermissionConfig) SanitizePermissions(mode os.FileMode, isDir bool) os.FileMode {
    if isDir {
        // Directories need execute permission to be accessible
        return mode | 0100
    }

    // Remove write permissions for group/other by default
    return mode &^ 0022
}
```

### Test Cases
- Template with executable shell scripts
- Template with various file permissions (644, 755, 600)
- Cross-platform permission handling
- Configuration override scenarios
- Binary file permission preservation

## Files to Modify
- `internal/generator/generator.go` (update file writing logic)
- `internal/generator/permissions.go` (new file)
- `internal/generator/permissions_test.go` (new file)
- `internal/generator/generator_test.go` (add permission tests)

## Configuration Options

**New Configuration Fields**:
```go
type Options struct {
    // ... existing fields ...
    PreservePermissions  bool
    DefaultFileMode      string // "0644"
    DefaultDirMode       string // "0755"
}
```

## Testing Strategy
- Unit tests for permission preservation logic
- Integration tests with various file types
- Cross-platform testing (especially Windows)
- Security testing for permission escalation
- Performance impact measurement

## Security Considerations
- Ensure no permission escalation opportunities
- Validate that templates can't set overly permissive permissions
- Handle symlinks securely
- Document permission behavior clearly

## Backward Compatibility
- Default behavior preserves permissions (opt-out)
- Legacy mode available with fixed permissions
- Configuration migration for existing users

## Definition of Done
- File permissions correctly preserved across platforms
- Configuration options implemented and tested
- Security review completed
- Documentation updated
- Performance impact acceptable
- All edge cases handled gracefully