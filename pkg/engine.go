package pkg

import (
	"context"

	"github.com/madstone-tech/ason/internal"
	internalEngine "github.com/madstone-tech/ason/internal/engine"
)

// GenerationError is returned when generation fails.
// It includes the phase and reason for the failure.
func NewGenerationError(phase, reason string) error {
	return &internal.GenerationError{
		Phase:  phase,
		Reason: reason,
	}
}

// DefaultEngine returns a new default template engine (Pongo2).
// This is the recommended engine for most users.
// The Pongo2 engine supports Jinja2-like template syntax with variable substitution and filters.
//
// Returns:
//   - A new Engine instance using Pongo2 for rendering
//
// Example:
//
//	engine := pkg.NewDefaultEngine()
//	gen, _ := pkg.NewGenerator(engine)
func NewDefaultEngine() Engine {
	return internalEngine.NewPongo2Engine()
}

// CustomEngine allows creation of custom template engine implementations.
// This interface is public and can be implemented by users to support different template formats.
//
// Example custom engine implementation:
//
//	type MyEngine struct{}
//	func (e *MyEngine) Render(template string, context map[string]interface{}) (string, error) {
//	    // Custom rendering logic
//	}
//	func (e *MyEngine) RenderFile(filePath string, context map[string]interface{}) (string, error) {
//	    // Custom file rendering logic
//	}
//
//	engine := &MyEngine{}
//	gen, _ := pkg.NewGenerator(engine)
//	gen.Generate(ctx, template, vars, output)
//
// The Engine interface is simple by design to support:
// - Pongo2/Jinja2-like engines
// - Custom template syntax engines
// - Mock engines for testing
//
// All engines must be thread-safe for concurrent use with the same Generator instance.
type CustomEngine interface {
	// Render renders a template string with context variables
	Render(template string, context map[string]interface{}) (string, error)

	// RenderFile renders a template from a file path
	RenderFile(filePath string, context map[string]interface{}) (string, error)
}

// RenderWithEngine is a helper function for rendering a template using a specific engine.
// This is useful for testing or one-off renders without creating a Generator.
//
// Parameters:
//   - ctx: Context for cancellation
//   - engine: The engine to use for rendering
//   - template: The template content
//   - context: Variables for substitution
//
// Returns:
//   - The rendered output
//   - An error if rendering fails or context is cancelled
//
// Example:
//
//	engine := pkg.NewDefaultEngine()
//	output, err := pkg.RenderWithEngine(ctx, engine, "Hello {{ name }}", map[string]interface{}{"name": "World"})
func RenderWithEngine(
	ctx context.Context,
	engine Engine,
	template string,
	context map[string]interface{},
) (string, error) {
	// Check context first
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if engine == nil {
		return "", &InvalidArgumentError{
			Argument: "engine",
			Reason:   "engine cannot be nil",
		}
	}

	return engine.Render(template, context)
}
