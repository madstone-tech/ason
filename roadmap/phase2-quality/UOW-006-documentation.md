# UOW-006: Add Package Documentation and API Documentation

## Overview
**Phase**: 2 - Code Quality
**Priority**: High
**Estimated Effort**: 6-8 hours
**Dependencies**: UOW-004, UOW-005

## Problem Description
The codebase lacks comprehensive documentation:
- Most packages missing package-level documentation comments
- Exported functions and types lack godoc comments
- No architecture decision records (ADRs)
- Missing examples for key functionality
- No API usage documentation for library consumers

## Acceptance Criteria
- [ ] All packages have comprehensive package-level documentation
- [ ] All exported functions and types have godoc comments
- [ ] Code examples included for major functionality
- [ ] Architecture decisions documented
- [ ] API usage guide created
- [ ] Documentation follows Go conventions
- [ ] `go doc` produces useful output for all packages

## Technical Approach

### Implementation Steps
1. Add package-level documentation to all packages
2. Add godoc comments to all exported symbols
3. Create architecture decision records
4. Add code examples and usage documentation
5. Verify documentation quality with tools

### Documentation Standards

**Package Documentation Format**:
```go
// Package name provides functionality for X.
//
// This package implements Y and supports Z. It is designed to be used
// in the following scenarios:
//   - Scenario 1
//   - Scenario 2
//
// Basic usage:
//   package main
//
//   import "your-module/internal/packagename"
//
//   func main() {
//       // Example usage
//   }
//
// For more advanced usage, see the examples in the package tests.
package packagename
```

### Code Changes

**Update**: `internal/generator/generator.go`
```go
// Package generator provides template processing and project generation functionality.
//
// The generator package handles the core logic of transforming template directories
// into fully-formed projects. It supports variable substitution, file processing,
// and directory structure preservation while maintaining file permissions and
// handling binary files appropriately.
//
// Key features:
//   - Template variable substitution using configurable engines
//   - Binary file detection and preservation
//   - File permission preservation
//   - Path traversal protection
//   - Dry-run functionality for preview
//
// Basic usage:
//   engine := engine.NewPongo2Engine()
//   options := &generator.Options{
//       Verbose: true,
//       DryRun:  false,
//   }
//   gen := generator.NewGenerator(engine, options)
//
//   context := map[string]interface{}{
//       "name": "my-project",
//       "version": "1.0.0",
//   }
//
//   err := gen.Generate("/path/to/template", "/path/to/output", context)
//
// For security considerations and advanced configuration, see the Options
// documentation and the examples in generator_test.go.
package generator

import (
    "fmt"
    "os"
    "path/filepath"
)

// Engine defines the interface for template rendering engines.
//
// Implementations should handle template syntax parsing, variable substitution,
// and error reporting. The interface is designed to be engine-agnostic,
// allowing for different template languages (Pongo2, text/template, etc.).
type Engine interface {
    // Render processes a template string with the given context and returns
    // the rendered result. Returns an error if template parsing or rendering fails.
    Render(template string, context map[string]interface{}) (string, error)

    // RenderFile processes a template file with the given context and returns
    // the rendered result. This is a convenience method that reads the file
    // and calls Render.
    RenderFile(filepath string, context map[string]interface{}) (string, error)
}

// Options configures the behavior of the Generator.
//
// Options control various aspects of project generation including verbosity,
// dry-run mode, variable handling, and security settings.
type Options struct {
    // Verbose enables detailed output during generation.
    // When true, the generator will print progress information for each
    // file and directory processed.
    Verbose bool

    // DryRun enables preview mode where no files are actually written.
    // Useful for testing templates and validating output before generation.
    DryRun bool

    // Variables contains template variables for substitution.
    // These variables are available in all template files and paths.
    Variables map[string]string

    // PreservePermissions controls whether original file permissions are maintained.
    // When true, generated files will have the same permissions as template files.
    // When false, default permissions (0644 for files, 0755 for directories) are used.
    PreservePermissions bool

    // MaxFileSize limits the size of files that will be processed as templates.
    // Files larger than this size are treated as binary and copied directly.
    // Default: 10MB
    MaxFileSize int64
}

// Generator handles template processing and project generation.
//
// A Generator combines a template engine with configuration options to transform
// template directories into project structures. It handles file processing,
// directory creation, permission preservation, and variable substitution.
//
// Generators are safe for concurrent use when using separate instances with
// different options. The same generator instance should not be used concurrently.
type Generator struct {
    engine  Engine
    options *Options
}

// NewGenerator creates a new Generator with the specified engine and options.
//
// The engine parameter defines how templates are processed. Common engines include
// Pongo2Engine for Django/Jinja2-style templates.
//
// If options is nil, default options will be used with sensible defaults:
//   - Verbose: false
//   - DryRun: false
//   - Variables: empty map
//   - PreservePermissions: true
//   - MaxFileSize: 10MB
//
// Example:
//   engine := engine.NewPongo2Engine()
//   options := &Options{Verbose: true}
//   gen := NewGenerator(engine, options)
func NewGenerator(engine Engine, options *Options) *Generator {
    if options == nil {
        options = &Options{
            Variables:           make(map[string]string),
            PreservePermissions: true,
            MaxFileSize:         10 << 20, // 10MB
        }
    }

    if options.Variables == nil {
        options.Variables = make(map[string]string)
    }

    return &Generator{
        engine:  engine,
        options: options,
    }
}

// Generate processes a template directory and generates a project at the output path.
//
// The templatePath should point to a directory containing template files. All files
// and subdirectories will be processed recursively. Files containing template syntax
// will be processed through the template engine, while binary files are copied directly.
//
// The outputPath specifies where the generated project should be created. The directory
// will be created if it doesn't exist. Files in existing directories may be overwritten.
//
// The context map provides variables available during template rendering. These variables
// can be referenced in both file contents and file/directory names.
//
// Returns an error if:
//   - Template directory doesn't exist or isn't readable
//   - Output directory can't be created
//   - Template processing fails
//   - File writing fails
//   - Path traversal attempts are detected
//
// Example:
//   context := map[string]interface{}{
//       "name": "my-service",
//       "port": 8080,
//   }
//   err := gen.Generate("/templates/service", "/projects/my-service", context)
func (g *Generator) Generate(templatePath, outputPath string, context map[string]interface{}) error {
    // Implementation...
}
```

