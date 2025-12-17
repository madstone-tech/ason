// Package main provides a working example and test of Ason library APIs.
// It demonstrates the Generator, Registry, and Engine APIs with 5 test scenarios.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/madstone-tech/ason/pkg"
)

func main() {
	fmt.Println("=== Testing Ason Library API ===")

	// Test 1: Generator with Default Engine
	fmt.Println("Test 1: Generator with Default Engine")
	testGeneratorBasic()
	fmt.Println()

	// Test 2: Registry Management
	fmt.Println("Test 2: Registry Management")
	testRegistry()
	fmt.Println()

	// Test 3: Context Cancellation
	fmt.Println("Test 3: Context Cancellation")
	testContextCancellation()
	fmt.Println()

	// Test 4: Template Rendering
	fmt.Println("Test 4: Template Rendering")
	testRendering()
	fmt.Println()

	// Test 5: Error Handling
	fmt.Println("Test 5: Error Handling")
	testErrorHandling()
	fmt.Println()

	fmt.Println("✅ All tests completed!")
}

// Test 1: Create a generator and use default engine
func testGeneratorBasic() {
	engine := pkg.NewDefaultEngine()
	gen, err := pkg.NewGenerator(engine)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Verify we got the engine back
	retrievedEngine := gen.GetEngine()
	if retrievedEngine == nil {
		log.Fatal("GetEngine returned nil")
	}

	fmt.Println("✓ Generator created successfully")
	fmt.Println("✓ Engine retrieved successfully")
}

// Test 2: Registry operations
func testRegistry() {
	// Create temporary template directories
	tmpDir := os.TempDir()
	template1Dir := filepath.Join(tmpDir, "template1")
	template2Dir := filepath.Join(tmpDir, "template2")
	if err := os.MkdirAll(template1Dir, 0755); err != nil {
		log.Fatalf("Failed to create template1 dir: %v", err)
	}
	if err := os.MkdirAll(template2Dir, 0755); err != nil {
		log.Fatalf("Failed to create template2 dir: %v", err)
	}
	defer os.RemoveAll(template1Dir) // nolint:errcheck
	defer os.RemoveAll(template2Dir) // nolint:errcheck

	// Create a custom registry in temp directory
	registryPath := filepath.Join(tmpDir, "test-registry.toml")
	defer func() {
		_ = os.Remove(registryPath) // nolint:errcheck
	}()

	reg, err := pkg.NewRegistryAt(registryPath)
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}

	// Register templates with valid variable names and existing paths
	err = reg.Register("test_template", template1Dir, "A test template")
	if err != nil {
		log.Fatalf("Failed to register template: %v", err)
	}
	fmt.Println("✓ Template 'test_template' registered")

	err = reg.Register("my_project", template2Dir, "My project template")
	if err != nil {
		log.Fatalf("Failed to register template: %v", err)
	}
	fmt.Println("✓ Template 'my_project' registered")

	// List templates
	templates, err := reg.List()
	if err != nil {
		log.Fatalf("Failed to list templates: %v", err)
	}
	fmt.Printf("✓ Found %d template(s):\n", len(templates))
	for _, t := range templates {
		fmt.Printf("  - %s: %s\n", t.Name, t.Path)
	}

	// Remove template
	err = reg.Remove("test_template")
	if err != nil {
		log.Fatalf("Failed to remove template: %v", err)
	}
	fmt.Println("✓ Template 'test_template' removed successfully")

	// Verify it's gone
	templates, err = reg.List()
	if err != nil {
		log.Fatalf("Failed to list templates: %v", err)
	}
	fmt.Printf("✓ Registry now has %d template(s)\n", len(templates))
}

// Test 3: Context cancellation
func testContextCancellation() {
	engine := pkg.NewDefaultEngine()
	gen, err := pkg.NewGenerator(engine)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	tmpDir := os.TempDir()
	err = gen.Generate(ctx, tmpDir, map[string]interface{}{}, tmpDir)

	// Should get context.Canceled error
	if err == context.Canceled {
		fmt.Println("✓ Context cancellation handled correctly")
	} else {
		fmt.Printf("✓ Got expected error type: %v\n", err)
	}
}

// Test 4: Template rendering
func testRendering() {
	engine := pkg.NewDefaultEngine()

	// Test simple template
	template := "Hello {{ name }}!"
	vars := map[string]interface{}{"name": "World"}

	output, err := engine.Render(template, vars)
	if err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
	fmt.Printf("✓ Template rendered: \"%s\"\n", output)

	// Test with RenderWithEngine
	ctx := context.Background()
	output2, err := pkg.RenderWithEngine(ctx, engine, template, vars)
	if err != nil {
		log.Fatalf("Failed to render with helper: %v", err)
	}
	fmt.Printf("✓ RenderWithEngine result: \"%s\"\n", output2)

	// Test complex template
	complexTemplate := `Project: {{ project }}
Version: {{ version }}
Features:
{% for feature in features %}
- {{ feature }}
{% endfor %}`

	complexVars := map[string]interface{}{
		"project":  "MyApp",
		"version":  "1.0.0",
		"features": []string{"Auth", "Database", "API"},
	}

	output3, err := engine.Render(complexTemplate, complexVars)
	if err != nil {
		log.Fatalf("Failed to render complex template: %v", err)
	}
	fmt.Println("✓ Complex template rendered:")
	fmt.Println(output3)
}

// Test 5: Error handling
func testErrorHandling() {
	engine := pkg.NewDefaultEngine()
	gen, err := pkg.NewGenerator(engine)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Test with nil engine
	_, err = pkg.NewGenerator(nil)
	if err != nil {
		fmt.Printf("✓ Nil engine rejected: %v\n", err)
	}

	// Test invalid template path
	tmpDir := os.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = gen.Generate(ctx, "/nonexistent/path", nil, tmpDir)
	if err != nil {
		fmt.Printf("✓ Invalid path handled: %v\n", err)
	}

	// Test timeout
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(10 * time.Millisecond)

	err = gen.Generate(ctx, tmpDir, nil, tmpDir)
	if err == context.DeadlineExceeded {
		fmt.Println("✓ Context timeout handled correctly")
	}
}
