# UOW-004: Remove Global State

## Overview
**Phase**: 2 - Code Quality
**Priority**: High
**Estimated Effort**: 8-10 hours
**Dependencies**: Phase 1 completion

## Problem Description
The global `opts` variable in `generator.go:214` creates tight coupling and makes the code difficult to test and maintain:

```go
// Problematic global state
var opts *Options // This creates coupling issues
```

This global state:
- Makes unit testing difficult
- Creates potential race conditions
- Reduces code reusability
- Makes dependency injection impossible

## Acceptance Criteria
- [ ] All global state removed from generator package
- [ ] Options passed through method parameters or struct fields
- [ ] Code remains fully functional
- [ ] Unit tests can run in isolation
- [ ] No breaking changes to public API
- [ ] Thread-safe operation ensured

## Technical Approach

### Implementation Steps
1. Refactor Generator struct to hold options
2. Update all functions to use struct-based options
3. Modify CLI commands to pass options explicitly
4. Update tests to work without global state
5. Verify thread safety

### Code Changes

**Update**: `internal/generator/generator.go`
```go
// Remove global opts variable and update Generator struct
type Generator struct {
    engine  Engine
    options *Options
}

// Options struct for generator configuration
type Options struct {
    Verbose     bool
    DryRun      bool
    Variables   map[string]string
    OutputPath  string
    // Add other configuration as needed
}

// NewGenerator creates a new generator with options
func NewGenerator(engine Engine, options *Options) *Generator {
    if options == nil {
        options = &Options{
            Variables: make(map[string]string),
        }
    }

    return &Generator{
        engine:  engine,
        options: options,
    }
}

// Update all methods to use g.options instead of global opts
func (g *Generator) Generate(templatePath, outputPath string, context map[string]interface{}) error {
    // Use g.options.Verbose instead of opts.Verbose
    if g.options.Verbose {
        fmt.Printf("🚀 Generating project from template: %s\n", templatePath)
    }

    return g.walkTemplateFiles(templatePath, outputPath, context, g.options.DryRun)
}

func (g *Generator) walkTemplateFiles(templatePath, outputPath string, context map[string]interface{}, dryRun bool) error {
    return filepath.Walk(templatePath, func(srcPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        srcRelPath, err := filepath.Rel(templatePath, srcPath)
        if err != nil {
            return fmt.Errorf("failed to get relative path: %w", err)
        }

        destRelPath := g.processPath(srcRelPath, context)
        destPath, err := validatePath(outputPath, destRelPath)
        if err != nil {
            return fmt.Errorf("path validation failed: %w", err)
        }

        if info.IsDir() {
            return g.createDirectory(destPath, info.Mode(), dryRun)
        }

        return g.processFile(srcPath, destPath, context, info.Mode(), dryRun)
    })
}

func (g *Generator) processPath(path string, context map[string]interface{}) string {
    // Use g.engine instead of assuming global state
    result, err := g.engine.Render(path, context)
    if err != nil {
        if g.options.Verbose {
            fmt.Printf("⚠️  Failed to process path template %s: %v\n", path, err)
        }
        return path
    }
    return result
}

func (g *Generator) createDirectory(destPath string, mode os.FileMode, dryRun bool) error {
    if g.options.Verbose {
        fmt.Printf("📁 Created directory: %s\n", filepath.Base(destPath))
    }

    if !dryRun {
        if err := os.MkdirAll(destPath, mode); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", destPath, err)
        }
    }

    return nil
}

func (g *Generator) processFile(srcPath, destPath string, context map[string]interface{}, mode os.FileMode, dryRun bool) error {
    if g.options.Verbose {
        fmt.Printf("📄 Processing file: %s\n", filepath.Base(srcPath))
    }

    content, err := os.ReadFile(srcPath)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", srcPath, err)
    }

    if isBinaryFile(content) {
        // Copy binary files without processing
        if !dryRun {
            if err := copyFile(srcPath, destPath, mode); err != nil {
                return fmt.Errorf("failed to copy binary file: %w", err)
            }
        }
        return nil
    }

    processedContent, err := g.engine.Render(string(content), context)
    if err != nil {
        return fmt.Errorf("failed to render template %s: %w", srcPath, err)
    }

    if !dryRun {
        if err := os.WriteFile(destPath, []byte(processedContent), mode); err != nil {
            return fmt.Errorf("failed to write file %s: %w", destPath, err)
        }
    }

    return nil
}
```

