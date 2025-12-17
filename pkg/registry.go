// Package pkg provides public library APIs for the Ason project scaffolding tool.
package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/madstone-tech/ason/internal"
	"github.com/madstone-tech/ason/internal/xdg"
)

// TemplateInfo represents metadata about a registered template.
// It is returned by Registry.List() and contains all necessary information
// about a template's location and configuration.
//
// Fields:
//   - Name: The unique identifier for the template in the registry
//   - Path: The filesystem path where the template is located
//   - Created: The timestamp when the template was registered
//   - Description: Human-readable description of the template's purpose
type TemplateInfo struct {
	Name        string    `json:"name" toml:"name"`
	Path        string    `json:"path" toml:"path"`
	Created     time.Time `json:"created" toml:"created"`
	Description string    `json:"description" toml:"description"`
}

// Registry manages a local registry of templates.
// Templates are stored in the XDG Base Directory location (~/.local/share/ason/)
// and can be registered, listed, and removed programmatically.
//
// A Registry instance is thread-safe and can be shared across goroutines.
// Use NewRegistry() for the default XDG location or NewRegistryAt(path) for custom paths.
//
// Example:
//
//	registry, err := pkg.NewRegistry()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	err = registry.Register("my-template", "/path/to/template", "My template description")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	templates, err := registry.List()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, tmpl := range templates {
//	    fmt.Printf("%s: %s\n", tmpl.Name, tmpl.Path)
//	}
type Registry struct {
	path string
	mu   sync.RWMutex
	// templates holds the in-memory registry cache
	templates map[string]*TemplateInfo
}

// NewRegistry creates a new Registry using the default XDG Base Directory location.
// The registry file is stored at ~/.local/share/ason/registry.toml
// If the registry file doesn't exist, it will be created on first write.
//
// Returns:
//   - A new Registry instance ready for use
//   - An error if the XDG directory cannot be determined
//
// Example:
//
//	registry, err := pkg.NewRegistry()
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewRegistry() (*Registry, error) {
	dataHome, err := xdg.DataHome()
	if err != nil {
		return nil, err
	}
	regPath := filepath.Join(dataHome, "ason", "registry.toml")
	return NewRegistryAt(regPath)
}

// NewRegistryAt creates a new Registry at the specified path.
// This is useful for testing or using non-standard registry locations.
//
// Parameters:
//   - registryPath: The path where the registry file will be stored
//
// Returns:
//   - A new Registry instance
//   - An error if the registry cannot be loaded
//
// Example:
//
//	registry, err := pkg.NewRegistryAt("/custom/path/registry.toml")
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewRegistryAt(registryPath string) (*Registry, error) {
	if registryPath == "" {
		return nil, &internal.InvalidPathError{
			Path:   registryPath,
			Reason: "registry path cannot be empty",
		}
	}

	reg := &Registry{
		path:      registryPath,
		templates: make(map[string]*TemplateInfo),
	}

	// Load existing registry if it exists
	if err := reg.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return reg, nil
}

// Register adds or updates a template in the registry.
// If a template with the same name already exists, it is overwritten silently.
// The template path is validated before registration.
//
// Parameters:
//   - name: The unique identifier for the template (alphanumeric + underscores)
//   - templatePath: The filesystem path to the template directory
//   - description: A human-readable description of the template
//
// Returns:
//   - nil if the template is successfully registered
//   - An error if validation fails or the registry cannot be saved
//
// Example:
//
//	err := registry.Register("golang-api", "/templates/golang-api", "Go API template")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (r *Registry) Register(name string, templatePath string, description string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate inputs
	if name == "" {
		return &internal.VariableValidationError{
			VariableName: "name",
			Reason:       "template name cannot be empty",
		}
	}

	if templatePath == "" {
		return &internal.InvalidPathError{
			Path:   templatePath,
			Reason: "template path cannot be empty",
		}
	}

	// Validate name format
	if err := internal.ValidateVariableName(name); err != nil {
		return err
	}

	// Validate path
	if err := internal.ValidatePath(templatePath); err != nil {
		return err
	}

	// Check if template path exists
	info, err := os.Stat(templatePath)
	if err != nil {
		return &internal.InvalidPathError{
			Path:   templatePath,
			Reason: fmt.Sprintf("template path does not exist: %v", err),
		}
	}

	if !info.IsDir() {
		return &internal.InvalidPathError{
			Path:   templatePath,
			Reason: "template path must be a directory",
		}
	}

	// Create or update template info
	r.templates[name] = &TemplateInfo{
		Name:        name,
		Path:        templatePath,
		Created:     time.Now(),
		Description: description,
	}

	// Persist to registry file
	return r.save()
}

