package registry

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/madstone-tech/ason/internal/xdg"
)

// Registry manages local templates
type Registry struct {
	path string
}

// TemplateEntry represents a template in the registry
type TemplateEntry struct {
	Name        string    `json:"name" toml:"name"`
	Path        string    `json:"path" toml:"path"`
	Description string    `json:"description" toml:"description"`
	Source      string    `json:"source" toml:"source"`
	Type        string    `json:"type" toml:"type"`
	Size        int64     `json:"size" toml:"size"`
	Files       int       `json:"files" toml:"files"`
	Added       time.Time `json:"added" toml:"added"`
	Variables   []string  `json:"variables,omitempty" toml:"variables,omitempty"`
}

// TemplateConfig represents the ason.toml configuration
type TemplateConfig struct {
	Name        string             `toml:"name,omitempty"`
	Description string             `toml:"description,omitempty"`
	Version     string             `toml:"version,omitempty"`
	Author      string             `toml:"author,omitempty"`
	Type        string             `toml:"type,omitempty"`
	Variables   []TemplateVariable `toml:"variables,omitempty"`
	Ignore      []string           `toml:"ignore,omitempty"`
	Tags        []string           `toml:"tags,omitempty"`
}

// TemplateVariable represents a template variable definition
type TemplateVariable struct {
	Name        string      `toml:"name"`
	Description string      `toml:"description,omitempty"`
	Required    bool        `toml:"required,omitempty"`
	Default     interface{} `toml:"default,omitempty"`
	Type        string      `toml:"type,omitempty"`
	Options     []string    `toml:"options,omitempty"`
	Example     string      `toml:"example,omitempty"`
}

// RegistryMetadata stores registry information
type RegistryMetadata struct {
	Templates map[string]TemplateEntry `json:"templates" toml:"templates"`
	Updated   time.Time                `json:"updated" toml:"updated"`
}

// NewRegistry creates a new template registry
func NewRegistry() (*Registry, error) {
	registryPath, err := xdg.DataHome()
	if err != nil {
		return nil, fmt.Errorf("failed to get data directory: %w", err)
	}

	// Create registry directory if it doesn't exist
	if err := os.MkdirAll(registryPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create registry directory: %w", err)
	}

	// Create templates subdirectory
	templatesPath := filepath.Join(registryPath, "templates")
	if err := os.MkdirAll(templatesPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create templates directory: %w", err)
	}

	return &Registry{
		path: registryPath,
	}, nil
}

// List returns all templates in the registry
func (r *Registry) List() ([]TemplateEntry, error) {
	meta, err := r.loadMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to load registry metadata: %w", err)
	}

	var templates []TemplateEntry
	for _, tmpl := range meta.Templates {
		templates = append(templates, tmpl)
	}

	return templates, nil
}

// Get returns the path to a template
func (r *Registry) Get(name string) (string, error) {
	meta, err := r.loadMetadata()
	if err != nil {
		return "", fmt.Errorf("failed to load registry metadata: %w", err)
	}

	if tmpl, exists := meta.Templates[name]; exists {
		return tmpl.Path, nil
	}

	return "", fmt.Errorf("template %s not found", name)
}

// Add adds a template to the registry
func (r *Registry) Add(name, sourcePath, description, templateType string) error {
	// Validate source path exists
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("source path does not exist: %s", sourcePath)
	}

	if !info.IsDir() {
		return fmt.Errorf("source path must be a directory: %s", sourcePath)
	}

	// Load existing metadata
	meta, err := r.loadMetadata()
	if err != nil {
		return fmt.Errorf("failed to load registry metadata: %w", err)
	}

	// Check if template already exists
	if _, exists := meta.Templates[name]; exists {
		return fmt.Errorf("template %s already exists", name)
	}

	// Calculate destination path
	destPath := filepath.Join(r.path, "templates", name)

	// Copy template to registry
	if err := r.copyTemplate(sourcePath, destPath); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Analyze template
	size, files, err := r.analyzeTemplate(destPath)
	if err != nil {
		return fmt.Errorf("failed to analyze template: %w", err)
	}

	// Load template config if exists
	config, err := r.loadTemplateConfig(destPath)
	if err != nil {
		// Not an error if config doesn't exist
		config = &TemplateConfig{}
	}

	// Use config values if not provided
	if description == "" && config.Description != "" {
		description = config.Description
	}
	if templateType == "" && config.Type != "" {
		templateType = config.Type
	}

	// Extract variable names from config
	var variables []string
	for _, v := range config.Variables {
		variables = append(variables, v.Name)
	}

	// Create template entry
	tmpl := TemplateEntry{
		Name:        name,
		Path:        destPath,
		Description: description,
		Source:      sourcePath,
		Type:        templateType,
		Size:        size,
		Files:       files,
		Added:       time.Now(),
		Variables:   variables,
	}

	// Add to metadata
	meta.Templates[name] = tmpl
	meta.Updated = time.Now()

	// Save metadata
	if err := r.saveMetadata(meta); err != nil {
		return fmt.Errorf("failed to save registry metadata: %w", err)
	}

	return nil
}

