# UOW-003: Add Input Validation and Sanitization

## Overview
**Phase**: 1 - Security
**Priority**: Critical
**Estimated Effort**: 6-8 hours
**Dependencies**: None

## Problem Description
The application lacks comprehensive input validation, creating potential security vulnerabilities:
- Direct user input usage in `commands.go:263` with `fmt.Scanln(&response)`
- No length limits on template names or paths
- Missing character restrictions on user inputs
- Potential injection attacks through template variables

## Acceptance Criteria
- [ ] All user inputs are validated and sanitized
- [ ] Template names follow strict naming conventions
- [ ] File paths are validated for safety
- [ ] Template variables are sanitized
- [ ] Input length limits are enforced
- [ ] Clear error messages for invalid inputs
- [ ] Comprehensive test coverage for validation logic

## Technical Approach

### Implementation Steps
1. Create centralized validation package
2. Implement input sanitization functions
3. Add validation to all user input points
4. Create validation middleware for CLI commands
5. Update error handling for validation failures

### Code Changes

**New File**: `internal/validation/validation.go`
```go
package validation

import (
    "fmt"
    "path/filepath"
    "regexp"
    "strings"
    "unicode"
)

const (
    MaxTemplateNameLength = 100
    MaxPathLength        = 1000
    MaxVariableNameLength = 50
    MaxVariableValueLength = 500
)

var (
    // Template name must be alphanumeric with hyphens/underscores
    templateNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
    // Variable names must be valid identifiers
    variableNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
    // Forbidden characters in file paths
    forbiddenPathChars = []string{"..", "<", ">", ":", "\"", "|", "?", "*"}
)

// ValidationError represents a validation failure
type ValidationError struct {
    Field   string
    Value   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// ValidateTemplateName validates a template name
func ValidateTemplateName(name string) error {
    if name == "" {
        return ValidationError{
            Field:   "template_name",
            Value:   name,
            Message: "template name cannot be empty",
        }
    }

    if len(name) > MaxTemplateNameLength {
        return ValidationError{
            Field:   "template_name",
            Value:   name,
            Message: fmt.Sprintf("template name too long (max %d characters)", MaxTemplateNameLength),
        }
    }

    if !templateNameRegex.MatchString(name) {
        return ValidationError{
            Field:   "template_name",
            Value:   name,
            Message: "template name must contain only alphanumeric characters, hyphens, and underscores",
        }
    }

    return nil
}

// ValidatePath validates a file or directory path
func ValidatePath(path string) error {
    if path == "" {
        return ValidationError{
            Field:   "path",
            Value:   path,
            Message: "path cannot be empty",
        }
    }

    if len(path) > MaxPathLength {
        return ValidationError{
            Field:   "path",
            Value:   path,
            Message: fmt.Sprintf("path too long (max %d characters)", MaxPathLength),
        }
    }

    // Check for forbidden characters
    for _, forbidden := range forbiddenPathChars {
        if strings.Contains(path, forbidden) {
            return ValidationError{
                Field:   "path",
                Value:   path,
                Message: fmt.Sprintf("path contains forbidden character: %s", forbidden),
            }
        }
    }

    // Ensure path is clean
    cleaned := filepath.Clean(path)
    if cleaned != path {
        return ValidationError{
            Field:   "path",
            Value:   path,
            Message: "path contains invalid sequences",
        }
    }

    return nil
}

// ValidateVariableName validates a template variable name
func ValidateVariableName(name string) error {
    if name == "" {
        return ValidationError{
            Field:   "variable_name",
            Value:   name,
            Message: "variable name cannot be empty",
        }
    }

    if len(name) > MaxVariableNameLength {
        return ValidationError{
            Field:   "variable_name",
            Value:   name,
            Message: fmt.Sprintf("variable name too long (max %d characters)", MaxVariableNameLength),
        }
    }

    if !variableNameRegex.MatchString(name) {
        return ValidationError{
            Field:   "variable_name",
            Value:   name,
            Message: "variable name must be a valid identifier",
        }
    }

    return nil
}

// SanitizeVariableValue sanitizes a template variable value
func SanitizeVariableValue(value string) (string, error) {
    if len(value) > MaxVariableValueLength {
        return "", ValidationError{
            Field:   "variable_value",
            Value:   value,
            Message: fmt.Sprintf("variable value too long (max %d characters)", MaxVariableValueLength),
        }
    }

    // Remove control characters except newlines and tabs
    sanitized := strings.Map(func(r rune) rune {
        if unicode.IsControl(r) && r != '\n' && r != '\t' {
            return -1
        }
        return r
    }, value)

    return sanitized, nil
}

// ValidateUserInput provides general user input validation
func ValidateUserInput(input string, maxLength int, allowEmpty bool) error {
    if !allowEmpty && input == "" {
        return ValidationError{
            Field:   "user_input",
            Value:   input,
            Message: "input cannot be empty",
        }
    }

    if len(input) > maxLength {
        return ValidationError{
            Field:   "user_input",
            Value:   input,
            Message: fmt.Sprintf("input too long (max %d characters)", maxLength),
        }
    }

    // Check for control characters
    for _, r := range input {
        if unicode.IsControl(r) && r != '\n' && r != '\t' {
            return ValidationError{
                Field:   "user_input",
                Value:   input,
                Message: "input contains invalid control characters",
            }
        }
    }

    return nil
}
```

