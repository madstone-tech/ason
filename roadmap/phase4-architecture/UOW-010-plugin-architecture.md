# UOW-010: Extend Engine Interface for Plugin Architecture

## Overview
**Phase**: 4 - Architecture
**Priority**: Medium
**Estimated Effort**: 10-12 hours
**Dependencies**: Phase 2 completion

## Problem Description
Current Engine interface is limited and doesn't support extensibility for custom template engines:
- Interface only supports basic rendering operations
- No plugin discovery or loading mechanism
- No way to extend functionality without modifying core code
- Limited to built-in template engines

The existing interface restricts future extensibility and third-party integrations.

## Acceptance Criteria
- [ ] Extended Engine interface supports plugin capabilities
- [ ] Plugin discovery and loading system implemented
- [ ] Plugin validation and security checks in place
- [ ] Documentation and examples for plugin development
- [ ] Backward compatibility with existing engines maintained
- [ ] Hot-reloading of plugins during development
- [ ] Plugin marketplace foundation established

## Technical Approach

### Implementation Strategy
1. Extend the Engine interface with plugin capabilities
2. Create plugin discovery and loading system
3. Implement plugin validation and sandboxing
4. Add configuration management for plugins
5. Create documentation and development tools

### Code Changes

**Update**: `internal/engine/engine.go`
```go
package engine

import (
    "context"
    "io"
    "time"
)

// Engine defines the core template rendering interface
type Engine interface {
    // Core rendering methods
    Render(template string, context map[string]interface{}) (string, error)
    RenderFile(filepath string, context map[string]interface{}) (string, error)

    // Plugin interface extensions
    PluginCapabilities
    ConfigurableEngine
    ContextAwareEngine
}

// PluginCapabilities defines plugin-specific functionality
type PluginCapabilities interface {
    // GetInfo returns plugin metadata
    GetInfo() *PluginInfo

    // Initialize sets up the plugin with configuration
    Initialize(config map[string]interface{}) error

    // Validate checks if the plugin can handle given content
    Validate(content []byte) error

    // Cleanup releases plugin resources
    Cleanup() error

    // GetSupportedExtensions returns file extensions this engine handles
    GetSupportedExtensions() []string
}

// ConfigurableEngine allows runtime configuration changes
type ConfigurableEngine interface {
    // SetConfig updates engine configuration
    SetConfig(key string, value interface{}) error

    // GetConfig retrieves current configuration
    GetConfig(key string) (interface{}, error)

    // ListConfigKeys returns all available configuration keys
    ListConfigKeys() []string
}

// ContextAwareEngine supports advanced context operations
type ContextAwareEngine interface {
    // RenderWithContext supports cancellation and timeouts
    RenderWithContext(ctx context.Context, template string, context map[string]interface{}) (string, error)

    // RenderStream supports streaming for large templates
    RenderStream(templateReader io.Reader, contextWriter io.Writer, context map[string]interface{}) error

    // PreprocessTemplate allows template analysis before rendering
    PreprocessTemplate(template string) (*TemplateMetadata, error)
}

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Description  string            `json:"description"`
    Author       string            `json:"author"`
    License      string            `json:"license"`
    Homepage     string            `json:"homepage"`
    Dependencies []string          `json:"dependencies"`
    Capabilities []string          `json:"capabilities"`
    Config       map[string]string `json:"config"`
    MinAsonVersion string          `json:"min_ason_version"`
}

// TemplateMetadata provides information about template structure
type TemplateMetadata struct {
    Variables    []string          `json:"variables"`
    Blocks       []string          `json:"blocks"`
    Includes     []string          `json:"includes"`
    Syntax       string            `json:"syntax"`
    Complexity   int               `json:"complexity"`
    EstimatedTime time.Duration    `json:"estimated_time"`
    Warnings     []string          `json:"warnings"`
}
```

