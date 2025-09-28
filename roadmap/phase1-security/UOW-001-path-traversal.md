# UOW-001: Fix Path Traversal Vulnerability

## Overview
**Phase**: 1 - Security
**Priority**: Critical
**Estimated Effort**: 4-6 hours
**Dependencies**: None

## Problem Description
The `generator.go:100` contains a path traversal vulnerability where malicious templates could use `../` sequences to write files outside the intended output directory.

```go
// Current vulnerable code
destPath := filepath.Join(outputPath, destRelPath)
```

This allows templates to potentially overwrite system files or write to unauthorized locations.

## Acceptance Criteria
- [ ] All file paths are validated to prevent directory traversal
- [ ] Malicious template paths are rejected with clear error messages
- [ ] Existing legitimate templates continue to work
- [ ] Unit tests cover all path traversal scenarios
- [ ] Security review confirms vulnerability is resolved

## Technical Approach

### Implementation Steps
1. Create `validatePath` function in `internal/generator/security.go`
2. Implement path cleaning and validation logic
3. Update `walkTemplateFiles` to use path validation
4. Add comprehensive error handling and logging
5. Create unit tests for attack scenarios

### Code Changes

**New File**: `internal/generator/security.go`
```go
package generator

import (
    "fmt"
    "path/filepath"
    "strings"
)

// validatePath ensures the destination path is within the allowed output directory
func validatePath(outputPath, destRelPath string) (string, error) {
    // Clean the paths to resolve any . or .. elements
    cleanOutput := filepath.Clean(outputPath)
    destPath := filepath.Join(cleanOutput, destRelPath)
    cleanDest := filepath.Clean(destPath)

    // Ensure the destination is within the output directory
    if !strings.HasPrefix(cleanDest, cleanOutput+string(filepath.Separator)) &&
       cleanDest != cleanOutput {
        return "", fmt.Errorf("invalid template path '%s': would write outside output directory", destRelPath)
    }

    return cleanDest, nil
}
```

**Update**: `internal/generator/generator.go`
```go
// In walkTemplateFiles function, replace:
destPath := filepath.Join(outputPath, destRelPath)

// With:
destPath, err := validatePath(outputPath, destRelPath)
if err != nil {
    return fmt.Errorf("path validation failed: %w", err)
}
```

### Test Cases
- Template with `../` in file paths
- Template with `..\` (Windows) in file paths
- Deeply nested `../../../../../../etc/passwd` attempts
- Legitimate nested directory structures
- Edge cases with symlinks and special characters

## Files to Modify
- `internal/generator/generator.go` (update existing logic)
- `internal/generator/security.go` (new file)
- `internal/generator/security_test.go` (new file)

## Testing Strategy
- Unit tests for `validatePath` function
- Integration tests with malicious templates
- Cross-platform testing (Windows/Unix path handling)
- Performance impact assessment

## Rollback Plan
If issues arise:
1. Revert path validation changes
2. Add temporary logging for suspicious paths
3. Implement gradual rollout with feature flag

## Security Validation
- [ ] Manual penetration testing with crafted templates
- [ ] Code review by security-conscious team member
- [ ] Automated security scanning tools validation
- [ ] Documentation of security assumptions

## Definition of Done
- Code implemented and reviewed
- All tests passing
- Security validation complete
- Documentation updated
- No performance regression
- Backward compatibility maintained