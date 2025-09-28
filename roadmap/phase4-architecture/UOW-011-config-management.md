# UOW-011: Implement Centralized Configuration Management System

## Overview
**Phase**: 4 - Architecture
**Priority**: Medium
**Estimated Effort**: 8-10 hours
**Dependencies**: UOW-010 (plugin architecture)

## Problem Description
The application lacks a centralized configuration management system:
- Configuration scattered across multiple files and packages
- No unified way to handle environment variables, files, and CLI flags
- Missing configuration validation and type safety
- No support for configuration profiles or environments
- Difficult to override settings for different use cases

## Acceptance Criteria
- [ ] Centralized configuration system with multiple sources
- [ ] Support for TOML, YAML, and JSON configuration files
- [ ] Environment variable integration with prefix support
- [ ] CLI flag override capabilities
- [ ] Configuration validation with clear error messages
- [ ] Profile-based configuration (dev, prod, test)
- [ ] Hot-reloading of configuration during development
- [ ] Backward compatibility with existing behavior

## Technical Approach

### Implementation Strategy
1. Create unified configuration structure
2. Implement multi-source configuration loading
3. Add validation and type safety
4. Create profile and environment support
5. Integrate with existing CLI and options

### Code Changes

**New File**: `internal/config/config.go`
```go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    "reflect"
    "strings"
    "time"

    "github.com/BurntSushi/toml"
    "github.com/spf13/viper"
    "your-module/internal/engine"
)

// Config represents the complete application configuration
type Config struct {
    // Core application settings
    Core     CoreConfig     `mapstructure:"core" toml:"core" yaml:"core" json:"core"`

    // Generator configuration
    Generator GeneratorConfig `mapstructure:"generator" toml:"generator" yaml:"generator" json:"generator"`

    // Registry configuration
    Registry RegistryConfig `mapstructure:"registry" toml:"registry" yaml:"registry" json:"registry"`

    // Plugin configuration
    Plugins  PluginConfig   `mapstructure:"plugins" toml:"plugins" yaml:"plugins" json:"plugins"`

    // Template configuration
    Templates TemplateConfig `mapstructure:"templates" toml:"templates" yaml:"templates" json:"templates"`

    // Development configuration
    Development DevelopmentConfig `mapstructure:"development" toml:"development" yaml:"development" json:"development"`
}

// CoreConfig contains fundamental application settings
type CoreConfig struct {
    // Verbose enables detailed output
    Verbose bool `mapstructure:"verbose" toml:"verbose" yaml:"verbose" json:"verbose"`

    // Profile specifies the configuration profile (dev, prod, test)
    Profile string `mapstructure:"profile" toml:"profile" yaml:"profile" json:"profile"`

    // LogLevel controls logging verbosity
    LogLevel string `mapstructure:"log_level" toml:"log_level" yaml:"log_level" json:"log_level"`

    // DataDir overrides the default data directory
    DataDir string `mapstructure:"data_dir" toml:"data_dir" yaml:"data_dir" json:"data_dir"`

    // ConfigDir specifies where to look for configuration files
    ConfigDir string `mapstructure:"config_dir" toml:"config_dir" yaml:"config_dir" json:"config_dir"`
}

// GeneratorConfig controls project generation behavior
type GeneratorConfig struct {
    // Engine specifies the default template engine
    Engine string `mapstructure:"engine" toml:"engine" yaml:"engine" json:"engine"`

    // DefaultVariables provides default template variables
    DefaultVariables map[string]string `mapstructure:"default_variables" toml:"default_variables" yaml:"default_variables" json:"default_variables"`

    // PreservePermissions controls file permission handling
    PreservePermissions bool `mapstructure:"preserve_permissions" toml:"preserve_permissions" yaml:"preserve_permissions" json:"preserve_permissions"`

    // ConcurrentProcessing enables parallel file processing
    ConcurrentProcessing bool `mapstructure:"concurrent_processing" toml:"concurrent_processing" yaml:"concurrent_processing" json:"concurrent_processing"`

    // MaxWorkers limits concurrent operations
    MaxWorkers int `mapstructure:"max_workers" toml:"max_workers" yaml:"max_workers" json:"max_workers"`

    // BufferSize controls I/O buffer sizes
    BufferSize int `mapstructure:"buffer_size" toml:"buffer_size" yaml:"buffer_size" json:"buffer_size"`

    // ProcessingTimeout sets maximum time for generation
    ProcessingTimeout time.Duration `mapstructure:"processing_timeout" toml:"processing_timeout" yaml:"processing_timeout" json:"processing_timeout"`
}

// RegistryConfig controls template registry behavior
type RegistryConfig struct {
    // AutoUpdate enables automatic template updates
    AutoUpdate bool `mapstructure:"auto_update" toml:"auto_update" yaml:"auto_update" json:"auto_update"`

    // UpdateInterval specifies how often to check for updates
    UpdateInterval time.Duration `mapstructure:"update_interval" toml:"update_interval" yaml:"update_interval" json:"update_interval"`

    // RemoteRegistries lists additional template sources
    RemoteRegistries []RemoteRegistry `mapstructure:"remote_registries" toml:"remote_registries" yaml:"remote_registries" json:"remote_registries"`

    // CacheTimeout controls how long to cache remote templates
    CacheTimeout time.Duration `mapstructure:"cache_timeout" toml:"cache_timeout" yaml:"cache_timeout" json:"cache_timeout"`
}

// RemoteRegistry represents a remote template registry
type RemoteRegistry struct {
    Name     string `mapstructure:"name" toml:"name" yaml:"name" json:"name"`
    URL      string `mapstructure:"url" toml:"url" yaml:"url" json:"url"`
    Username string `mapstructure:"username" toml:"username" yaml:"username" json:"username"`
    Token    string `mapstructure:"token" toml:"token" yaml:"token" json:"token"`
    Enabled  bool   `mapstructure:"enabled" toml:"enabled" yaml:"enabled" json:"enabled"`
}

// PluginConfig controls plugin system behavior
type PluginConfig struct {
    // Enabled controls whether plugins are loaded
    Enabled bool `mapstructure:"enabled" toml:"enabled" yaml:"enabled" json:"enabled"`

    // Directories lists where to search for plugins
    Directories []string `mapstructure:"directories" toml:"directories" yaml:"directories" json:"directories"`

    // TrustedAuthors lists trusted plugin authors
    TrustedAuthors []string `mapstructure:"trusted_authors" toml:"trusted_authors" yaml:"trusted_authors" json:"trusted_authors"`

    // MaxPlugins limits the number of loaded plugins
    MaxPlugins int `mapstructure:"max_plugins" toml:"max_plugins" yaml:"max_plugins" json:"max_plugins"`

    // PluginTimeout sets maximum time for plugin operations
    PluginTimeout time.Duration `mapstructure:"plugin_timeout" toml:"plugin_timeout" yaml:"plugin_timeout" json:"plugin_timeout"`

    // SandboxEnabled enables plugin sandboxing
    SandboxEnabled bool `mapstructure:"sandbox_enabled" toml:"sandbox_enabled" yaml:"sandbox_enabled" json:"sandbox_enabled"`
}

// TemplateConfig controls template processing behavior
type TemplateConfig struct {
    // Syntax specifies enabled template syntaxes
    Syntax []string `mapstructure:"syntax" toml:"syntax" yaml:"syntax" json:"syntax"`

    // Extensions maps file extensions to template engines
    Extensions map[string]string `mapstructure:"extensions" toml:"extensions" yaml:"extensions" json:"extensions"`

    // IgnorePatterns lists patterns to ignore during generation
    IgnorePatterns []string `mapstructure:"ignore_patterns" toml:"ignore_patterns" yaml:"ignore_patterns" json:"ignore_patterns"`

    // BinaryExtensions lists extensions to treat as binary
    BinaryExtensions []string `mapstructure:"binary_extensions" toml:"binary_extensions" yaml:"binary_extensions" json:"binary_extensions"`

    // MaxFileSize limits template file sizes
    MaxFileSize int64 `mapstructure:"max_file_size" toml:"max_file_size" yaml:"max_file_size" json:"max_file_size"`
}

// DevelopmentConfig contains development-specific settings
type DevelopmentConfig struct {
    // HotReload enables automatic reloading during development
    HotReload bool `mapstructure:"hot_reload" toml:"hot_reload" yaml:"hot_reload" json:"hot_reload"`

    // WatchPatterns specifies which files to watch for changes
    WatchPatterns []string `mapstructure:"watch_patterns" toml:"watch_patterns" yaml:"watch_patterns" json:"watch_patterns"`

    // EnableDebug enables debug logging and features
    EnableDebug bool `mapstructure:"enable_debug" toml:"enable_debug" yaml:"enable_debug" json:"enable_debug"`

    // ProfilerEnabled enables performance profiling
    ProfilerEnabled bool `mapstructure:"profiler_enabled" toml:"profiler_enabled" yaml:"profiler_enabled" json:"profiler_enabled"`
}

// Manager handles configuration loading and management
type Manager struct {
    config    *Config
    viper     *viper.Viper
    sources   []Source
    validator *Validator
}

// Source represents a configuration source
type Source interface {
    Name() string
    Load(v *viper.Viper) error
    Priority() int
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
    v := viper.New()

    // Set defaults
    setDefaults(v)

    return &Manager{
        config:    &Config{},
        viper:     v,
        sources:   []Source{},
        validator: NewValidator(),
    }
}

// setDefaults establishes default configuration values
func setDefaults(v *viper.Viper) {
    // Core defaults
    v.SetDefault("core.verbose", false)
    v.SetDefault("core.profile", "default")
    v.SetDefault("core.log_level", "info")

    // Generator defaults
    v.SetDefault("generator.engine", "pongo2")
    v.SetDefault("generator.preserve_permissions", true)
    v.SetDefault("generator.concurrent_processing", true)
    v.SetDefault("generator.max_workers", 4)
    v.SetDefault("generator.buffer_size", 65536) // 64KB
    v.SetDefault("generator.processing_timeout", "5m")

    // Registry defaults
    v.SetDefault("registry.auto_update", false)
    v.SetDefault("registry.update_interval", "24h")
    v.SetDefault("registry.cache_timeout", "1h")

    // Plugin defaults
    v.SetDefault("plugins.enabled", true)
    v.SetDefault("plugins.max_plugins", 10)
    v.SetDefault("plugins.plugin_timeout", "30s")
    v.SetDefault("plugins.sandbox_enabled", true)

    // Template defaults
    v.SetDefault("templates.syntax", []string{"pongo2"})
    v.SetDefault("templates.max_file_size", 104857600) // 100MB

    // Development defaults
    v.SetDefault("development.hot_reload", false)
    v.SetDefault("development.enable_debug", false)
    v.SetDefault("development.profiler_enabled", false)
}

// AddSource adds a configuration source
func (m *Manager) AddSource(source Source) {
    m.sources = append(m.sources, source)
}

// Load loads configuration from all sources
func (m *Manager) Load() error {
    // Sort sources by priority
    sources := make([]Source, len(m.sources))
    copy(sources, m.sources)

    // Load from each source in priority order
    for _, source := range sources {
        if err := source.Load(m.viper); err != nil {
            return fmt.Errorf("failed to load from source %s: %w", source.Name(), err)
        }
    }

    // Unmarshal into config struct
    if err := m.viper.Unmarshal(m.config); err != nil {
        return fmt.Errorf("failed to unmarshal configuration: %w", err)
    }

    // Validate configuration
    if err := m.validator.Validate(m.config); err != nil {
        return fmt.Errorf("configuration validation failed: %w", err)
    }

    return nil
}

// Get returns the current configuration
func (m *Manager) Get() *Config {
    return m.config
}

// Set updates a configuration value
func (m *Manager) Set(key string, value interface{}) {
    m.viper.Set(key, value)
}

// GetString retrieves a string configuration value
func (m *Manager) GetString(key string) string {
    return m.viper.GetString(key)
}

// GetBool retrieves a boolean configuration value
func (m *Manager) GetBool(key string) bool {
    return m.viper.GetBool(key)
}

// GetInt retrieves an integer configuration value
func (m *Manager) GetInt(key string) int {
    return m.viper.GetInt(key)
}

// GetDuration retrieves a duration configuration value
func (m *Manager) GetDuration(key string) time.Duration {
    return m.viper.GetDuration(key)
}

// WriteConfig saves current configuration to file
func (m *Manager) WriteConfig(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := toml.NewEncoder(file)
    return encoder.Encode(m.config)
}

// Watch starts watching for configuration changes
func (m *Manager) Watch(callback func(*Config)) error {
    m.viper.WatchConfig()
    m.viper.OnConfigChange(func(e fsnotify.Event) {
        if err := m.viper.Unmarshal(m.config); err == nil {
            callback(m.config)
        }
    })
    return nil
}
```

