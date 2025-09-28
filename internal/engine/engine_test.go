package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPongo2Engine(t *testing.T) {
	engine := NewPongo2Engine()
	if engine == nil {
		t.Fatal("NewPongo2Engine() returned nil")
	}
}

func TestPongo2Engine_Render(t *testing.T) {
	engine := NewPongo2Engine()

	tests := []struct {
		name     string
		template string
		context  map[string]interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple template",
			template: "Hello {{ name }}!",
			context:  map[string]interface{}{"name": "World"},
			want:     "Hello World!",
			wantErr:  false,
		},
		{
			name:     "template with missing variable",
			template: "Hello {{ name }}!",
			context:  map[string]interface{}{},
			want:     "Hello !",
			wantErr:  false,
		},
		{
			name:     "template with loop",
			template: "{% for item in items %}{{ item }}{% endfor %}",
			context:  map[string]interface{}{"items": []string{"a", "b", "c"}},
			want:     "abc",
			wantErr:  false,
		},
		{
			name:     "invalid template syntax",
			template: "Hello {{ name",
			context:  map[string]interface{}{"name": "World"},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Render(tt.template, tt.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pongo2Engine.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pongo2Engine.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPongo2Engine_RenderFile(t *testing.T) {
	engine := NewPongo2Engine()

	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "ason_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		fileContent string
		context     map[string]interface{}
		want        string
		wantErr     bool
	}{
		{
			name:        "valid template file",
			fileContent: "Hello {{ name }}!",
			context:     map[string]interface{}{"name": "World"},
			want:        "Hello World!",
			wantErr:     false,
		},
		{
			name:        "template file with conditional",
			fileContent: "{% if show %}Visible{% else %}Hidden{% endif %}",
			context:     map[string]interface{}{"show": true},
			want:        "Visible",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpFile := filepath.Join(tmpDir, "test.tmpl")
			err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write temp file: %v", err)
			}

			got, err := engine.RenderFile(tmpFile, tt.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pongo2Engine.RenderFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pongo2Engine.RenderFile() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		_, err := engine.RenderFile("/non/existent/file.tmpl", map[string]interface{}{})
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})
}
