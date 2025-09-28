# ※ Installation Guide

> *Prepare your system to receive the sacred rattle's power*

This guide walks you through installing Ason on your system using various methods, from pre-built binaries to building from source.

## Quick Installation

### Pre-built Binaries (Recommended)

The fastest way to get Ason running is using pre-built binaries:

#### Linux
```bash
# x86_64
curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Linux_x86_64.tar.gz | tar xz
sudo mv ason /usr/local/bin/

# ARM64
curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Linux_arm64.tar.gz | tar xz
sudo mv ason /usr/local/bin/

# Verify installation
ason --version
```

#### macOS
```bash
# Intel Macs
curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Darwin_x86_64.tar.gz | tar xz
sudo mv ason /usr/local/bin/

# Apple Silicon (M1/M2)
curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Darwin_arm64.tar.gz | tar xz
sudo mv ason /usr/local/bin/

# Verify installation
ason --version
```

#### Windows
```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/madstone-tech/ason/releases/latest/download/ason_Windows_x86_64.zip" -OutFile "ason.zip"
Expand-Archive -Path "ason.zip" -DestinationPath "."

# Move to a directory in PATH
Move-Item "ason.exe" "C:\Program Files\ason\"

# Add to PATH or verify installation
ason --version
```

### Package Managers

#### Homebrew (macOS/Linux)
```bash
# Add the tap
brew tap madstone-tech/tap

# Install Ason
brew install ason

# Verify installation
ason --version
```

#### Scoop (Windows)
```powershell
# Add the bucket (coming soon)
scoop bucket add madstone-tech https://github.com/madstone-tech/scoop-bucket

# Install Ason
scoop install ason

# Verify installation
ason --version
```

#### APT (Debian/Ubuntu)
```bash
# Download and install DEB package
curl -sLO https://github.com/madstone-tech/ason/releases/latest/download/ason_amd64.deb
sudo dpkg -i ason_amd64.deb

# Or install dependencies if needed
sudo apt-get install -f

# Verify installation
ason --version
```

#### YUM/DNF (RedHat/CentOS/Fedora)
```bash
# Download and install RPM package
curl -sLO https://github.com/madstone-tech/ason/releases/latest/download/ason_x86_64.rpm
sudo rpm -i ason_x86_64.rpm

# Or with dnf
sudo dnf install ason_x86_64.rpm

# Verify installation
ason --version
```

#### AUR (Arch Linux)
```bash
# Using yay
yay -S ason

# Using paru
paru -S ason

# Manual installation
git clone https://aur.archlinux.org/ason.git
cd ason
makepkg -si
```

### Go Install

If you have Go installed:

```bash
# Install latest version
go install github.com/madstone-tech/ason@latest

# Install specific version
go install github.com/madstone-tech/ason@v1.0.0

# Verify installation (ensure $GOPATH/bin is in PATH)
ason --version
```

## Building from Source

### Prerequisites

- **Go 1.25.1+** (required)
- **Git** (for cloning)
- **Task** (optional, for development tasks)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/madstone-tech/ason.git
cd ason

# Install Task (optional but recommended)
go install github.com/go-task/task/v3/cmd/task@latest

# Complete setup (builds, installs, and sets up completion)
task setup

# Or manual build
go mod tidy
go build -o ason .

# Install to GOPATH/bin
task install
```

### Build Options

#### Quick Build
```bash
# Build for current platform
task build

# Or with go directly
go build -o ason .
```

#### Multi-platform Build
```bash
# Build for all platforms (requires GoReleaser)
task release:local

# Install GoReleaser if not present
go install github.com/goreleaser/goreleaser@latest
```

#### Development Build with Debugging
```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o ason-debug .

# Build with race detection
go build -race -o ason-race .
```

#### Custom Build with Version Information
```bash
# Build with version information
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" -o ason .
```

## Docker Installation

### Using Pre-built Image

```bash
# Pull the image
docker pull ghcr.io/madstone-tech/ason:latest

# Run Ason in container
docker run --rm -v $(pwd):/workspace ghcr.io/madstone-tech/ason:latest --version

# Create alias for easier use
echo 'alias ason="docker run --rm -v \$(pwd):/workspace ghcr.io/madstone-tech/ason:latest"' >> ~/.bashrc
source ~/.bashrc
```

### Building Docker Image

```bash
# Clone repository
git clone https://github.com/madstone-tech/ason.git
cd ason

# Build Docker image
task docker:build

# Or build manually
docker build -f Dockerfile.goreleaser -t ason:local .

# Run containerized Ason
task docker:run -- --version
```

## Verification

### Check Installation

```bash
# Verify Ason is installed and accessible
ason --version

# Should output something like:
# ※ Ason v1.0.0

# Check available commands
ason --help

# Test basic functionality
ason list
```

### System Information

```bash
# Check where Ason is installed
which ason

# Check version details
ason --version

# Verify Go version (if building from source)
go version
```

## Shell Completion Setup

After installation, set up shell completion for the best experience:

### Automatic Installation
```bash
# Install completion for current shell
./scripts/install-completion.sh

# Or if using task
task completion:install
```

### Manual Installation

#### Bash
```bash
# Generate and install completion
ason completion bash > ~/.local/share/bash-completion/completions/ason

# Add to ~/.bashrc if needed
echo '. ~/.local/share/bash-completion/completions/ason' >> ~/.bashrc
```

#### Zsh (Oh My Zsh)
```bash
# Create plugin directory
mkdir -p ~/.oh-my-zsh/custom/plugins/ason

