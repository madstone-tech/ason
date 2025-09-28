# Contributing to Ason

Thank you for your interest in contributing to Ason! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Process](#contributing-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Security](#security)
- [Community](#community)

## Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow. Please be respectful, inclusive, and constructive in all interactions.

### Our Standards

- Use welcoming and inclusive language
- Be respectful of differing viewpoints and experiences
- Gracefully accept constructive criticism
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- **Go**: Version 1.25 or higher
- **Git**: For version control
- **Task**: For development tasks (recommended, install with `go install github.com/go-task/task/v3/cmd/task@latest`)

### First Time Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:

   ```bash
   git clone https://github.com/YOUR-USERNAME/ason.git
   cd ason
   ```

3. **Add upstream remote**:

   ```bash
   git remote add upstream https://github.com/madstone-tech/ason.git
   ```

4. **Install dependencies**:

   ```bash
   go mod download
   ```

5. **Run the tests** to ensure everything works:

   ```bash
   task test
   # or
   go test ./...
   ```

## Development Setup

### Project Structure

```
ason/
├── cmd/                 # CLI command implementations
│   ├── root.go         # Root command setup
│   ├── new.go          # Project generation command
│   └── ...
├── internal/           # Internal packages
│   ├── engine/         # Template rendering engines
│   ├── generator/      # Project generation logic
│   ├── registry/       # Template registry management
│   └── ...
├── examples/           # Example templates and usage
├── docs/              # Documentation
├── roadmap/           # Implementation roadmap
└── test/              # Integration tests
```

### Development Commands

Using Task (recommended):

```bash
# Setup development environment
task setup

# Build the binary
task build

# Run tests
task test

# Run linting
task lint

# Clean build artifacts
task clean

# Show all available tasks
task
```

Using Make (legacy):

```bash
make build    # Build binary
make test     # Run tests
make clean    # Clean artifacts
```

Using Go directly:

```bash
go build -o ason .           # Build
go test ./...                # Test
go run . new template output # Run locally
```

### Running Ason Locally

```bash
# Build and test locally
go build -o ason .
./ason --help

# Test with a template
./ason new examples/simple-template ./test-output
```

## Contributing Process

### 1. Planning Your Contribution

**For Bug Fixes:**

- Check existing issues to avoid duplicates
- Create an issue if one doesn't exist
- Discuss the approach before implementing

**For New Features:**

- Check the [roadmap](./roadmap/) and existing issues
- Create a feature request issue first
- Wait for maintainer feedback before implementation
- Consider backward compatibility

**For Security Issues:**

- **Critical vulnerabilities**: Email <security@madstone.tech> privately
- **Minor security improvements**: Create a security issue

### 2. Development Workflow

1. **Create a feature branch**:

   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-number-description
   ```

2. **Make your changes** following the coding standards

3. **Write tests** for your changes

4. **Update documentation** if needed

5. **Commit your changes** with clear messages:

   ```bash
   git add .
   git commit -m "feat: add new template validation feature"
   ```

6. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** using the PR template

### 3. Pull Request Guidelines

- **Use the PR template** and fill out all relevant sections
- **Keep PRs focused** - one feature or fix per PR
- **Write clear commit messages** following [Conventional Commits](https://conventionalcommits.org/)
- **Update tests** and ensure all tests pass
- **Update documentation** for user-facing changes
- **Rebase your branch** to keep a clean history

### 4. Review Process

- Maintainers will review your PR
- Address feedback promptly
- Keep the PR updated with the main branch
- Once approved, a maintainer will merge your PR

## Coding Standards

### Go Style Guide

We follow standard Go conventions plus some additional guidelines:

- **gofmt**: All code must be formatted with `gofmt`
- **golint**: Code should pass `golint` checks
- **govet**: Code should pass `go vet` checks
- **Naming**: Use clear, descriptive names for functions, variables, and types
- **Comments**: Public APIs must have godoc comments
- **Error handling**: Always handle errors appropriately

### Code Style

```go
// Good: Clear function name and documentation
// ProcessTemplate renders a template with the given context and writes to output.
func ProcessTemplate(templatePath, outputPath string, context map[string]interface{}) error {
    if templatePath == "" {
        return fmt.Errorf("template path cannot be empty")
    }

    // Implementation...
    return nil
}

// Good: Proper error handling
content, err := os.ReadFile(path)
if err != nil {
    return fmt.Errorf("failed to read template file %s: %w", path, err)
}

// Good: Clear variable names
templateEngine := engine.NewPongo2Engine()
generatorOptions := &generator.Options{
    Verbose: true,
    DryRun:  false,
}
```

### Package Guidelines

- **internal/**: Use for implementation details not exposed to users
- **cmd/**: CLI command implementations only
- **Interfaces**: Keep interfaces small and focused
- **Dependencies**: Minimize external dependencies
- **Testing**: Each package should have comprehensive tests

### Commit Message Format

We use [Conventional Commits](https://conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes
- `security`: Security improvements

**Examples:**

```
feat(generator): add concurrent file processing
fix(registry): resolve template path validation issue
docs: update CLI usage examples
test(engine): add comprehensive template rendering tests
```

## Testing Guidelines

### Test Organization

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test complete workflows
- **Table-driven tests**: Use for testing multiple scenarios
- **Mocking**: Mock external dependencies appropriately

### Test Structure

```go
func TestGenerateProject(t *testing.T) {
    tests := []struct {
        name        string
        templateDir string
        outputDir   string
        context     map[string]interface{}
        wantErr     bool
        wantFiles   []string
    }{
        {
            name:        "simple template",
            templateDir: "testdata/simple",
            outputDir:   "testdata/output",
            context:     map[string]interface{}{"name": "test"},
            wantErr:     false,
            wantFiles:   []string{"README.md", "main.go"},
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation...
        })
    }
}
```

### Test Requirements

- **Coverage**: Aim for >80% test coverage
- **Edge cases**: Test error conditions and edge cases
- **Isolation**: Tests should not depend on each other
- **Cleanup**: Clean up test artifacts
- **Fast**: Tests should run quickly

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...

# Run specific test
go test -run TestGenerateProject ./internal/generator

# Run benchmarks
go test -bench=. ./...
```

## Documentation

### Code Documentation

- **Public APIs**: Must have godoc comments
- **Complex logic**: Comment non-obvious code
- **Examples**: Provide examples in godoc when helpful

```go
// Generator handles template processing and project generation.
//
// A Generator combines a template engine with configuration options to transform
// template directories into project structures. It handles file processing,
// directory creation, and variable substitution.
//
// Example usage:
//   engine := engine.NewPongo2Engine()
//   options := &Options{Verbose: true}
//   gen := NewGenerator(engine, options)
//   err := gen.Generate(templatePath, outputPath, context)
type Generator struct {
    // ...
}
```

### User Documentation

- **README**: Keep the main README up to date
- **CLI help**: Ensure help text is clear and accurate
- **Examples**: Provide working examples
- **Changelog**: Document user-facing changes

### Documentation Updates

When making changes, update:

- Godoc comments for changed APIs
- README if CLI behavior changes
- Examples if usage patterns change
- CHANGELOG.md for user-facing changes

## Security

### Security Best Practices

- **Input validation**: Validate all user inputs
- **Path traversal**: Prevent directory traversal attacks
- **Secrets**: Never commit secrets or sensitive data
- **Dependencies**: Keep dependencies updated
- **Error messages**: Don't leak sensitive information

### Reporting Security Issues

- **Critical vulnerabilities**: Email <security@madstone.io>
- **Minor improvements**: Create a security issue on GitHub
- **Include**: Detailed description and reproduction steps
- **Response time**: We aim to respond within 48 hours

### Security Review

All security-related changes require:

- Additional review by maintainers
- Security testing
- Documentation of security implications

## Community

### Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and community discussion
- **Documentation**: Check existing docs first
- **Code**: Read the source code for implementation details

### Communication

- **Be respectful**: Follow the code of conduct
- **Be clear**: Provide context and details
- **Be patient**: Maintainers volunteer their time
- **Be helpful**: Help others when you can

### Recognition

Contributors are recognized in:

- Release notes for significant contributions
- GitHub contributor graphs
- Community discussions and feedback

## Development Roadmap

Check our [implementation roadmap](./roadmap/) for:

- Planned features and improvements
- Current development priorities
- Ways to contribute to specific areas

### Current Focus Areas

1. **Phase 1**: Critical security fixes
2. **Phase 2**: Code quality improvements
3. **Phase 3**: Performance optimizations
4. **Phase 4**: Architecture enhancements
5. **Phase 5**: Testing improvements

## Questions?

If you have questions about contributing:

1. Check this guide and existing documentation
2. Search existing issues and discussions
3. Create a question issue using the question template
4. Join community discussions

Thank you for contributing to Ason! 🚀

