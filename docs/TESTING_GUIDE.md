# Testing the Ason Library API

This guide shows you how to test the new Ason library export features (v0.3.0).

## Quick Start

### 1. Run the Official Test Suite

```bash
cd /path/to/ason
go test ./tests -v
```

Expected output: All 27+ tests pass

### 2. Run the Interactive Demo

```bash
go run /tmp/test_library.go
```

This runs 5 test scenarios demonstrating all major features.

## Test Scenarios

### Test 1: Generator with Default Engine

Tests basic generator creation and engine retrieval.

```go
engine := pkg.NewDefaultEngine()
gen, err := pkg.NewGenerator(engine)
if err != nil {
    log.Fatal(err)
}

// Get the engine back
retrievedEngine := gen.GetEngine()
```

**What it verifies:**
- Generator creation succeeds
- Engine is properly stored
- GetEngine() returns the correct engine

---

### Test 2: Registry Management

Tests template registration, listing, and removal.

```go
// Create registry
reg, err := pkg.NewRegistry()  // Uses XDG location
// Or: reg, err := pkg.NewRegistryAt("/custom/path")

// Register a template (path must exist!)
err = reg.Register("my_template", "/path/to/template", "Description")

// List all templates (alphabetically sorted)
templates, err := reg.List()
for _, t := range templates {
    println(t.Name, t.Path, t.Created, t.Description)
}

// Remove a template (idempotent)
err = reg.Remove("my_template")
```

**What it verifies:**
- Template registration with validation
- Template listing with alphabetical sort
- Template removal (safe to remove twice)
- XDG-compliant storage
- TOML persistence

**Important Notes:**
- Template names must be valid variable names (letters, underscores, alphanumeric)
  - Valid: `my_template`, `template_1`, `asonApp`
  - Invalid: `my-template`, `123app`, `my.template`
- Paths must exist on the filesystem
- Default registry location: `~/.local/share/ason/registry.toml`

---

### Test 3: Context Cancellation

Tests that Generate() respects context cancellation.

```go
gen, _ := pkg.NewGenerator(engine)

// Create a cancelled context
ctx, cancel := context.WithCancel(context.Background())
cancel()  // Cancel immediately

err := gen.Generate(ctx, templateDir, vars, outputDir)
if err == context.Canceled {
    println("✓ Cancellation handled correctly")
}
```

**What it verifies:**
- Context cancellation is detected early
- Generate() returns context.Canceled error
- No partial files are written

---

### Test 4: Template Rendering

Tests the Pongo2 template engine rendering.

```go
engine := pkg.NewDefaultEngine()

// Simple template
output, err := engine.Render("Hello {{ name }}!", 
    map[string]interface{}{"name": "World"})
// Result: "Hello World!"

// Complex template with loops
complexTemplate := `Features:
{% for item in items %}
- {{ item }}
{% endfor %}`

output, err := engine.Render(complexTemplate,
    map[string]interface{}{
        "items": []string{"Auth", "Database", "API"},
    })
```

**Supported Pongo2 Features:**
- Variable interpolation: `{{ variable }}`
- Filters: `{{ value|upper }}`, `{{ date|date:"Y-m-d" }}`
- Loops: `{% for item in items %} ... {% endfor %}`
- Conditionals: `{% if condition %} ... {% endif %}`
- Template inheritance and includes
- Custom filters (via custom Engine implementation)

---

### Test 5: Error Handling

Tests proper error handling and validation.

```go
// Nil engine rejected
gen, err := pkg.NewGenerator(nil)
// err != nil: "invalid argument engine: engine cannot be nil"

// Invalid template path
err := gen.Generate(ctx, "/nonexistent", vars, outputDir)
// err != nil: "generation error during rendering: ..."

// Context timeout
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
defer cancel()
time.Sleep(10 * time.Millisecond)
err := gen.Generate(ctx, templateDir, vars, outputDir)
// err == context.DeadlineExceeded
```

**Error Types Returned:**
- `context.Canceled` - Operation was cancelled
- `context.DeadlineExceeded` - Context timeout
- `*InvalidArgumentError` - Invalid input (nil engine, bad variable names)
- `*InvalidPathError` - Path traversal or invalid path
- `*VariableValidationError` - Bad variable names/values
- `*GenerationError` - Template rendering failed
- `*EngineError` - Engine-specific error

All errors support `errors.Unwrap()` for error chaining.

---

## Running the Demo

A complete working example is available at `/tmp/test_library.go`:

```bash
go run /tmp/test_library.go
```

Output shows all 5 test scenarios running successfully.