// Remove removes a template from the registry
func (r *Registry) Remove(name string, backup bool, backupDir string) error {
	// Load existing metadata
	meta, err := r.loadMetadata()
	if err != nil {
		return fmt.Errorf("failed to load registry metadata: %w", err)
	}

	// Check if template exists
	tmpl, exists := meta.Templates[name]
	if !exists {
		return fmt.Errorf("template %s not found", name)
	}

	// Create backup if requested
	if backup {
		if err := r.createBackup(tmpl, backupDir); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Remove template directory
	if err := os.RemoveAll(tmpl.Path); err != nil {
		return fmt.Errorf("failed to remove template directory: %w", err)
	}

	// Remove from metadata
	delete(meta.Templates, name)
	meta.Updated = time.Now()

	// Save metadata
	if err := r.saveMetadata(meta); err != nil {
		return fmt.Errorf("failed to save registry metadata: %w", err)
	}

	return nil
}

// loadMetadata loads the registry metadata
func (r *Registry) loadMetadata() (*RegistryMetadata, error) {
	metaPath := filepath.Join(r.path, "registry.toml")

	// If metadata doesn't exist, return empty metadata
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		return &RegistryMetadata{
			Templates: make(map[string]TemplateEntry),
			Updated:   time.Now(),
		}, nil
	}

	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var meta RegistryMetadata
	if err := toml.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata file: %w", err)
	}

	if meta.Templates == nil {
		meta.Templates = make(map[string]TemplateEntry)
	}

	return &meta, nil
}

// saveMetadata saves the registry metadata
func (r *Registry) saveMetadata(meta *RegistryMetadata) error {
	metaPath := filepath.Join(r.path, "registry.toml")

	data, err := toml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// loadTemplateConfig loads the ason.toml config from a template
func (r *Registry) loadTemplateConfig(templatePath string) (*TemplateConfig, error) {
	tomlPath := filepath.Join(templatePath, "ason.toml")
	if _, err := os.Stat(tomlPath); err != nil {
		return nil, fmt.Errorf("no ason.toml found in template")
	}

	data, err := os.ReadFile(tomlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ason.toml: %w", err)
	}

	var config TemplateConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse ason.toml: %w", err)
	}
	return &config, nil
}

// copyTemplate recursively copies a template directory
func (r *Registry) copyTemplate(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip hidden files and directories (except .gitignore, .env.example)
		if strings.HasPrefix(info.Name(), ".") && info.Name() != ".gitignore" && info.Name() != ".env.example" {
			return nil
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return r.copyFile(path, dstPath)
	})
}

// copyFile copies a single file
func (r *Registry) copyFile(src, dst string) error {
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

// analyzeTemplate analyzes a template directory
func (r *Registry) analyzeTemplate(templatePath string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}

		return nil
	})

	return totalSize, fileCount, err
}

// createBackup creates a backup of a template
func (r *Registry) createBackup(tmpl TemplateEntry, backupDir string) error {
	if backupDir == "" {
		backupDir = filepath.Join(r.path, "backups")
	}

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02-150405")
	// For now, just copy the directory (TODO: implement tar.gz compression)
	backupDirPath := filepath.Join(backupDir, fmt.Sprintf("%s-%s", tmpl.Name, timestamp))
	return r.copyTemplate(tmpl.Path, backupDirPath)
}
