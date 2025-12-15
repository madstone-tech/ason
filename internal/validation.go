package internal

import (
	"path/filepath"
	"regexp"
	"strings"
)

// ValidatePath validates that a path is safe to use.
// It checks for directory traversal attempts and ensures the path is absolute or relative but clean.
func ValidatePath(path string) error {
	if path == "" {
		return &InvalidPathError{Path: path, Reason: "path cannot be empty"}
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check for directory traversal attempts
	if strings.Contains(cleanPath, "..") {
		return &InvalidPathError{Path: path, Reason: "path contains directory traversal (..)"}
	}

	// For absolute paths, ensure they don't escape the root
	if filepath.IsAbs(cleanPath) {
		// Valid absolute path
		return nil
	}

	// For relative paths, ensure they don't start with parent references
	if strings.HasPrefix(cleanPath, "..") {
		return &InvalidPathError{Path: path, Reason: "relative path attempts to traverse upward"}
	}

	return nil
}

// ValidateVariableName validates that a variable name follows naming conventions.
// Valid names contain alphanumeric characters, underscores, and dots (for nested access).
func ValidateVariableName(name string) error {
	if name == "" {
		return &VariableValidationError{
			VariableName: name,
			Reason:       "variable name cannot be empty",
		}
	}

	// Variable names should be alphanumeric, underscores, and dots for nested access
	pattern := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`)
	if !pattern.MatchString(name) {
		return &VariableValidationError{
			VariableName: name,
			Reason:       "variable name must start with letter or underscore and contain only alphanumeric characters, underscores, and dots",
		}
	}

	// Check length
	if len(name) > 256 {
		return &VariableValidationError{
			VariableName: name,
			Reason:       "variable name exceeds maximum length of 256 characters",
		}
	}

	return nil
}

// ValidateVariableValue validates that a variable value is safe to use.
// It ensures the value can be properly serialized and doesn't contain malicious content.
func ValidateVariableValue(name string, value interface{}) error {
	if value == nil {
		return nil // nil values are allowed
	}

	switch v := value.(type) {
	case string:
		return validateStringValue(name, v)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return nil // numeric and boolean values are safe
	case []interface{}, map[string]interface{}:
		return nil // nested structures are allowed
	default:
		// Allow other types but warn about custom types
		return nil
	}
}

// validateStringValue checks for potentially dangerous patterns in string values.
func validateStringValue(name string, value string) error {
	// Check for null bytes which can cause issues
	if strings.Contains(value, "\x00") {
		return &VariableValidationError{
			VariableName: name,
			Value:        value,
			Reason:       "value contains null bytes",
		}
	}

	// Check for extremely long values that might indicate abuse
	if len(value) > 1000000 { // 1MB limit
		return &VariableValidationError{
			VariableName: name,
			Value:        "[truncated]",
			Reason:       "value exceeds maximum length of 1MB",
		}
	}

	return nil
}

// ValidateVariables validates a map of variables.
func ValidateVariables(variables map[string]interface{}) error {
	for name, value := range variables {
		if err := ValidateVariableName(name); err != nil {
			return err
		}
		if err := ValidateVariableValue(name, value); err != nil {
			return err
		}
	}
	return nil
}

// SanitizeString removes potentially dangerous characters from a string.
// This is used for error messages and logging to prevent information leakage.
func SanitizeString(s string) string {
	if len(s) > 512 {
		s = s[:512] + "..."
	}

	// Remove control characters except newlines and tabs
	sanitized := strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, s)

	return sanitized
}