**New File**: `internal/config/sources.go`
```go
package config

import (
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/viper"
)

// FileSource loads configuration from files
type FileSource struct {
    paths    []string
    name     string
    priority int
}

// NewFileSource creates a file-based configuration source
func NewFileSource(name string, paths []string, priority int) *FileSource {
    return &FileSource{
        paths:    paths,
        name:     name,
        priority: priority,
    }
}

// Name returns the source name
func (fs *FileSource) Name() string {
    return fs.name
}

// Priority returns the source priority
func (fs *FileSource) Priority() int {
    return fs.priority
}

// Load loads configuration from files
func (fs *FileSource) Load(v *viper.Viper) error {
    for _, path := range fs.paths {
        if _, err := os.Stat(path); os.IsNotExist(err) {
            continue
        }

        // Determine config type from extension
        ext := filepath.Ext(path)
        switch ext {
        case ".toml":
            v.SetConfigType("toml")
        case ".yaml", ".yml":
            v.SetConfigType("yaml")
        case ".json":
            v.SetConfigType("json")
        default:
            continue
        }

        v.SetConfigFile(path)
        if err := v.MergeInConfig(); err != nil {
            return err
        }
    }

    return nil
}

// EnvSource loads configuration from environment variables
type EnvSource struct {
    prefix   string
    name     string
    priority int
}

// NewEnvSource creates an environment variable configuration source
func NewEnvSource(name, prefix string, priority int) *EnvSource {
    return &EnvSource{
        prefix:   prefix,
        name:     name,
        priority: priority,
    }
}

// Name returns the source name
func (es *EnvSource) Name() string {
    return es.name
}

// Priority returns the source priority
func (es *EnvSource) Priority() int {
    return es.priority
}

// Load loads configuration from environment variables
func (es *EnvSource) Load(v *viper.Viper) error {
    v.SetEnvPrefix(es.prefix)
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    v.AutomaticEnv()

    // Manually bind known environment variables
    envVars := []string{
        "VERBOSE",
        "PROFILE",
        "LOG_LEVEL",
        "DATA_DIR",
        "CONFIG_DIR",
        "GENERATOR_ENGINE",
        "GENERATOR_CONCURRENT_PROCESSING",
        "GENERATOR_MAX_WORKERS",
        "PLUGINS_ENABLED",
        "PLUGINS_SANDBOX_ENABLED",
    }

    for _, envVar := range envVars {
        fullKey := es.prefix + "_" + envVar
        configKey := strings.ToLower(strings.ReplaceAll(envVar, "_", "."))

        if value := os.Getenv(fullKey); value != "" {
            v.Set(configKey, value)
        }
    }

    return nil
}

// CLISource loads configuration from CLI flags
type CLISource struct {
    flags    map[string]interface{}
    name     string
    priority int
}

// NewCLISource creates a CLI flag configuration source
func NewCLISource(name string, flags map[string]interface{}, priority int) *CLISource {
    return &CLISource{
        flags:    flags,
        name:     name,
        priority: priority,
    }
}

// Name returns the source name
func (cs *CLISource) Name() string {
    return cs.name
}

// Priority returns the source priority
func (cs *CLISource) Priority() int {
    return cs.priority
}

// Load loads configuration from CLI flags
func (cs *CLISource) Load(v *viper.Viper) error {
    for key, value := range cs.flags {
        if value != nil {
            v.Set(key, value)
        }
    }
    return nil
}
```