**New File**: `internal/plugin/plugin.go`
```go
package plugin

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "plugin"
    "strings"
    "sync"
    "time"

    "your-module/internal/engine"
)

// Manager handles plugin discovery, loading, and lifecycle
type Manager struct {
    plugins       map[string]*LoadedPlugin
    pluginDirs    []string
    config        *Config
    validator     *Validator
    mu            sync.RWMutex
}

// LoadedPlugin wraps a plugin with metadata and lifecycle management
type LoadedPlugin struct {
    Info       *engine.PluginInfo
    Engine     engine.Engine
    PluginPath string
    LoadTime   time.Time
    Active     bool
    Config     map[string]interface{}
}

// Config configures plugin management behavior
type Config struct {
    EnableHotReload    bool
    PluginTimeout      time.Duration
    MaxPlugins         int
    TrustedAuthors     []string
    AllowedCapabilities []string
    SandboxEnabled     bool
}

// NewManager creates a new plugin manager
func NewManager(config *Config) *Manager {
    if config == nil {
        config = DefaultConfig()
    }

    return &Manager{
        plugins:    make(map[string]*LoadedPlugin),
        pluginDirs: getDefaultPluginDirs(),
        config:     config,
        validator:  NewValidator(config),
    }
}

// DefaultConfig returns sensible defaults for plugin management
func DefaultConfig() *Config {
    return &Config{
        EnableHotReload:     false,
        PluginTimeout:       30 * time.Second,
        MaxPlugins:          10,
        TrustedAuthors:      []string{},
        AllowedCapabilities: []string{"render", "validate", "preprocess"},
        SandboxEnabled:      true,
    }
}

// getDefaultPluginDirs returns default plugin search directories
func getDefaultPluginDirs() []string {
    homeDir, _ := os.UserHomeDir()
    return []string{
        "/usr/local/lib/ason/plugins",
        filepath.Join(homeDir, ".ason", "plugins"),
        "./plugins",
    }
}

// DiscoverPlugins scans configured directories for plugins
func (m *Manager) DiscoverPlugins() ([]string, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    var discovered []string

    for _, dir := range m.pluginDirs {
        if _, err := os.Stat(dir); os.IsNotExist(err) {
            continue
        }

        plugins, err := m.scanPluginDir(dir)
        if err != nil {
            return nil, fmt.Errorf("failed to scan plugin directory %s: %w", dir, err)
        }

        discovered = append(discovered, plugins...)
    }

    return discovered, nil
}

// scanPluginDir scans a directory for plugin files
func (m *Manager) scanPluginDir(dir string) ([]string, error) {
    var plugins []string

    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        // Look for .so files (compiled plugins)
        if strings.HasSuffix(path, ".so") {
            plugins = append(plugins, path)
        }

        return nil
    })

    return plugins, err
}

// LoadPlugin loads a specific plugin from path
func (m *Manager) LoadPlugin(pluginPath string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Validate plugin before loading
    if err := m.validator.ValidatePlugin(pluginPath); err != nil {
        return fmt.Errorf("plugin validation failed: %w", err)
    }

    // Load the plugin
    p, err := plugin.Open(pluginPath)
    if err != nil {
        return fmt.Errorf("failed to open plugin: %w", err)
    }

    // Look for the required symbols
    engineSym, err := p.Lookup("Engine")
    if err != nil {
        return fmt.Errorf("plugin missing Engine symbol: %w", err)
    }

    infoSym, err := p.Lookup("Info")
    if err != nil {
        return fmt.Errorf("plugin missing Info symbol: %w", err)
    }

    // Type assert to expected interfaces
    pluginEngine, ok := engineSym.(engine.Engine)
    if !ok {
        return fmt.Errorf("Engine symbol does not implement engine.Engine interface")
    }

    pluginInfo, ok := infoSym.(*engine.PluginInfo)
    if !ok {
        return fmt.Errorf("Info symbol is not *engine.PluginInfo")
    }

    // Additional validation
    if err := m.validator.ValidatePluginInfo(pluginInfo); err != nil {
        return fmt.Errorf("plugin info validation failed: %w", err)
    }

    // Initialize the plugin
    if err := pluginEngine.Initialize(make(map[string]interface{})); err != nil {
        return fmt.Errorf("plugin initialization failed: %w", err)
    }

    // Store the loaded plugin
    loadedPlugin := &LoadedPlugin{
        Info:       pluginInfo,
        Engine:     pluginEngine,
        PluginPath: pluginPath,
        LoadTime:   time.Now(),
        Active:     true,
        Config:     make(map[string]interface{}),
    }

    m.plugins[pluginInfo.Name] = loadedPlugin

    return nil
}

// GetPlugin retrieves a loaded plugin by name
func (m *Manager) GetPlugin(name string) (*LoadedPlugin, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    plugin, exists := m.plugins[name]
    if !exists {
        return nil, fmt.Errorf("plugin %s not found", name)
    }

    if !plugin.Active {
        return nil, fmt.Errorf("plugin %s is not active", name)
    }

    return plugin, nil
}

// ListPlugins returns all loaded plugins
func (m *Manager) ListPlugins() []*LoadedPlugin {
    m.mu.RLock()
    defer m.mu.RUnlock()

    plugins := make([]*LoadedPlugin, 0, len(m.plugins))
    for _, plugin := range m.plugins {
        plugins = append(plugins, plugin)
    }

    return plugins
}

// UnloadPlugin removes a plugin from memory
func (m *Manager) UnloadPlugin(name string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    plugin, exists := m.plugins[name]
    if !exists {
        return fmt.Errorf("plugin %s not found", name)
    }

    // Cleanup plugin resources
    if err := plugin.Engine.Cleanup(); err != nil {
        return fmt.Errorf("plugin cleanup failed: %w", err)
    }

    plugin.Active = false
    delete(m.plugins, name)

    return nil
}

// ConfigurePlugin updates plugin configuration
func (m *Manager) ConfigurePlugin(name string, config map[string]interface{}) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    plugin, exists := m.plugins[name]
    if !exists {
        return fmt.Errorf("plugin %s not found", name)
    }

    // Merge configuration
    for key, value := range config {
        plugin.Config[key] = value
        if err := plugin.Engine.SetConfig(key, value); err != nil {
            return fmt.Errorf("failed to set config %s: %w", key, err)
        }
    }

    return nil
}

// GetEngineForFile returns the most appropriate engine for a file
func (m *Manager) GetEngineForFile(filePath string) (engine.Engine, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    ext := filepath.Ext(filePath)

    // Check plugins for file extension support
    for _, plugin := range m.plugins {
        if !plugin.Active {
            continue
        }

        extensions := plugin.Engine.GetSupportedExtensions()
        for _, supportedExt := range extensions {
            if ext == supportedExt {
                return plugin.Engine, nil
            }
        }
    }

    return nil, fmt.Errorf("no engine found for file extension %s", ext)
}

// Shutdown gracefully shuts down all plugins
func (m *Manager) Shutdown(ctx context.Context) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    for name, plugin := range m.plugins {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := plugin.Engine.Cleanup(); err != nil {
                fmt.Printf("Warning: failed to cleanup plugin %s: %v\n", name, err)
            }
        }
    }

    m.plugins = make(map[string]*LoadedPlugin)
    return nil
}
```