**Update**: `cmd/commands.go`
```go
import (
    "bufio"
    "os"
    "strings"
    "your-module/internal/validation"
)

// Replace direct fmt.Scanln usage with validated input
func getValidatedUserInput(prompt string, maxLength int) (string, error) {
    fmt.Print(prompt)

    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
        return "", fmt.Errorf("failed to read input: %w", err)
    }

    input = strings.TrimSpace(input)

    if err := validation.ValidateUserInput(input, maxLength, false); err != nil {
        return "", err
    }

    return input, nil
}

// Update confirmation prompt
func confirmAction(message string) (bool, error) {
    for {
        response, err := getValidatedUserInput(message+" (y/N): ", 10)
        if err != nil {
            return false, err
        }

        switch strings.ToLower(response) {
        case "y", "yes":
            return true, nil
        case "n", "no", "":
            return false, nil
        default:
            fmt.Println("Please enter 'y' for yes or 'n' for no.")
        }
    }
}
```

**Update**: `cmd/new.go`
```go
// Add validation to template and output path arguments
func (o *newOptions) Validate() error {
    if err := validation.ValidateTemplateName(o.template); err != nil {
        return fmt.Errorf("invalid template name: %w", err)
    }

    if o.output != "" {
        if err := validation.ValidatePath(o.output); err != nil {
            return fmt.Errorf("invalid output path: %w", err)
        }
    }

    // Validate variable assignments
    for key, value := range o.variables {
        if err := validation.ValidateVariableName(key); err != nil {
            return fmt.Errorf("invalid variable name '%s': %w", key, err)
        }

        sanitized, err := validation.SanitizeVariableValue(value)
        if err != nil {
            return fmt.Errorf("invalid variable value for '%s': %w", key, err)
        }
        o.variables[key] = sanitized
    }

    return nil
}
```

### Input Validation Points
1. **Template Names**: CLI arguments and registry operations
2. **File Paths**: Template and output paths
3. **Variable Names**: Template variable keys
4. **Variable Values**: Template variable values
5. **User Prompts**: Interactive confirmation inputs
6. **Registry Names**: Template registry names

## Files to Modify
- `internal/validation/validation.go` (new file)
- `internal/validation/validation_test.go` (new file)
- `cmd/commands.go` (update input handling)
- `cmd/new.go` (add validation calls)
- `cmd/add.go` (add template name validation)
- `cmd/remove.go` (add template name validation)
- `internal/prompt/prompt.go` (add input validation)

## Testing Strategy
- Unit tests for all validation functions
- Edge case testing (empty, too long, special characters)
- Security testing with malicious inputs
- Integration tests with CLI commands
- Cross-platform character encoding tests

## Security Considerations
- Prevent injection attacks through template variables
- Limit resource consumption with length restrictions
- Handle Unicode and encoding edge cases
- Validate all user-controlled inputs
- Sanitize outputs to prevent information leakage

## Error Handling
- Consistent error message format
- Clear guidance for users on fixing validation errors
- Logging of validation failures for security monitoring
- Graceful degradation for non-critical validation failures

## Definition of Done
- All user inputs validated and sanitized
- Comprehensive test coverage (>95%)
- Security review completed
- Error messages user-friendly and informative
- Performance impact negligible
- Documentation updated with validation rules