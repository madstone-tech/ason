package internal

import (
	"path/filepath"
	"strings"
)

// CleanPath returns a cleaned version of the path that is safe to use.
// It resolves . and .. elements and removes duplicate slashes.
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// CheckTraversalAttack detects and prevents directory traversal attacks.
// It returns an error if the path attempts to escape its intended directory.
// basePath is the base directory that should not be escaped.
// targetPath is the path to validate.
func CheckTraversalAttack(basePath, targetPath string) error {
	// Clean both paths
	cleanBase := filepath.Clean(basePath)
	cleanTarget := filepath.Clean(targetPath)

	// If basePath is relative, make it absolute for comparison
	if !filepath.IsAbs(cleanBase) {
		absBase, err := filepath.Abs(cleanBase)
		if err != nil {
			return &InvalidPathError{Path: basePath, Reason: "cannot determine absolute path"}
		}
		cleanBase = absBase
	}

	// If targetPath is relative, join it with basePath
	var fullTarget string
	if filepath.IsAbs(cleanTarget) {
		fullTarget = cleanTarget
	} else {
		fullTarget = filepath.Join(cleanBase, cleanTarget)
	}

	// Evaluate symlinks if they exist
	evalBase, err := filepath.EvalSymlinks(cleanBase)
	if err != nil {
		// If EvalSymlinks fails, at least try cleaning the path
		evalBase = cleanBase
	}

	evalTarget, err := filepath.EvalSymlinks(fullTarget)
	if err != nil {
		// If the target doesn't exist yet, we can't evaluate it
		// but we can still check if the intended path escapes the base
		evalTarget = fullTarget
	}

	// Check if the evaluated target is within the base directory
	rel, err := filepath.Rel(evalBase, evalTarget)
	if err != nil {
		return &InvalidPathError{Path: targetPath, Reason: "cannot compute relative path"}
	}

	// If the relative path starts with "..", it escapes the base
	if strings.HasPrefix(rel, "..") {
		return &InvalidPathError{Path: targetPath, Reason: "path escapes base directory"}
	}

	return nil
}

// SanitizeError removes sensitive information from an error message.
// It prevents exposing internal paths and system details.
func SanitizeError(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()
	sanitized := SanitizeString(msg)
	return sanitized
}

// IsPathSafe checks if a path is safe to use in file operations.
// It combines path validation and traversal attack checking.
func IsPathSafe(basePath, targetPath string) bool {
	if err := ValidatePath(targetPath); err != nil {
		return false
	}

	if basePath != "" {
		if err := CheckTraversalAttack(basePath, targetPath); err != nil {
			return false
		}
	}

	return true
}