**New File**: `internal/plugin/validator.go`
```go
package plugin

import (
    "crypto/sha256"
    "fmt"
    "io"
    "os"
    "strings"

    "your-module/internal/engine"
)

// Validator provides security validation for plugins
type Validator struct {
    config *Config
}

// NewValidator creates a new plugin validator
func NewValidator(config *Config) *Validator {
    return &Validator{
        config: config,
    }
}

// ValidatePlugin performs security validation on a plugin file
func (v *Validator) ValidatePlugin(pluginPath string) error {
    // Check file permissions
    if err := v.validateFilePermissions(pluginPath); err != nil {
        return fmt.Errorf("file permission validation failed: %w", err)
    }

    // Calculate file hash
    hash, err := v.calculateFileHash(pluginPath)
    if err != nil {
        return fmt.Errorf("hash calculation failed: %w", err)
    }

    // TODO: Check against known good/bad hashes
    _ = hash

    // Validate file size (prevent extremely large plugins)
    if err := v.validateFileSize(pluginPath); err != nil {
        return fmt.Errorf("file size validation failed: %w", err)
    }

    return nil
}

// ValidatePluginInfo validates plugin metadata
func (v *Validator) ValidatePluginInfo(info *engine.PluginInfo) error {
    if info.Name == "" {
        return fmt.Errorf("plugin name cannot be empty")
    }

    if info.Version == "" {
        return fmt.Errorf("plugin version cannot be empty")
    }

    // Check trusted authors if configured
    if len(v.config.TrustedAuthors) > 0 {
        trusted := false
        for _, author := range v.config.TrustedAuthors {
            if author == info.Author {
                trusted = true
                break
            }
        }
        if !trusted {
            return fmt.Errorf("plugin author %s is not trusted", info.Author)
        }
    }

    // Validate capabilities
    for _, capability := range info.Capabilities {
        if !v.isAllowedCapability(capability) {
            return fmt.Errorf("capability %s is not allowed", capability)
        }
    }

    return nil
}

// validateFilePermissions checks that plugin file has appropriate permissions
func (v *Validator) validateFilePermissions(pluginPath string) error {
    stat, err := os.Stat(pluginPath)
    if err != nil {
        return err
    }

    mode := stat.Mode()

    // Plugin should not be world-writable
    if mode&0002 != 0 {
        return fmt.Errorf("plugin file is world-writable")
    }

    // Plugin should be executable
    if mode&0111 == 0 {
        return fmt.Errorf("plugin file is not executable")
    }

    return nil
}

// calculateFileHash computes SHA256 hash of plugin file
func (v *Validator) calculateFileHash(pluginPath string) (string, error) {
    file, err := os.Open(pluginPath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// validateFileSize ensures plugin file is within reasonable size limits
func (v *Validator) validateFileSize(pluginPath string) error {
    stat, err := os.Stat(pluginPath)
    if err != nil {
        return err
    }

    // Limit plugin size to 100MB
    maxSize := int64(100 * 1024 * 1024)
    if stat.Size() > maxSize {
        return fmt.Errorf("plugin file too large: %d bytes (max %d)", stat.Size(), maxSize)
    }

    return nil
}

// isAllowedCapability checks if a capability is in the allowed list
func (v *Validator) isAllowedCapability(capability string) bool {
    for _, allowed := range v.config.AllowedCapabilities {
        if strings.EqualFold(capability, allowed) {
            return true
        }
    }
    return false
}
```