// List returns all templates registered in the registry.
// Templates are returned in alphabetical order by name.
// The registry is read-locked during this operation to allow concurrent reads.
//
// Returns:
//   - A slice of TemplateInfo for all registered templates
//   - An error if the registry cannot be read
//
// Example:
//
//	templates, err := registry.List()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, tmpl := range templates {
//	    fmt.Printf("%s: %s\n", tmpl.Name, tmpl.Path)
//	}
func (r *Registry) List() ([]TemplateInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]TemplateInfo, 0, len(r.templates))
	for _, tmpl := range r.templates {
		result = append(result, *tmpl)
	}

	// Sort by name for consistent results
	// (simple bubble sort for small registries)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Name > result[j].Name {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// Remove deletes a template from the registry.
// If the template doesn't exist, no error is returned (idempotent operation).
// The registry is write-locked during this operation to prevent concurrent writes.
//
// Parameters:
//   - name: The unique identifier of the template to remove
//
// Returns:
//   - nil if the template is successfully removed
//   - An error if the registry cannot be saved
//
// Example:
//
//	err := registry.Remove("golang-api")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (r *Registry) Remove(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Remove the template (idempotent - no error if doesn't exist)
	delete(r.templates, name)

	// Persist to registry file
	return r.save()
}

// load reads the registry from the filesystem.
// This is called during NewRegistryAt() initialization.
// The format is TOML, stored at the specified registry path.
func (r *Registry) load() error {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return err // Will be os.IsNotExist if file doesn't exist
	}

	// Parse TOML registry file
	type RegistryFile struct {
		Templates map[string]TemplateInfo `toml:"templates"`
	}

	var regFile RegistryFile
	if err := toml.Unmarshal(data, &regFile); err != nil {
		return &internal.GenerationError{
			Phase:  "registry_load",
			Reason: fmt.Sprintf("failed to parse registry file: %v", err),
			Cause:  err,
		}
	}

	// Load templates into memory map
	for name, tmpl := range regFile.Templates {
		info := tmpl
		r.templates[name] = &info
	}

	return nil
}

// save persists the registry to the filesystem.
// Uses atomic write pattern: write to temp file, then rename.
func (r *Registry) save() error {
	// Ensure directory exists
	dirPath := filepath.Dir(r.path)
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}

	// Prepare TOML data
	type RegistryFile struct {
		Templates map[string]TemplateInfo `toml:"templates"`
	}

	regFile := RegistryFile{
		Templates: make(map[string]TemplateInfo),
	}

	for name, tmpl := range r.templates {
		regFile.Templates[name] = *tmpl
	}

	// Marshal to TOML
	data, err := toml.Marshal(regFile)
	if err != nil {
		return &internal.GenerationError{
			Phase:  "registry_save",
			Reason: fmt.Sprintf("failed to marshal registry: %v", err),
			Cause:  err,
		}
	}

	// Write to temporary file first (atomic write)
	tmpPath := r.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return err
	}

	// Rename temp file to actual registry file (atomic)
	if err := os.Rename(tmpPath, r.path); err != nil {
		// Clean up temp file on error
		_ = os.Remove(tmpPath)
		return err
	}

	return nil
}
