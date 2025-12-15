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

// TestNewGeneratorWithValidEngine tests creating a generator with a valid engine.
func TestNewGeneratorWithValidEngine(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, err := pkg.NewGenerator(engine)

	require.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, engine, gen.GetEngine())
}

// TestNewGeneratorWithNilEngine tests that NewGenerator rejects nil engine.
func TestNewGeneratorWithNilEngine(t *testing.T) {
	gen, err := pkg.NewGenerator(nil)

	assert.Error(t, err)
	assert.Nil(t, gen)
	assert.Contains(t, err.Error(), "engine cannot be nil")
}

// TestNewGeneratorWithDefaultEngine tests creating a generator with the default engine.
func TestNewGeneratorWithDefaultEngine(t *testing.T) {
	engine := pkg.NewDefaultEngine()
	gen, err := pkg.NewGenerator(engine)

	require.NoError(t, err)
	assert.NotNil(t, gen)
}

// TestGeneratorGetEngine tests retrieving the engine from a generator.
func TestGeneratorGetEngine(t *testing.T) {
	mockEngine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(mockEngine)

	retrievedEngine := gen.GetEngine()
	assert.Equal(t, mockEngine, retrievedEngine)
}

// TestGenerateWithNilContext tests that Generate respects context cancellation.
func TestGenerateWithCancelledContext(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := gen.Generate(ctx, "/template", map[string]interface{}{}, "/output")

	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

// TestGenerateWithTimeoutContext tests that Generate respects context timeout.
func TestGenerateWithTimeoutContext(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give the timeout time to fire
	time.Sleep(10 * time.Millisecond)

	err := gen.Generate(ctx, "/template", map[string]interface{}{}, "/output")

	require.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestGenerateWithValidContext tests Generate with a valid context.
func TestGenerateWithValidContext(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := gen.Generate(ctx, "/template", map[string]interface{}{}, "/output")

	// For now, we expect no error since validation isn't implemented yet
	// This will change as we implement the full Generate method
	assert.NoError(t, err)
}

// TestGenerateWithEmptyTemplatePath tests that empty template path is handled.
func TestGenerateWithEmptyTemplatePath(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx := context.Background()
	err := gen.Generate(ctx, "", map[string]interface{}{}, "/output")

	// Expected to fail on empty path
	// This will be tested once validation is implemented
	// For now just checking it doesn't panic
	_ = err
}

// TestGenerateWithEmptyOutputPath tests that empty output path is handled.
func TestGenerateWithEmptyOutputPath(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx := context.Background()
	err := gen.Generate(ctx, "/template", map[string]interface{}{}, "")

	// Expected to fail on empty path
	// For now just checking it doesn't panic
	_ = err
}

// TestGenerateWithNilVariables tests that Generate handles nil variables.
func TestGenerateWithNilVariables(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx := context.Background()
	err := gen.Generate(ctx, "/template", nil, "/output")

	// Should handle nil variables gracefully
	assert.NoError(t, err)
}

// TestGenerateWithEmptyVariables tests that Generate handles empty variables.
func TestGenerateWithEmptyVariables(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	ctx := context.Background()
	err := gen.Generate(ctx, "/template", map[string]interface{}{}, "/output")

	// Should handle empty variables gracefully
	assert.NoError(t, err)
}

// TestRenderWithEngineValidInputs tests RenderWithEngine with valid inputs.
func TestRenderWithEngineValidInputs(t *testing.T) {
	engine := mocks.NewMockEngine()
	engine.SetRenderResponse("rendered output")

	ctx := context.Background()
	output, err := pkg.RenderWithEngine(ctx, engine, "template", map[string]interface{}{})

	require.NoError(t, err)
	assert.Equal(t, "rendered output", output)
}

// TestRenderWithEngineNilEngine tests RenderWithEngine with nil engine.
func TestRenderWithEngineNilEngine(t *testing.T) {
	ctx := context.Background()
	output, err := pkg.RenderWithEngine(ctx, nil, "template", map[string]interface{}{})

	assert.Error(t, err)
	assert.Empty(t, output)
	assert.Contains(t, err.Error(), "engine cannot be nil")
}

// TestRenderWithEngineCancelledContext tests RenderWithEngine with cancelled context.
func TestRenderWithEngineCancelledContext(t *testing.T) {
	engine := mocks.NewMockEngine()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	output, err := pkg.RenderWithEngine(ctx, engine, "template", map[string]interface{}{})

	assert.Error(t, err)
	assert.Empty(t, output)
	assert.Equal(t, context.Canceled, err)
}

// TestRenderWithEngineEngineError tests RenderWithEngine when engine fails.
func TestRenderWithEngineEngineError(t *testing.T) {
	engine := mocks.NewMockEngine()
	engine.SetRenderError(pkg.NewGenerationError("render", "template error"))

	ctx := context.Background()
	output, err := pkg.RenderWithEngine(ctx, engine, "template", map[string]interface{}{})

	assert.Error(t, err)
	assert.Empty(t, output)
	assert.Contains(t, err.Error(), "template error")
}

// TestDefaultEngineNotNil tests that NewDefaultEngine returns a non-nil engine.
func TestDefaultEngineNotNil(t *testing.T) {
	engine := pkg.NewDefaultEngine()
	assert.NotNil(t, engine)
}

// TestDefaultEngineCanRender tests that the default engine can render templates.
func TestDefaultEngineCanRender(t *testing.T) {
	engine := pkg.NewDefaultEngine()
	output, err := engine.Render("Hello {{ name }}", map[string]interface{}{"name": "World"})

	require.NoError(t, err)
	assert.Equal(t, "Hello World", output)
}
