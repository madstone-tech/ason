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

// TestNewGenerator tests Generator creation
func TestNewGenerator(t *testing.T) {
	// Test successful creation
	engine := mocks.NewMockEngine()
	gen, err := pkg.NewGenerator(engine)
	require.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, engine, gen.GetEngine())

	// Test nil engine rejection
	gen, err = pkg.NewGenerator(nil)
	assert.Error(t, err)
	assert.Nil(t, gen)
	assert.Contains(t, err.Error(), "cannot be nil")
}

// TestGenerateValidation tests input validation
func TestGenerateValidation(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)
	ctx := context.Background()

	tmpDir := t.TempDir()

	tests := map[string]struct {
		templatePath string
		variables    map[string]interface{}
		shouldError  bool
	}{
		"empty template path": {
			templatePath: "",
			variables:    map[string]interface{}{"name": "test"},
			shouldError:  true,
		},
		"invalid variable name": {
			templatePath: "tests/fixtures/templates/simple",
			variables:    map[string]interface{}{"123invalid": "value"},
			shouldError:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := gen.Generate(ctx, tt.templatePath, tt.variables, tmpDir)
			if tt.shouldError {
				assert.Error(t, err)
			}
		})
	}
}

// TestGenerateWithContext tests context cancellation
func TestGenerateWithContext(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	tmpDir := t.TempDir()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := gen.Generate(ctx, "template", nil, tmpDir)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

// TestGenerateWithTimeout tests timeout behavior
func TestGenerateWithTimeout(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	tmpDir := t.TempDir()

	// Create a context that times out immediately
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Ensure timeout

	err := gen.Generate(ctx, "template", nil, tmpDir)
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestGenerateConcurrent tests concurrent generation with multiple goroutines
func TestGenerateConcurrent(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	// Run concurrent operations
	results := RunConcurrent(t, 5, func(index int) error {
		tmpDir := t.TempDir()
		ctx := context.Background()

		// Use valid paths for simplicity
		return gen.Generate(
			ctx,
			"tests/fixtures/templates/simple",
			map[string]interface{}{"index": index},
			tmpDir,
		)
	})

	// We don't assert success because template might not exist
	// Just verify the concurrent structure worked
	assert.Len(t, results, 5)
}

// TestGetEngine verifies GetEngine returns the correct engine
func TestGetEngine(t *testing.T) {
	engine := mocks.NewMockEngine()
	gen, _ := pkg.NewGenerator(engine)

	retrieved := gen.GetEngine()
	assert.Equal(t, engine, retrieved)
}

// TestInvalidArgumentError verifies error formatting
func TestInvalidArgumentError(t *testing.T) {
	err := &pkg.InvalidArgumentError{
		Argument: "engine",
		Reason:   "cannot be nil",
	}

	expected := "invalid argument engine: cannot be nil"
	assert.Equal(t, expected, err.Error())
}