**New File**: `internal/config/validator.go`
```go
package config

import (
    "fmt"
    "path/filepath"
    "strings"
    "time"
)

// Validator validates configuration values
type Validator struct{}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
    return &Validator{}
}

// Validate performs comprehensive configuration validation
func (v *Validator) Validate(config *Config) error {
    if err := v.validateCore(&config.Core); err != nil {
        return fmt.Errorf("core configuration invalid: %w", err)
    }

    if err := v.validateGenerator(&config.Generator); err != nil {
        return fmt.Errorf("generator configuration invalid: %w", err)
    }

    if err := v.validateRegistry(&config.Registry); err != nil {
        return fmt.Errorf("registry configuration invalid: %w", err)
    }

    if err := v.validatePlugins(&config.Plugins); err != nil {
        return fmt.Errorf("plugins configuration invalid: %w", err)
    }

    if err := v.validateTemplates(&config.Templates); err != nil {
        return fmt.Errorf("templates configuration invalid: %w", err)
    }

    return nil
}

// validateCore validates core configuration
func (v *Validator) validateCore(core *CoreConfig) error {
    // Validate profile
    validProfiles := []string{"default", "dev", "development", "prod", "production", "test", "testing"}
    if !v.isValidChoice(core.Profile, validProfiles) {
        return fmt.Errorf("invalid profile %s, must be one of: %s", core.Profile, strings.Join(validProfiles, ", "))
    }

    // Validate log level
    validLogLevels := []string{"debug", "info", "warn", "warning", "error", "fatal", "panic"}
    if !v.isValidChoice(core.LogLevel, validLogLevels) {
        return fmt.Errorf("invalid log level %s, must be one of: %s", core.LogLevel, strings.Join(validLogLevels, ", "))
    }

    // Validate data directory
    if core.DataDir != "" {
        if !filepath.IsAbs(core.DataDir) {
            return fmt.Errorf("data directory must be an absolute path: %s", core.DataDir)
        }
    }

    return nil
}

// validateGenerator validates generator configuration
func (v *Validator) validateGenerator(gen *GeneratorConfig) error {
    // Validate engine
    validEngines := []string{"pongo2", "go-template", "mustache"}
    if !v.isValidChoice(gen.Engine, validEngines) {
        return fmt.Errorf("invalid engine %s, must be one of: %s", gen.Engine, strings.Join(validEngines, ", "))
    }

    // Validate max workers
    if gen.MaxWorkers < 1 {
        return fmt.Errorf("max workers must be at least 1, got %d", gen.MaxWorkers)
    }
    if gen.MaxWorkers > 100 {
        return fmt.Errorf("max workers cannot exceed 100, got %d", gen.MaxWorkers)
    }

    // Validate buffer size
    if gen.BufferSize < 1024 {
        return fmt.Errorf("buffer size must be at least 1024 bytes, got %d", gen.BufferSize)
    }
    if gen.BufferSize > 10*1024*1024 {
        return fmt.Errorf("buffer size cannot exceed 10MB, got %d", gen.BufferSize)
    }

    // Validate processing timeout
    if gen.ProcessingTimeout < time.Second {
        return fmt.Errorf("processing timeout must be at least 1 second, got %v", gen.ProcessingTimeout)
    }
    if gen.ProcessingTimeout > time.Hour {
        return fmt.Errorf("processing timeout cannot exceed 1 hour, got %v", gen.ProcessingTimeout)
    }

    return nil
}

// validateRegistry validates registry configuration
func (v *Validator) validateRegistry(reg *RegistryConfig) error {
    // Validate update interval
    if reg.AutoUpdate && reg.UpdateInterval < time.Minute {
        return fmt.Errorf("update interval must be at least 1 minute when auto-update is enabled, got %v", reg.UpdateInterval)
    }

    // Validate cache timeout
    if reg.CacheTimeout < time.Minute {
        return fmt.Errorf("cache timeout must be at least 1 minute, got %v", reg.CacheTimeout)
    }

    // Validate remote registries
    for i, remote := range reg.RemoteRegistries {
        if remote.Name == "" {
            return fmt.Errorf("remote registry %d: name cannot be empty", i)
        }
        if remote.URL == "" {
            return fmt.Errorf("remote registry %s: URL cannot be empty", remote.Name)
        }
        if !strings.HasPrefix(remote.URL, "http://") && !strings.HasPrefix(remote.URL, "https://") {
            return fmt.Errorf("remote registry %s: URL must start with http:// or https://", remote.Name)
        }
    }

    return nil
}

// validatePlugins validates plugin configuration
func (v *Validator) validatePlugins(plugins *PluginConfig) error {
    // Validate max plugins
    if plugins.MaxPlugins < 1 {
        return fmt.Errorf("max plugins must be at least 1, got %d", plugins.MaxPlugins)
    }
    if plugins.MaxPlugins > 100 {
        return fmt.Errorf("max plugins cannot exceed 100, got %d", plugins.MaxPlugins)
    }

    // Validate plugin timeout
    if plugins.PluginTimeout < time.Second {
        return fmt.Errorf("plugin timeout must be at least 1 second, got %v", plugins.PluginTimeout)
    }

    // Validate plugin directories
    for _, dir := range plugins.Directories {
        if !filepath.IsAbs(dir) {
            return fmt.Errorf("plugin directory must be absolute path: %s", dir)
        }
    }

    return nil
}

// validateTemplates validates template configuration
func (v *Validator) validateTemplates(templates *TemplateConfig) error {
    // Validate syntax types
    validSyntax := []string{"pongo2", "go-template", "mustache"}
    for _, syntax := range templates.Syntax {
        if !v.isValidChoice(syntax, validSyntax) {
            return fmt.Errorf("invalid template syntax %s, must be one of: %s", syntax, strings.Join(validSyntax, ", "))
        }
    }

    // Validate max file size
    if templates.MaxFileSize < 1024 {
        return fmt.Errorf("max file size must be at least 1024 bytes, got %d", templates.MaxFileSize)
    }
    if templates.MaxFileSize > 1024*1024*1024 {
        return fmt.Errorf("max file size cannot exceed 1GB, got %d", templates.MaxFileSize)
    }

    return nil
}

// isValidChoice checks if value is in the list of valid choices
func (v *Validator) isValidChoice(value string, choices []string) bool {
    for _, choice := range choices {
        if strings.EqualFold(value, choice) {
            return true
        }
    }
    return false
}
```

## Files to Create/Modify
- `internal/config/config.go` (new file)
- `internal/config/sources.go` (new file)
- `internal/config/validator.go` (new file)
- `internal/config/config_test.go` (new file)
- `cmd/config.go` (new CLI commands for config management)
- `examples/config/ason.toml` (example configuration file)
- Update existing packages to use centralized config

## Integration Points
- Update `cmd/` packages to use config manager
- Modify `internal/generator/` to read from config
- Update `internal/registry/` to use config settings
- Integrate with plugin system configuration

## Configuration Profiles
- **default**: Basic settings for general use
- **dev/development**: Development-friendly settings with hot-reload
- **prod/production**: Production settings with security focus
- **test/testing**: Testing environment with mocked services

## Definition of Done
- Centralized configuration system handles all app settings
- Multiple configuration sources work with proper precedence
- Configuration validation prevents invalid settings
- Profile-based configuration supports different environments
- Hot-reloading works during development
- CLI commands allow configuration management
- Documentation covers all configuration options
- Backward compatibility maintained with existing behavior