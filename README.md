# ※ Ason - Shake Your Projects Into Being

> *The sacred rattle that transforms templates into living code*

A powerful, lightweight project scaffolding tool that catalyzes the transformation of templates into fully-formed projects. Built with Go's simplicity and minimal dependencies.

## Installation

### macOS (Homebrew - Recommended)

```bash
# Add the tap and install
brew tap madstone-tech/tap
brew install ason

# Verify installation
ason --version
```

Shell completions are automatically installed with Homebrew.

### Linux

#### Download Binary (All distributions)
```bash
# AMD64
curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Linux_x86_64.tar.gz | tar xz
sudo mv ason /usr/local/bin/

# ARM64
curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Linux_arm64.tar.gz | tar xz
sudo mv ason /usr/local/bin/
```

### Windows

Download the latest `.zip` from [releases](https://github.com/madstone-tech/ason/releases/latest), extract, and add to your PATH.

### Using Go

```bash
go install github.com/madstone-tech/ason@latest
```

## Quick Start

```bash
# View available commands
ason --help

# Register a template
ason register my-template ./path/to/template

# List registered templates
ason list

# Create a project from a template
ason new my-template my-project

# Or use a local template directly
ason new ./path/to/template my-project
```

## Features

- 🪇 **Rhythmic Generation**: Fast, lightweight operation with minimal dependencies
- 🎭 **Jinja2-like Templating**: Uses Pongo2 for familiar template syntax
- 📿 **Template Registry**: Local management of your template collection
- 💫 **Interactive Prompts**: Beautiful terminal UI with Bubble Tea
- 🔮 **Terraform Integration**: Seamless infrastructure-as-code support
- ⚡ **Shell Autocompletion**: Tab completion for bash, zsh, and fish

## Shell Autocompletion

Ason supports autocompletion for bash, zsh, and fish shells, providing intelligent suggestions for:

- Template names from your registry
- File and directory paths
- Command flags and options
- Variable names (name=, version=, etc.)

### Quick Installation

```bash
# Install completion for your current shell
./scripts/install-completion.sh

# Or install for specific shells
./scripts/install-completion.sh bash zsh fish

# Or install for all available shells
./scripts/install-completion.sh all
```

### Manual Installation

#### Bash
```bash
# Generate completion script
ason completion bash > ~/.local/share/bash-completion/completions/ason

# Or for system-wide installation (requires sudo)
sudo ason completion bash > /usr/share/bash-completion/completions/ason
```

#### Zsh

**For Oh My Zsh users (recommended):**
```bash
# Create custom plugin
mkdir -p ~/.oh-my-zsh/custom/plugins/ason
ason completion zsh > ~/.oh-my-zsh/custom/plugins/ason/_ason

# Add to plugins list in ~/.zshrc
# plugins=(... ason)
```

**Alternative for Oh My Zsh:**
```bash
# Direct to completions directory (if it exists)
ason completion zsh > ~/.oh-my-zsh/completions/_ason
```

**For standard Zsh setup:**
```bash
mkdir -p ~/.zfunc
ason completion zsh > ~/.zfunc/_ason

# Add to ~/.zshrc if not already present
echo 'fpath=(~/.zfunc $fpath)' >> ~/.zshrc
echo 'autoload -U compinit' >> ~/.zshrc
echo 'compinit' >> ~/.zshrc
```

#### Fish
```bash
# Install fish completion
ason completion fish > ~/.config/fish/completions/ason.fish
```

### Testing Completion

After installation, test completion by typing:
```bash
ason <TAB>          # Shows available commands
ason new <TAB>      # Shows available templates and directories
ason remove <TAB>   # Shows templates in registry
ason --<TAB>        # Shows available flags
```

## Building and Releasing

This project uses [Taskfile](https://taskfile.dev) for automation and [GoReleaser](https://goreleaser.com) for building and releasing.

### Development Tasks

```bash
# Install Task (if not already installed)
go install github.com/go-task/task/v3/cmd/task@latest

# Show all available tasks
task

# Setup and installation
task setup              # Complete development environment setup
task install            # Install binary to GOPATH/bin
task uninstall          # Remove binary from GOPATH/bin

# Common development tasks
task build              # Build for current platform
task test               # Run tests
task test:coverage      # Run tests with coverage
task lint               # Run linters
task clean              # Clean build artifacts
```

### Building for Multiple Platforms

```bash
# Build for all platforms (local, no publish)
task release:local

# Create snapshot release (test release process)
task release:snapshot

# Test release process (dry run)
task release:test
```

### Creating Releases

```bash
# Create and push a git tag
task git:tag TAG=v1.0.0

# Publish release (requires GITHUB_TOKEN)
export GITHUB_TOKEN=your_token_here
task release:publish
```

### Release Artifacts

The release process automatically creates:
- **Homebrew formula** (macOS)
- **Universal binaries** for macOS (Intel + Apple Silicon)
- **Binary archives** for Linux (amd64, arm64) and Windows
- **Source code** archives

## About the Name

Ason is named after the sacred rattle used in Haitian Vodou ceremonies - a tool that catalyzes transformation and invokes change. Just as the ason's rhythm activates spiritual work, this tool activates the transformation of templates into living projects.

## License

MIT
