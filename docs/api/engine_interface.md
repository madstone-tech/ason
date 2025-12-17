# Engine Interface Documentation

The `Engine` interface allows you to implement custom template rendering engines for Ason. This enables support for template syntax beyond the default Pongo2/Jinja2.

## Interface Definition

```go
type Engine interface {
    // Render renders a template string with the given context variables
    Render(template string, context map[string]interface{}) (string, error)
    
    // RenderFile renders a template from a file path
    RenderFile(filePath string, context map[string]interface{}) (string, error)
}
```

## Method Signatures

### Render(template string, context map[string]interface{}) (string, error)

Renders a template provided as a string.

**Parameters:**
- `template`: The template content as a string
- `context`: Variables available for substitution in the template

**Returns:**
- The rendered output string
- An error if rendering fails (syntax error, missing variable, etc.)

**Contract:**
- Must be thread-safe for concurrent use
- Should respect context cancellation if applicable
- Should not modify the input template or context

### RenderFile(filePath string, context map[string]interface{}) (string, error)

Renders a template stored in a file.

**Parameters:**
- `filePath`: Path to the template file
- `context`: Variables available for substitution in the template

**Returns:**
- The rendered output string
- An error if the file cannot be read or rendering fails

**Contract:**
- Must handle file I/O errors (permission denied, not found, etc.)
- Must be thread-safe for concurrent use
- Should respect context cancellation if applicable

## Thread Safety

Your Engine implementation **MUST be thread-safe**. The Generator can and will use a single Engine instance across multiple concurrent goroutines without additional synchronization.

**Recommendations:**
- Use only thread-safe operations in your implementation
- Avoid shared mutable state
- Use `sync.Mutex` or `sync.RWMutex` if state must be shared
- Pre-compile templates during initialization when possible

## Context Cancellation

While the Engine interface doesn't explicitly take a `context.Context` parameter, your implementation should be mindful that it may be called within a cancellation context. If your engine performs long-running operations (like processing very large templates), consider checking for cancellation.

## Error Handling

Your Engine should return descriptive errors that include:
- What operation failed (render, compile, etc.)
- Why it failed (syntax error, missing variable, etc.)
- Original error if available for debugging

Example error patterns:
```go
fmt.Errorf("failed to render template: %w", underlyingError)
fmt.Errorf("undefined variable: %s", variableName)
fmt.Errorf("syntax error at line %d: %s", lineNum, message)
```

## Example: Custom Engine Implementation

See `examples/custom_engine.go` for a complete working example of a custom Engine implementation.

## Built-in Engines

### DefaultEngine (Pongo2)

The default engine uses Pongo2, which provides Jinja2-like template syntax:
- Variable substitution: `{{ variable_name }}`
- Filters: `{{ value | upper }}`
- Conditionals: `{% if condition %} ... {% endif %}`
- Loops: `{% for item in list %} ... {% endfor %}`

To use the default engine:
```go
engine := pkg.NewDefaultEngine()
gen, err := pkg.NewGenerator(engine)
```

## Performance Considerations

- **Pre-compilation**: If your engine supports it, compile templates during initialization rather than on each render
- **Caching**: Consider caching compiled templates if rendering the same template multiple times
- **Memory**: Be mindful of memory usage, especially for large templates
- **Concurrency**: Ensure your implementation scales well with concurrent usage

## Testing Your Engine

When testing a custom Engine:
1. Verify thread-safety with concurrent renders
2. Test error handling for invalid templates
3. Benchmark performance on typical templates
4. Verify output correctness with various input types
5. Test with the Generator's concurrency helpers

Example test pattern:
```go
func TestCustomEngine(t *testing.T) {
    engine := NewCustomEngine()
    gen, _ := pkg.NewGenerator(engine)
    
    // Test successful render
    err := gen.Generate(ctx, templatePath, vars, outputPath)
    require.NoError(t, err)
}
```

## Migration from Other Template Engines

If you're coming from a different template engine:

1. **Mustache** → Implement a MustacheEngine wrapping a Go Mustache library
2. **Handlebars** → Implement a HandlebarEngine wrapping a Go Handlebars library
3. **HTML Templates** → Implement an HTMLEngine wrapping Go's html/template
4. **Go Templates** → Implement a GoTemplateEngine wrapping Go's text/template
5. **Custom Format** → Implement your own parsing and rendering logic

See `docs/examples/custom_engine_guide.md` for detailed implementation guidelines.
