// Package pkg provides public library APIs for the Ason project scaffolding tool.
package pkg

import (
	"context"
	"sync"
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

	// Stub implementation - will be filled in with actual logic in T015
	return nil
}