# Generate completion
ason completion zsh > ~/.oh-my-zsh/custom/plugins/ason/_ason

# Add to plugins in ~/.zshrc
# plugins=(... ason)
```

#### Fish
```bash
# Generate completion
ason completion fish > ~/.config/fish/completions/ason.fish
```

## Configuration

### Initial Configuration

Ason works out of the box, but you can customize its behavior:

```bash
# Create configuration directory
mkdir -p ~/.ason

# Example configuration file (~/.ason/config.yaml)
cat > ~/.ason/config.yaml << EOF
# Ason Configuration
default_template_dir: ~/templates
registry_path: ~/.ason/templates
auto_backup: true
completion_cache: true

# Default variables
defaults:
  author: "Your Name"
  email: "your.email@example.com"
  license: "MIT"
EOF
```

### Environment Variables

```bash
# Set default template directory
export ASON_TEMPLATE_DIR=~/my-templates

# Set custom registry path
export ASON_REGISTRY_PATH=~/.my-ason-templates

# Enable debug mode
export ASON_DEBUG=true

# Add to shell profile
echo 'export ASON_TEMPLATE_DIR=~/my-templates' >> ~/.bashrc
```

## System Requirements

### Minimum Requirements
- **OS**: Linux, macOS, Windows
- **Architecture**: x86_64, ARM64
- **Memory**: 64MB RAM
- **Disk**: 10MB free space

### Recommended Requirements
- **Memory**: 128MB+ RAM for large templates
- **Disk**: 100MB+ for template registry
- **Network**: Internet access for downloading templates

### Supported Platforms

| Platform | Architecture | Status |
|----------|-------------|---------|
| Linux | x86_64 | ✅ Full support |
| Linux | ARM64 | ✅ Full support |
| Linux | ARM | ✅ Community support |
| macOS | x86_64 (Intel) | ✅ Full support |
| macOS | ARM64 (M1/M2) | ✅ Full support |
| Windows | x86_64 | ✅ Full support |
| Windows | ARM64 | 🚧 Experimental |
| FreeBSD | x86_64 | 🚧 Experimental |

## Troubleshooting Installation

### Common Issues

#### Permission Denied
```bash
# If you get permission denied when installing to /usr/local/bin
sudo mv ason /usr/local/bin/
sudo chmod +x /usr/local/bin/ason

# Or install to user directory
mkdir -p ~/bin
mv ason ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
```

#### Command Not Found
```bash
# Check if installation directory is in PATH
echo $PATH

# Add to PATH if needed
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Or find where ason was installed
find / -name "ason" 2>/dev/null
```

#### Go Module Issues (when building from source)
```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy

# Build again
go build -o ason .
```

#### macOS Gatekeeper Issues
```bash
# If macOS blocks the binary
sudo xattr -rd com.apple.quarantine /usr/local/bin/ason

# Or allow in System Preferences > Security & Privacy
```

#### Windows Antivirus Issues
```bash
# If Windows Defender blocks the executable
# Add exception in Windows Security settings
# Or download from official releases page
```

### Verification Steps

```bash
# Check binary exists and is executable
ls -la $(which ason)

# Test basic functionality
ason --help
ason --version
ason list

# Check file permissions
stat $(which ason)

# Test template operations
mkdir test-template
echo "# {{ project_name }}" > test-template/README.md
ason new test-template test-output --dry-run
rm -rf test-template test-output
```

### Getting Help

If you encounter issues:

1. **Check the documentation**: [Troubleshooting Guide](../troubleshooting/common-issues.md)
2. **Search existing issues**: [GitHub Issues](https://github.com/madstone-tech/ason/issues)
3. **Create a new issue**: Include:
   - OS and architecture
   - Installation method used
   - Error messages
   - Output of `ason --version`

## Updating Ason

### Package Manager Updates

```bash
# Homebrew
brew update && brew upgrade ason

# APT
sudo apt update && sudo apt upgrade ason

# Go install
go install github.com/madstone-tech/ason@latest
```

### Manual Updates

```bash
# Download latest binary
curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_$(uname -s)_$(uname -m).tar.gz | tar xz

# Replace existing binary
sudo mv ason /usr/local/bin/

# Verify update
ason --version
```

### Development Updates

```bash
# Update from source
cd ason
git pull origin main
task clean
task build
task install
```

## Uninstallation

### Removing Ason

```bash
# Remove binary
sudo rm /usr/local/bin/ason

# Or from GOPATH
rm $GOPATH/bin/ason

# Remove configuration and templates (optional)
rm -rf ~/.ason

# Remove completion scripts
rm ~/.local/share/bash-completion/completions/ason
rm ~/.oh-my-zsh/custom/plugins/ason/_ason
rm ~/.config/fish/completions/ason.fish
```

### Package Manager Removal

```bash
# Homebrew
brew uninstall ason
brew untap madstone-tech/tap

# APT
sudo apt remove ason

# YUM/DNF
sudo dnf remove ason
```

## Next Steps

After successful installation:

1. **[Quick Start Guide](quick-start.md)** - Get up and running in 5 minutes
2. **[Create Your First Template](first-template.md)** - Learn template creation
3. **[Configuration Guide](configuration.md)** - Customize Ason for your workflow

---

*The sacred rattle is now prepared for transformation. Let the journey begin! 🪇*