**Update**: `cmd/new.go`
```go
// Update newCmd to create generator with options
func (o *newOptions) run() error {
    // Create generator options from CLI options
    genOptions := &generator.Options{
        Verbose:   o.verbose,
        DryRun:    o.dryRun,
        Variables: o.variables,
    }

    // Create template engine
    engine := engine.NewPongo2Engine()

    // Create generator with options
    gen := generator.NewGenerator(engine, genOptions)

    // Get template path
    templatePath, err := o.getTemplatePath()
    if err != nil {
        return err
    }

    // Prepare context
    context := make(map[string]interface{})
    for key, value := range o.variables {
        context[key] = value
    }

    // Generate project
    return gen.Generate(templatePath, o.output, context)
}
```

**Update**: `internal/generator/generator_test.go`
```go
// Update tests to use new constructor pattern
func TestGenerator_Generate(t *testing.T) {
    tests := []struct {
        name        string
        options     *generator.Options
        templateDir string
        outputDir   string
        context     map[string]interface{}
        wantErr     bool
    }{
        {
            name: "successful generation",
            options: &generator.Options{
                Verbose: false,
                DryRun:  false,
                Variables: map[string]string{
                    "name": "test-project",
                },
            },
            templateDir: "testdata/simple-template",
            outputDir:   "testdata/output",
            context: map[string]interface{}{
                "name": "test-project",
            },
            wantErr: false,
        },
        {
            name: "verbose mode",
            options: &generator.Options{
                Verbose: true,
                DryRun:  false,
            },
            templateDir: "testdata/simple-template",
            outputDir:   "testdata/output",
            context:     map[string]interface{}{},
            wantErr:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            engine := engine.NewPongo2Engine()
            gen := generator.NewGenerator(engine, tt.options)

            err := gen.Generate(tt.templateDir, tt.outputDir, tt.context)

            if (err != nil) != tt.wantErr {
                t.Errorf("Generator.Generate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Migration Strategy
1. **Phase 1**: Add new constructor and options struct
2. **Phase 2**: Update internal methods to use struct options
3. **Phase 3**: Update CLI commands to use new pattern
4. **Phase 4**: Remove global variable
5. **Phase 5**: Update all tests

## Files to Modify
- `internal/generator/generator.go` (major refactoring)
- `internal/generator/options.go` (new file)
- `cmd/new.go` (update to use new generator pattern)
- `internal/generator/generator_test.go` (update tests)

## Testing Strategy
- Unit tests for all generator methods with different options
- Integration tests with CLI commands
- Concurrent execution tests to verify thread safety
- Backward compatibility tests
- Performance regression tests

## Thread Safety Considerations
- Generator instances are not shared between goroutines
- Each generator has its own options instance
- No shared mutable state between operations
- Safe for concurrent use with separate instances

## Backward Compatibility
- Public API remains unchanged
- Internal API changes only affect internal packages
- CLI behavior remains identical
- Configuration options preserved

## Benefits After Implementation
- **Testability**: Unit tests can run in isolation
- **Thread Safety**: No shared global state
- **Maintainability**: Clear dependency injection
- **Reusability**: Generator can be used as a library
- **Debugging**: Easier to trace option usage

## Definition of Done
- All global state removed from generator package
- Generator constructor pattern implemented
- All tests passing with new pattern
- Thread safety verified
- Performance impact negligible
- Code review completed
- Documentation updated