**Update**: `internal/registry/registry.go`
```go
// Package registry provides local template storage and management functionality.
//
// The registry package implements a local template repository that follows the
// XDG Base Directory Specification for cross-platform compatibility. Templates
// are stored in the user's data directory and managed through a TOML metadata file.
//
// Key features:
//   - XDG-compliant storage locations
//   - TOML-based metadata management
//   - Template CRUD operations (Create, Read, Update, Delete)
//   - Path validation and normalization
//   - Atomic operations for data consistency
//
// Storage locations:
//   - Registry metadata: $XDG_DATA_HOME/ason/registry.toml
//   - Templates: $XDG_DATA_HOME/ason/templates/{template-name}/
//
// Basic usage:
//   reg, err := registry.New()
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   // Add a template
//   err = reg.Add("my-template", "/path/to/template")
//
//   // List templates
//   templates, err := reg.List()
//
//   // Get template path
//   path, err := reg.Get("my-template")
//
// The registry ensures thread-safe operations and maintains data consistency
// through atomic file operations and proper locking mechanisms.
package registry

// Registry manages a collection of project templates stored locally.
//
// The Registry provides a high-level interface for template storage and retrieval.
// It handles all filesystem operations, metadata management, and ensures data
// consistency across operations.
//
// Registry operations are thread-safe and can be used concurrently from multiple
// goroutines. All metadata updates are atomic to prevent corruption.
type Registry struct {
    dataDir    string
    metaPath   string
    templates  map[string]Template
    mu         sync.RWMutex
}

// Template represents a stored template with metadata.
//
// Templates contain the name used for registry operations and the filesystem
// path where the template content is stored. Additional metadata may be added
// in future versions.
type Template struct {
    // Name is the unique identifier for the template in the registry.
    // Names must be valid identifiers containing only alphanumeric characters,
    // hyphens, and underscores.
    Name string `toml:"name"`

    // Path is the filesystem path where the template content is stored.
    // This is typically within the registry's template directory but may
    // point to external locations for linked templates.
    Path string `toml:"path"`

    // Description provides a human-readable description of the template.
    Description string `toml:"description,omitempty"`

    // Version indicates the template version for compatibility tracking.
    Version string `toml:"version,omitempty"`

    // CreatedAt tracks when the template was added to the registry.
    CreatedAt time.Time `toml:"created_at"`

    // UpdatedAt tracks when the template was last modified.
    UpdatedAt time.Time `toml:"updated_at"`
}

// New creates a new Registry instance using XDG-compliant storage locations.
//
// The registry will be initialized in the user's data directory as defined by
// the XDG Base Directory Specification. If the registry directory doesn't exist,
// it will be created with appropriate permissions.
//
// Returns an error if:
//   - XDG data directory cannot be determined
//   - Registry directory cannot be created
//   - Existing registry metadata is corrupted
//   - Filesystem permissions prevent access
//
// Example:
//   reg, err := registry.New()
//   if err != nil {
//       return fmt.Errorf("failed to initialize registry: %w", err)
//   }
func New() (*Registry, error) {
    // Implementation...
}

// Add registers a new template in the registry.
//
// The template content at sourcePath will be copied to the registry's storage
// directory. If a template with the same name already exists, an error is returned.
//
// Parameters:
//   - name: Unique identifier for the template (must be a valid identifier)
//   - sourcePath: Path to the template directory or file to register
//
// Returns an error if:
//   - Template name is invalid or already exists
//   - Source path doesn't exist or isn't readable
//   - Template content cannot be copied
//   - Registry metadata cannot be updated
//
// Example:
//   err := reg.Add("golang-service", "/path/to/service-template")
//   if err != nil {
//       return fmt.Errorf("failed to add template: %w", err)
//   }
func (r *Registry) Add(name, sourcePath string) error {
    // Implementation...
}
```

