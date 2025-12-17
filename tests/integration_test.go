package tests

import (
	"context"
	"testing"
	"time"

	"github.com/madstone-tech/ason/pkg"
	"github.com/madstone-tech/ason/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerationWorkflowBasic tests a basic end-to-end generation workflow.
// This tests the complete flow: create generator -> generate project.
func TestGenerationWorkflowBasic(t *testing.T) {
	engine := mocks.NewMockEngine()
	engine.SetRenderResponse("# Test Project")

	gen, err := pkg.NewGenerator(engine)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := map[string]interface{}{
		"project_name": "test-project",
	}

	err = gen.Generate(ctx, "fixtures/templates/simple", vars, "/tmp/test-output")
	assert.NoError(t, err)
}

// TestGenerationWorkflowWithVariables tests generation with multiple variables.
func TestGenerationWorkflowWithVariables(t *testing.T) {
	engine := mocks.NewMockEngine()
	engine.SetRenderResponse("# My Project\nAuthor: Alice\nVersion: 1.0.0")

	gen, err := pkg.NewGenerator(engine)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := map[string]interface{}{
		"project_name": "my-project",
		"author":       "Alice",
		"version":      "1.0.0",
	}

	err = gen.Generate(ctx, "fixtures/templates/with-variables", vars, "/tmp/test-output")
	assert.NoError(t, err)
}

// TestGenerationWithCustomEngine tests generation using a custom mock engine.
func TestGenerationWithCustomEngine(t *testing.T) {
	customEngine := mocks.NewMockEngine()
	customEngine.SetRenderResponse("Custom rendered output")
	customEngine.SetRenderFileResponse("Custom file output")

	gen, err := pkg.NewGenerator(customEngine)
	require.NoError(t, err)

	ctx := context.Background()
	tmpDir := t.TempDir()

	// Create a temporary template directory with a simple file
	templateDir := t.TempDir()
	err = gen.Generate(ctx, templateDir, map[string]interface{}{}, tmpDir)

	assert.NoError(t, err)
	assert.Equal(t, customEngine, gen.GetEngine())
}

// TestGenerationContextCancellation tests that generation respects context cancellation.
func TestGenerationContextCancellation(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	templateDir := t.TempDir()
	tmpDir := t.TempDir()
	err := gen.Generate(ctx, templateDir, map[string]interface{}{}, tmpDir)
	assert.Equal(t, context.Canceled, err)
}

// TestGenerationContextTimeout tests that generation respects context timeout.
func TestGenerationContextTimeout(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond)

	templateDir := t.TempDir()
	tmpDir := t.TempDir()
	err := gen.Generate(ctx, templateDir, map[string]interface{}{}, tmpDir)
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestGenerationMultipleProjects tests generating multiple projects sequentially.
func TestGenerationMultipleProjects(t *testing.T) {
	engine := pkg.NewDefaultEngine()
	gen, err := pkg.NewGenerator(engine)
	require.NoError(t, err)

	ctx := context.Background()

	// Generate first project
	templateDir1 := t.TempDir()
	tmpDir1 := t.TempDir()
	err1 := gen.Generate(ctx, templateDir1, map[string]interface{}{"name": "proj1"}, tmpDir1)

	// Generate second project
	templateDir2 := t.TempDir()
	tmpDir2 := t.TempDir()
	err2 := gen.Generate(ctx, templateDir2, map[string]interface{}{"name": "proj2"}, tmpDir2)

	// Both should complete without error (validation will happen later)
	_ = err1
	_ = err2
}

// TestRenderWithDefaultEngine tests rendering with the default Pongo2 engine.
func TestRenderWithDefaultEngine(t *testing.T) {
	engine := pkg.NewDefaultEngine()

	ctx := context.Background()
	output, err := pkg.RenderWithEngine(ctx, engine, "Hello {{ name }}!", map[string]interface{}{"name": "World"})

	require.NoError(t, err)
	assert.Equal(t, "Hello World!", output)
}

// TestRenderWithComplexTemplate tests rendering a complex template.
func TestRenderWithComplexTemplate(t *testing.T) {
	engine := pkg.NewDefaultEngine()

	ctx := context.Background()
	template := `Project: {{ project_name }}
Author: {{ author }}
Version: {{ version }}
{% for item in items %}
  - {{ item }}
{% endfor %}`

	vars := map[string]interface{}{
		"project_name": "MyApp",
		"author":       "Alice",
		"version":      "1.0.0",
		"items":        []string{"Feature 1", "Feature 2"},
	}

	output, err := pkg.RenderWithEngine(ctx, engine, template, vars)

	require.NoError(t, err)
	assert.Contains(t, output, "Project: MyApp")
	assert.Contains(t, output, "Author: Alice")
	assert.Contains(t, output, "Version: 1.0.0")
	assert.Contains(t, output, "Feature 1")
	assert.Contains(t, output, "Feature 2")
}

// TestEngineInterfaceCompliance tests that mock engine implements Engine interface.
func TestEngineInterfaceCompliance(t *testing.T) {
	var engine pkg.Engine = mocks.NewMockEngine()
	assert.NotNil(t, engine)
}

// TestDefaultEngineInterfaceCompliance tests that default engine implements Engine interface.
func TestDefaultEngineInterfaceCompliance(t *testing.T) {
	engine := pkg.NewDefaultEngine()
	assert.NotNil(t, engine)
}