**New File**: `examples/plugins/simple-engine/main.go`
```go
// Example plugin implementation
package main

import (
    "context"
    "fmt"
    "io"
    "strings"

    "your-module/internal/engine"
)

// SimpleEngine implements a basic string replacement template engine
type SimpleEngine struct {
    config map[string]interface{}
}

// Ensure SimpleEngine implements all required interfaces
var _ engine.Engine = (*SimpleEngine)(nil)

// Plugin info - required for plugin discovery
var Info = &engine.PluginInfo{
    Name:        "simple-engine",
    Version:     "1.0.0",
    Description: "Simple string replacement template engine",
    Author:      "Ason Team",
    License:     "MIT",
    Capabilities: []string{"render", "validate"},
    Config: map[string]string{
        "prefix": "Template variable prefix (default: $)",
        "suffix": "Template variable suffix (default: $)",
    },
}

// Engine is the exported symbol that Ason will look for
var Engine = &SimpleEngine{
    config: make(map[string]interface{}),
}

// Render implements basic string replacement
func (e *SimpleEngine) Render(template string, context map[string]interface{}) (string, error) {
    result := template

    prefix := e.getConfigString("prefix", "$")
    suffix := e.getConfigString("suffix", "$")

    for key, value := range context {
        placeholder := prefix + key + suffix
        replacement := fmt.Sprintf("%v", value)
        result = strings.ReplaceAll(result, placeholder, replacement)
    }

    return result, nil
}

// RenderFile reads a file and renders it
func (e *SimpleEngine) RenderFile(filepath string, context map[string]interface{}) (string, error) {
    content, err := os.ReadFile(filepath)
    if err != nil {
        return "", err
    }

    return e.Render(string(content), context)
}

// GetInfo returns plugin metadata
func (e *SimpleEngine) GetInfo() *engine.PluginInfo {
    return Info
}

// Initialize sets up the plugin
func (e *SimpleEngine) Initialize(config map[string]interface{}) error {
    e.config = make(map[string]interface{})
    for k, v := range config {
        e.config[k] = v
    }
    return nil
}

// Validate checks if content can be processed
func (e *SimpleEngine) Validate(content []byte) error {
    // Simple validation - check for our variable syntax
    prefix := e.getConfigString("prefix", "$")
    if strings.Contains(string(content), prefix) {
        return nil
    }
    return fmt.Errorf("no template variables found")
}

// Cleanup releases resources
func (e *SimpleEngine) Cleanup() error {
    e.config = nil
    return nil
}

// GetSupportedExtensions returns file extensions this engine handles
func (e *SimpleEngine) GetSupportedExtensions() []string {
    return []string{".txt", ".md", ".cfg", ".conf"}
}

// SetConfig updates configuration
func (e *SimpleEngine) SetConfig(key string, value interface{}) error {
    e.config[key] = value
    return nil
}

// GetConfig retrieves configuration
func (e *SimpleEngine) GetConfig(key string) (interface{}, error) {
    value, exists := e.config[key]
    if !exists {
        return nil, fmt.Errorf("config key %s not found", key)
    }
    return value, nil
}

// ListConfigKeys returns available configuration keys
func (e *SimpleEngine) ListConfigKeys() []string {
    return []string{"prefix", "suffix"}
}

// RenderWithContext supports context cancellation
func (e *SimpleEngine) RenderWithContext(ctx context.Context, template string, context map[string]interface{}) (string, error) {
    select {
    case <-ctx.Done():
        return "", ctx.Err()
    default:
        return e.Render(template, context)
    }
}

// RenderStream processes streaming content
func (e *SimpleEngine) RenderStream(templateReader io.Reader, contextWriter io.Writer, context map[string]interface{}) error {
    content, err := io.ReadAll(templateReader)
    if err != nil {
        return err
    }

    result, err := e.Render(string(content), context)
    if err != nil {
        return err
    }

    _, err = contextWriter.Write([]byte(result))
    return err
}

// PreprocessTemplate analyzes template before rendering
func (e *SimpleEngine) PreprocessTemplate(template string) (*engine.TemplateMetadata, error) {
    prefix := e.getConfigString("prefix", "$")
    suffix := e.getConfigString("suffix", "$")

    var variables []string
    lines := strings.Split(template, "\n")

    for _, line := range lines {
        start := 0
        for {
            prefixIdx := strings.Index(line[start:], prefix)
            if prefixIdx == -1 {
                break
            }
            prefixIdx += start

            suffixIdx := strings.Index(line[prefixIdx+len(prefix):], suffix)
            if suffixIdx == -1 {
                break
            }
            suffixIdx += prefixIdx + len(prefix)

            variable := line[prefixIdx+len(prefix) : suffixIdx]
            variables = append(variables, variable)

            start = suffixIdx + len(suffix)
        }
    }

    return &engine.TemplateMetadata{
        Variables:  variables,
        Blocks:     []string{},
        Includes:   []string{},
        Syntax:     "simple",
        Complexity: len(variables),
    }, nil
}

// Helper function to get string config with default
func (e *SimpleEngine) getConfigString(key, defaultValue string) string {
    if value, exists := e.config[key]; exists {
        if str, ok := value.(string); ok {
            return str
        }
    }
    return defaultValue
}

// main is required for plugin compilation but not used
func main() {}
```

## Files to Create/Modify
- `internal/engine/engine.go` (extend interface)
- `internal/plugin/plugin.go` (new file)
- `internal/plugin/validator.go` (new file)
- `internal/plugin/plugin_test.go` (new file)
- `examples/plugins/simple-engine/` (new directory with example)
- `docs/plugin-development.md` (new file)
- `cmd/plugin.go` (new CLI commands for plugin management)

## Plugin Development Kit
- Template project structure for new plugins
- Build scripts and Makefile
- Testing utilities and helpers
- Documentation generator for plugin APIs

## Security Considerations
- Plugin sandboxing and resource limits
- Code signing and verification
- Capability-based security model
- Plugin isolation and cleanup

## Testing Strategy
- Unit tests for plugin loading and management
- Integration tests with example plugins
- Security tests for malicious plugins
- Performance tests for plugin overhead

## Definition of Done
- Extended Engine interface supports all plugin capabilities
- Plugin discovery and loading system works reliably
- Security validation prevents malicious plugins
- Example plugin demonstrates all features
- Documentation covers plugin development process
- CLI commands allow plugin management
- Backward compatibility maintained with existing engines