**Create**: `docs/architecture/ADR-001-template-engine-abstraction.md`
```md
# ADR-001: Template Engine Abstraction

## Status
Accepted

## Context
Ason needs to support template processing for project generation. The initial implementation uses Pongo2 (Django/Jinja2-style templates), but future requirements may demand different template engines (Go text/template, Handlebars, etc.).

## Decision
We will implement a template engine abstraction layer through the `Engine` interface that allows pluggable template processors.

### Interface Design
```go
type Engine interface {
    Render(template string, context map[string]interface{}) (string, error)
    RenderFile(filepath string, context map[string]interface{}) (string, error)
}
```

### Implementation Strategy
- Each engine is implemented as a separate package
- Generator accepts any Engine implementation
- Default engine remains Pongo2 for backward compatibility
- Future engines can be added without breaking changes

## Consequences

### Positive
- Flexibility to support different template syntaxes
- Easy to test with mock engines
- Clean separation of concerns
- Future-proof architecture

### Negative
- Additional abstraction layer complexity
- All engines must conform to the same interface
- Error handling must be generic across engines

## Implementation Notes
- Engine instances should be stateless and thread-safe
- Context maps should use consistent key naming conventions
- Error messages should be engine-agnostic where possible
```

**Create**: `docs/API.md`
```md
# Ason API Documentation

## Overview
Ason can be used as both a CLI tool and a Go library. This document covers the library API for programmatic usage.

## Quick Start

```go
package main

import (
    "log"
    "github.com/your-org/ason/internal/engine"
    "github.com/your-org/ason/internal/generator"
)

func main() {
    // Create template engine
    eng := engine.NewPongo2Engine()

    // Configure generator options
    opts := &generator.Options{
        Verbose: true,
        DryRun:  false,
        Variables: map[string]string{
            "name": "my-project",
        },
    }

    // Create generator
    gen := generator.NewGenerator(eng, opts)

    // Generate project
    context := map[string]interface{}{
        "name": "my-project",
        "version": "1.0.0",
    }

    err := gen.Generate("/path/to/template", "/path/to/output", context)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Core Packages

### generator
The generator package provides the main project generation functionality.

#### Types
- `Generator`: Main generator struct
- `Options`: Configuration options
- `Engine`: Template engine interface

#### Functions
- `NewGenerator(engine Engine, options *Options) *Generator`
- `(g *Generator) Generate(templatePath, outputPath string, context map[string]interface{}) error`

### registry
The registry package manages local template storage.

#### Types
- `Registry`: Template registry manager
- `Template`: Template metadata

#### Functions
- `New() (*Registry, error)`
- `(r *Registry) Add(name, sourcePath string) error`
- `(r *Registry) List() ([]Template, error)`
- `(r *Registry) Get(name string) (string, error)`
- `(r *Registry) Remove(name string) error`

### engine
The engine package provides template rendering implementations.

#### Available Engines
- `Pongo2Engine`: Django/Jinja2-style templates

#### Functions
- `NewPongo2Engine() Engine`

## Examples

### Custom Template Engine
```go
type CustomEngine struct{}

func (e *CustomEngine) Render(template string, context map[string]interface{}) (string, error) {
    // Custom rendering logic
    return processedTemplate, nil
}

func (e *CustomEngine) RenderFile(filepath string, context map[string]interface{}) (string, error) {
    content, err := os.ReadFile(filepath)
    if err != nil {
        return "", err
    }
    return e.Render(string(content), context)
}
```

### Registry Management
```go
reg, err := registry.New()
if err != nil {
    log.Fatal(err)
}

// Add template from directory
err = reg.Add("web-app", "/templates/react-app")

// List all templates
templates, err := reg.List()

// Get template path
templatePath, err := reg.Get("web-app")

// Remove template
err = reg.Remove("web-app")
```
```

## Files to Create/Modify
- `internal/generator/generator.go` (add package docs and comments)
- `internal/registry/registry.go` (add package docs and comments)
- `internal/engine/engine.go` (add package docs and comments)
- `internal/template/template.go` (add package docs and comments)
- `internal/prompt/prompt.go` (add package docs and comments)
- `internal/xdg/xdg.go` (add package docs and comments)
- `docs/architecture/ADR-001-template-engine-abstraction.md` (new)
- `docs/architecture/ADR-002-registry-storage.md` (new)
- `docs/API.md` (new)
- `examples/` directory with code examples (new)

## Documentation Quality Checks
- Run `go doc` on all packages to verify output
- Use `golint` to check comment style
- Verify examples compile and run
- Check that all exported symbols are documented
- Ensure package documentation is comprehensive

## Definition of Done
- All packages have comprehensive documentation
- All exported symbols have godoc comments
- API documentation is complete and accurate
- Architecture decisions are documented
- Examples are provided and tested
- Documentation follows Go conventions
- `go doc` produces useful output for all packages