---

## Testing Your Own Code

### 1. Using the Generator API

```go
package main

import (
    "context"
    "log"
    "github.com/madstone-tech/ason/pkg"
)

func main() {
    // Create engine and generator
    engine := pkg.NewDefaultEngine()
    gen, err := pkg.NewGenerator(engine)
    if err != nil {
        log.Fatal(err)
    }

    // Define variables
    variables := map[string]interface{}{
        "project_name": "my-app",
        "author":       "Alice",
        "version":      "1.0.0",
    }

    // Generate project
    ctx := context.Background()
    err = gen.Generate(ctx, "./template", variables, "./output")
    if err != nil {
        log.Fatal(err)
    }

    println("✓ Project generated successfully")
}
```

### 2. Using the Registry API

```go
package main

import (
    "log"
    "github.com/madstone-tech/ason/pkg"
)

func main() {
    // Create registry (uses ~/.local/share/ason/registry.toml)
    reg, err := pkg.NewRegistry()
    if err != nil {
        log.Fatal(err)
    }

    // Register templates
    err = reg.Register("golang_service", 
        "/Users/me/templates/golang-service", 
        "Go microservice template")
    if err != nil {
        log.Fatal(err)
    }

    // List registered templates
    templates, err := reg.List()
    if err != nil {
        log.Fatal(err)
    }

    for _, t := range templates {
        printf("%s (%s): %s\n", t.Name, t.Path, t.Description)
    }
}
```

### 3. Custom Engine Implementation

```go
package main

import (
    "log"
    "github.com/madstone-tech/ason/pkg"
)

// MyCustomEngine implements pkg.Engine interface
type MyCustomEngine struct{}

func (e *MyCustomEngine) Render(template string, 
    context map[string]interface{}) (string, error) {
    // Your custom template rendering logic
    return template, nil
}

func (e *MyCustomEngine) RenderFile(filePath string,
    context map[string]interface{}) (string, error) {
    // Your custom file rendering logic
    return "", nil
}

func main() {
    // Use custom engine
    engine := &MyCustomEngine{}
    gen, err := pkg.NewGenerator(engine)
    if err != nil {
        log.Fatal(err)
    }

    variables := map[string]interface{}{}
    ctx := context.Background()
    err = gen.Generate(ctx, "./template", variables, "./output")
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## Common Issues & Solutions

### Issue: "variable name must start with letter..."
**Solution:** Template names must be valid Go identifiers. Use underscores instead of hyphens:
- ❌ `my-template` 
- ✓ `my_template`

### Issue: "template path does not exist"
**Solution:** The registry validates that paths exist. Create the directory first:
```bash
mkdir -p /path/to/template
```

### Issue: Template rendering returns blank
**Solution:** Verify variables match template variable names:
```go
// Template uses {{ project_name }}
variables := map[string]interface{}{
    "project_name": "MyApp",  // Must match!
}
```

### Issue: Changes not persisting in Registry
**Solution:** The Registry persists to disk automatically. Check:
- File permissions at `~/.local/share/ason/registry.toml`
- Directory exists: `mkdir -p ~/.local/share/ason/`

---

## Test Coverage

The project includes comprehensive test coverage:

- **Generator Tests** (tests/test_generator.go)
  - Creation and initialization
  - Validation (empty paths, invalid variables)
  - Context cancellation and timeouts
  - Concurrent operations

- **Registry Tests** (tests/test_registry.go)
  - CRUD operations
  - XDG compliance
  - Concurrent read/write
  - TOML persistence

- **Integration Tests** (tests/integration_test.go)
  - End-to-end workflows
  - Engine compliance
  - Error handling

- **Concurrency Tests** (tests/concurrency_helpers.go)
  - Thread safety validation
  - Race condition detection
  - Concurrent operation patterns

Run all tests:
```bash
go test ./tests -v
go test ./tests -race  # Detect race conditions
```

---

## Performance Testing

The library is designed for high performance:

```bash
# Build binary
go build -o /tmp/ason .

# Run benchmarks (if available)
go test -bench=. ./...

# Profile memory usage
go test -memprofile=mem.prof ./tests
go tool pprof mem.prof

# Profile CPU usage
go test -cpuprofile=cpu.prof ./tests
go tool pprof cpu.prof
```

---

## Next Steps

1. **Try the demo**: `go run /tmp/test_library.go`
2. **Run official tests**: `go test ./tests -v`
3. **Implement custom engine** if needed
4. **Integrate into your application** using the public API
5. **Read the documentation**: `docs/api/engine_interface.md`

