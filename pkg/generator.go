// Package pkg provides public library APIs for the Ason project scaffolding tool.
package pkg

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/madstone-tech/ason/internal"
)

// Engine defines the interface that template engines must implement.
// It is used to render template content with variable substitution.
type Engine interface {
	// Render renders a template string with the given context variables.
	// The template string should be in the template engine's format (e.g., Pongo2/Jinja2).
	// The context map provides variables that can be referenced in the template.
	// Returns the rendered output or an error if rendering fails.
	//
	// Parameters:
	//   - template: The template content as a string
	//   - context: Variables available for substitution in the template
	//
	// Returns:
	//   - The rendered template output
	//   - An error if rendering fails (template syntax error, missing variables, etc.)
	Render(template string, context map[string]interface{}) (string, error)

	// RenderFile renders a template from a file with the given context variables.
	// This is used for template files (e.g., *.tmpl) in the project template.
	// The filePath should be an absolute or relative path to the template file.
	// Returns the rendered output or an error if the file cannot be read or rendering fails.
	//
	// Parameters:
	//   - filePath: Path to the template file
	//   - context: Variables available for substitution in the template
	//
	// Returns:
	//   - The rendered template output
	//   - An error if the file cannot be read or rendering fails
	RenderFile(filePath string, context map[string]interface{}) (string, error)
}

// Generator provides methods to generate projects from templates.
// A Generator is created once and can be used to generate multiple projects.
// It is safe for concurrent use when handling multiple project generations.
//
// Example:
//
//	engine := pkg.NewDefaultEngine()
//	gen := pkg.NewGenerator(engine)
//	err := gen.Generate(ctx, "/path/to/template", variables, "/output/path")
type Generator struct {
	// engine is the template rendering engine (Pongo2, custom, etc.)
	engine Engine

	// mu protects concurrent access to shared state
	mu sync.RWMutex
}

// NewGenerator creates a new Generator with the specified template engine.
// The engine is used for rendering template files and strings.
// Returns an error if the engine is nil.
//
// Parameters:
//   - engine: The template engine to use for rendering. Must not be nil.
//
// Returns:
//   - A new Generator instance ready for use
//   - An error if the engine is nil
//
// Example:
//
//	engine := pkg.NewDefaultEngine()
//	gen, err := pkg.NewGenerator(engine)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewGenerator(engine Engine) (*Generator, error) {
	if engine == nil {
		return nil, &InvalidArgumentError{
			Argument: "engine",
			Reason:   "engine cannot be nil",
		}
	}

	return &Generator{
		engine: engine,
	}, nil
}

// GetEngine returns the template engine used by this generator.
// This is primarily useful for testing or inspecting the engine configuration.
func (g *Generator) GetEngine() Engine {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.engine
}

// InvalidArgumentError is returned when a required argument is invalid.
type InvalidArgumentError struct {
	Argument string
	Reason   string
}

func (e *InvalidArgumentError) Error() string {
	return "invalid argument " + e.Argument + ": " + e.Reason
}

// Generate renders a project template and writes the output to the specified directory.
// It is safe for concurrent use but respects the context for cancellation.
//
// The function validates all inputs, loads the template configuration,
// processes all template files, and writes the output to the specified directory.
// Binary files are preserved without rendering.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - templatePath: Path to the template directory (must be absolute or relative)
//   - variables: Variables available for substitution in templates
//   - outputPath: Path where the generated project will be written (must be absolute or relative)
//
// Returns:
//   - nil if generation succeeds
//   - An error if any step fails (validation, loading, rendering, or writing)
//
// The function validates:
//   - templatePath is not empty and doesn't contain directory traversal attempts
//   - outputPath is not empty and doesn't contain directory traversal attempts
//   - All variables have valid names and values
//   - Context is not cancelled
//
// Example:
//
//	gen, _ := pkg.NewGenerator(engine)
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	err := gen.Generate(ctx,
//	    "/path/to/template",
//	    map[string]interface{}{"project_name": "my-app"},
//	    "/output/path")
func (g *Generator) Generate(
	ctx context.Context,
	templatePath string,
	variables map[string]interface{},
	outputPath string,
) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Validate context is not cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Input validation (paths, variables) using validation utilities (FR-010, FR-011)
	if err := validateInputs(templatePath, outputPath, variables); err != nil {
		return err
	}

	// Render and write the template to the output directory
	if err := g.renderAndWrite(ctx, templatePath, outputPath, variables); err != nil {
		return wrapError("rendering", err)
	}

	return nil
}

// validateInputs validates all input parameters
func validateInputs(templatePath, outputPath string, variables map[string]interface{}) error {
	// Validate paths
	if templatePath == "" {
		return &internal.InvalidPathError{
			Path:   templatePath,
			Reason: "template path cannot be empty",
		}
	}

	if outputPath == "" {
		return &internal.InvalidPathError{
			Path:   outputPath,
			Reason: "output path cannot be empty",
		}
	}

	// Validate paths don't contain traversal attacks
	if err := internal.ValidatePath(templatePath); err != nil {
		return err
	}

	if err := internal.ValidatePath(outputPath); err != nil {
		return err
	}

	// Validate all variables
	if variables != nil {
		for name, value := range variables {
			if err := internal.ValidateVariableName(name); err != nil {
				return err
			}
			if err := internal.ValidateVariableValue(name, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// renderAndWrite renders and writes the template to the output directory
func (g *Generator) renderAndWrite(
	ctx context.Context,
	templatePath string,
	outputPath string,
	variables map[string]interface{},
) error {
	// Create output directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	// Walk the template directory and process all files
	return filepath.WalkDir(templatePath, func(path string, d fs.DirEntry, err error) error {
		// Check context frequently (NFR-P-004: context cancellation support)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Skip ason.toml (config file, not template)
		if d.Name() == "ason.toml" {
			return nil
		}

		// Skip hidden files except .gitignore
		if strings.HasPrefix(d.Name(), ".") && d.Name() != ".gitignore" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		destPath := filepath.Join(outputPath, relPath)

		if d.IsDir() {
			// Create directory (FR-001: write output preserving structure)
			return os.MkdirAll(destPath, 0755)
		}

		// Process file (render template or copy binary) (AC-003: preserve binary files)
		return g.processFileForGeneration(ctx, path, destPath, variables)
	})
}

// processFileForGeneration processes a single file during generation
func (g *Generator) processFileForGeneration(
	ctx context.Context,
	srcPath string,
	destPath string,
	variables map[string]interface{},
) error {
	// Check context
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Check if file should be templated (binary files are copied as-is)
	ext := strings.ToLower(filepath.Ext(srcPath))
	binaryExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true,
		".pdf": true, ".zip": true, ".exe": true, ".bin": true, ".so": true,
		".dylib": true, ".dll": true, ".woff": true, ".woff2": true, ".ttf": true,
	}

	if binaryExts[ext] {
		// Copy binary files as-is (AC-003: binary file preservation)
		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	}

	// Render text files (FR-001: template rendering)
	srcContent, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	rendered, err := g.engine.Render(string(srcContent), variables)
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, []byte(rendered), 0644)
}

// wrapError wraps an error with generation context
func wrapError(phase string, err error) error {
	return &internal.GenerationError{
		Phase:  phase,
		Reason: err.Error(),
		Cause:  err,
	}
}
