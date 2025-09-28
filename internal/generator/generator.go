package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/madstone-tech/ason/internal/engine"
	"github.com/madstone-tech/ason/internal/template"
)

// Generator handles template generation
type Generator struct {
	template *Template
	engine   engine.Engine
}

// Options for generation
type Options struct {
	SkipHooks bool
	DryRun    bool
	Verbose   bool
}

// Template represents a template with its configuration
type Template struct {
	Path   string
	Config *template.Config
}

// New creates a new generator
func New(tmpl *Template, eng engine.Engine) *Generator {
	return &Generator{
		template: tmpl,
		engine:   eng,
	}
}

// Generate generates a project from the template
func (g *Generator) Generate(outputPath string, context map[string]interface{}, opts Options) error {
	if opts.DryRun {
		fmt.Printf("DRY RUN: Would generate project at %s\n", outputPath)
		if err := g.walkTemplateFiles(g.template.Path, outputPath, context, true); err != nil {
			return err
		}
		return nil
	}

	// Create output directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("※ Generating project at %s...\n", outputPath)

	// Process all template files
	if err := g.walkTemplateFiles(g.template.Path, outputPath, context, false); err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	return nil
}

// walkTemplateFiles recursively processes all files in the template
func (g *Generator) walkTemplateFiles(templatePath, outputPath string, context map[string]interface{}, dryRun bool) error {
	return filepath.Walk(templatePath, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from template root
		relPath, err := filepath.Rel(templatePath, srcPath)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		// Skip the template root directory
		if relPath == "." {
			return nil
		}

		// Skip hidden files except .gitignore and .env.example
		if strings.HasPrefix(filepath.Base(srcPath), ".") &&
			filepath.Base(srcPath) != ".gitignore" &&
			filepath.Base(srcPath) != ".env.example" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Process template variables in the path
		destRelPath, err := g.processString(relPath, context)
		if err != nil {
			return fmt.Errorf("failed to process path %s: %w", relPath, err)
		}

		destPath := filepath.Join(outputPath, destRelPath)

		if dryRun {
			if info.IsDir() {
				fmt.Printf("[DRY RUN] Would create directory: %s\n", destPath)
			} else {
				fmt.Printf("[DRY RUN] Would process file: %s → %s\n", srcPath, destPath)
			}
			return nil
		}

		if info.IsDir() {
			// Create directory
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			if opts.Verbose {
				fmt.Printf("📁 Created directory: %s\n", destRelPath)
			}
		} else {
			// Process file
			if err := g.processFile(srcPath, destPath, context); err != nil {
				return fmt.Errorf("failed to process file %s: %w", srcPath, err)
			}
			fmt.Printf("💫 Transformed: %s\n", destRelPath)
		}

		return nil
	})
}

// processFile processes a single file through the template engine
func (g *Generator) processFile(srcPath, destPath string, context map[string]interface{}) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Check if file should be processed as a template
	if g.shouldProcessAsTemplate(srcPath) {
		// Read source file
		srcContent, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read source file: %w", err)
		}

		// Process through template engine
		processedContent, err := g.engine.Render(string(srcContent), context)
		if err != nil {
			return fmt.Errorf("failed to process template: %w", err)
		}

		// Write processed content
		if err := os.WriteFile(destPath, []byte(processedContent), 0644); err != nil {
			return fmt.Errorf("failed to write processed file: %w", err)
		}
	} else {
		// Copy binary files as-is
		if err := g.copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	}

	return nil
}

// shouldProcessAsTemplate determines if a file should be processed as a template
func (g *Generator) shouldProcessAsTemplate(filePath string) bool {
	// Skip binary file extensions
	ext := strings.ToLower(filepath.Ext(filePath))
	binaryExts := []string{
		".png", ".jpg", ".jpeg", ".gif", ".ico", ".pdf", ".zip", ".tar.gz",
		".exe", ".bin", ".so", ".dylib", ".dll", ".woff", ".woff2", ".ttf",
		".eot", ".mp3", ".mp4", ".avi", ".mov", ".webm", ".ogg",
	}

	for _, binExt := range binaryExts {
		if ext == binExt {
			return false
		}
	}

	return true
}

// copyFile copies a file from src to dst
func (g *Generator) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// processString processes a string through the template engine
func (g *Generator) processString(input string, context map[string]interface{}) (string, error) {
	// Only process if the string contains template syntax
	if !strings.Contains(input, "{{") {
		return input, nil
	}

	return g.engine.Render(input, context)
}

var opts Options // Make opts available to the package
