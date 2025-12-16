# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2025-12-16

### Added - Library Export (Feature 002)

#### Core Generation API
- **Generator struct**: Provides the main interface for programmatic template generation
  - `NewGenerator(engine Engine) (*Generator, error)` - Creates a new generator with a template engine
  - `Generate(ctx context.Context, templatePath string, variables map[string]interface{}, outputPath string) error` - Generates projects from templates
  - `GetEngine() Engine` - Retrieves the configured engine
- Context cancellation support for graceful timeout/cancellation handling
- Full input validation (empty paths, invalid variable names)
- Path traversal prevention for security
- Binary file preservation during generation
- Thread-safe implementation with RWMutex protection

#### Registry Management API
- **Registry struct**: Template registry management with XDG compliance
  - `NewRegistry() (*Registry, error)` - Uses default XDG location (~/.local/share/ason/registry.toml)
  - `NewRegistryAt(registryPath string) (*Registry, error)` - Custom registry path
  - `Register(name, templatePath, description string) error` - Register templates
  - `List() ([]TemplateInfo, error)` - List all registered templates (alphabetically sorted)
  - `Remove(name string) error` - Remove templates (idempotent)
- **TemplateInfo struct**: Template metadata (Name, Path, Created, Description)
- Atomic TOML persistence (temp file + rename pattern)
- Thread-safe concurrent operations (multiple concurrent reads, serialized writes)
- Proper XDG Base Directory compliance

#### Engine Interface
- **Engine interface**: Pluggable template engine support
  - `Render(template string, context map[string]interface{}) (string, error)` - Render template strings
  - `RenderFile(filePath string, context map[string]interface{}) (string, error)` - Render template files
- **NewDefaultEngine()**: Pongo2 template engine implementation
- **RenderWithEngine()**: Context-aware rendering helper function
- Comprehensive engine interface documentation in `docs/api/engine_interface.md`
- Custom engine implementation guidelines

#### Error Types
- **TemplateNotFoundError**: When template path doesn't exist
- **InvalidPathError**: For path traversal or invalid paths
- **VariableValidationError**: For invalid variable names/values
- **GenerationError**: For template rendering failures
- **EngineError**: For engine-specific failures
- All errors support `Unwrap()` for error chaining

#### API Locations
- Public API: `pkg/generator.go`, `pkg/engine.go`, `pkg/registry.go`
- Documentation: `docs/api/engine_interface.md`
- Examples: Included in inline godoc comments

### Improved
- Test infrastructure with proper temporary directory handling
- Fixed .golangci.yml configuration syntax
- Removed duplicate test files, consolidated test suites

### Changed
- **BREAKING (Minor)**: None - fully backward compatible with CLI

### Security
- Input validation for all paths and variable names
- Path traversal prevention checks
- Registry operations are atomic with temporary file pattern

## [0.2.2] - 2025-10-22

### Fixed
- Fixed version injection timing issue where version was set after Cobra command initialization
- Updated tests to use XDG-compliant paths and TOML format for registry

## [0.2.1] - 2025-10-22

### Fixed
- Version number now correctly displays from GoReleaser ldflags injection (partial fix)

## [0.2.0] - 2025-10-22

### Changed
- Renamed `ason add` command to `ason register` for better semantic clarity (backward compatible via alias)
- Improved README installation instructions with focus on Homebrew for macOS
- Updated all documentation to use `register` as the primary command name

### Added
- Created CHANGELOG.md following Keep a Changelog format
- Added `.specify/`, `.opencode/`, and `specs/` to .gitignore

## [0.1.0] - Previous Release

See git tags for release history.

[0.3.0]: https://github.com/madstone-tech/ason/compare/v0.2.2...v0.3.0
[0.2.2]: https://github.com/madstone-tech/ason/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/madstone-tech/ason/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/madstone-tech/ason/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/madstone-tech/ason/releases/tag/v0.1.0
