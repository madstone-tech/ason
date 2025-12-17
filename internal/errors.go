// Package internal contains internal utility functions and types for Ason.
package internal

import "fmt"

// TemplateNotFoundError is returned when a template cannot be found.
type TemplateNotFoundError struct {
	TemplateName string
	Path         string
}

func (e *TemplateNotFoundError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("template not found: %s (path: %s)", e.TemplateName, e.Path)
	}
	return fmt.Sprintf("template not found: %s", e.TemplateName)
}

// InvalidPathError is returned when a path is invalid or unsafe.
type InvalidPathError struct {
	Path   string
	Reason string
}

func (e *InvalidPathError) Error() string {
	return fmt.Sprintf("invalid path: %s (%s)", e.Path, e.Reason)
}

// VariableValidationError is returned when a variable fails validation.
type VariableValidationError struct {
	VariableName string
	Value        interface{}
	Reason       string
}

func (e *VariableValidationError) Error() string {
	return fmt.Sprintf("invalid variable %s: %s", e.VariableName, e.Reason)
}

// GenerationError is returned when project generation fails.
type GenerationError struct {
	Phase  string // e.g., "loading", "rendering", "writing"
	Reason string
	Cause  error
}

func (e *GenerationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("generation error during %s: %s (cause: %v)", e.Phase, e.Reason, e.Cause)
	}
	return fmt.Sprintf("generation error during %s: %s", e.Phase, e.Reason)
}

func (e *GenerationError) Unwrap() error {
	return e.Cause
}

// RegistryError is returned when registry operations fail.
type RegistryError struct {
	Operation string // e.g., "register", "list", "remove", "load"
	Name      string
	Reason    string
	Cause     error
}

func (e *RegistryError) Error() string {
	reason := e.Reason
	if e.Name != "" {
		reason = fmt.Sprintf("%s (%s)", e.Reason, e.Name)
	}
	if e.Cause != nil {
		return fmt.Sprintf("registry error during %s: %s (cause: %v)", e.Operation, reason, e.Cause)
	}
	return fmt.Sprintf("registry error during %s: %s", e.Operation, reason)
}

func (e *RegistryError) Unwrap() error {
	return e.Cause
}

// EngineError is returned when the template engine fails.
type EngineError struct {
	Operation string // e.g., "render", "compile"
	Reason    string
	Cause     error
}

func (e *EngineError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("engine error during %s: %s (cause: %v)", e.Operation, e.Reason, e.Cause)
	}
	return fmt.Sprintf("engine error during %s: %s", e.Operation, e.Reason)
}

func (e *EngineError) Unwrap() error {
	return e.Cause
}
