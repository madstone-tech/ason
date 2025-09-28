# ※ ason completion

> *Empower your shell with the wisdom of the sacred rattle*

The `ason completion` command generates shell completion scripts that provide intelligent tab completion for Ason commands, making your workflow faster and more intuitive.

## Synopsis

```bash
ason completion [SHELL] [flags]
```

## Description

The `completion` command generates autocompletion scripts for your shell, enabling tab completion for commands, template names, file paths, and variable suggestions. This dramatically improves the user experience by reducing typing and preventing errors.

## Arguments

### SHELL
The shell for which to generate completion script.

**Supported shells:**
- `bash` - Bash shell completion
- `zsh` - Zsh shell completion
- `fish` - Fish shell completion
- `powershell` - PowerShell completion (Windows)

## Flags

### --no-descriptions
Disable completion descriptions for faster performance.

```bash
# Generate without descriptions
ason completion zsh --no-descriptions
```

### Global Flags
- `-h, --help` - Show help for the command
- `-v, --version` - Show Ason version

## Examples

### Generate Completion Scripts

```bash
# Generate bash completion
ason completion bash

# Generate zsh completion
ason completion zsh

# Generate fish completion
ason completion fish

# Generate without descriptions (faster)
ason completion bash --no-descriptions
```

### Install Completion Scripts

```bash
# Bash - system-wide
sudo ason completion bash > /usr/share/bash-completion/completions/ason

# Bash - user-specific
ason completion bash > ~/.local/share/bash-completion/completions/ason

# Zsh - Oh My Zsh
mkdir -p ~/.oh-my-zsh/custom/plugins/ason
ason completion zsh > ~/.oh-my-zsh/custom/plugins/ason/_ason

# Zsh - standard
mkdir -p ~/.zfunc
ason completion zsh > ~/.zfunc/_ason

# Fish
ason completion fish > ~/.config/fish/completions/ason.fish
```

## Installation Guide

### Automatic Installation

Use the provided installation script for easy setup:

```bash
# Install for current shell
./scripts/install-completion.sh

# Install for specific shells
./scripts/install-completion.sh bash zsh fish

# Install for all available shells
./scripts/install-completion.sh all
```

### Manual Installation

#### Bash Completion

**Prerequisites:**
```bash
# Install bash-completion if not already installed
# Ubuntu/Debian
sudo apt install bash-completion

# macOS with Homebrew
brew install bash-completion
```

**Installation:**
```bash
# User-specific installation (recommended)
mkdir -p ~/.local/share/bash-completion/completions
ason completion bash > ~/.local/share/bash-completion/completions/ason

# System-wide installation (requires sudo)
sudo ason completion bash > /usr/share/bash-completion/completions/ason

# Verify installation
. ~/.local/share/bash-completion/completions/ason
```

**Add to ~/.bashrc:**
```bash
# Ensure bash-completion is loaded
if [ -f /usr/share/bash-completion/bash_completion ]; then
  . /usr/share/bash-completion/bash_completion
elif [ -f /etc/bash_completion ]; then
  . /etc/bash_completion
fi

# Load user completions
if [ -d ~/.local/share/bash-completion/completions ]; then
  for completion in ~/.local/share/bash-completion/completions/*; do
    [ -r "$completion" ] && . "$completion"
  done
fi
```

#### Zsh Completion

**For Oh My Zsh users (recommended):**
```bash
# Create custom plugin directory
mkdir -p ~/.oh-my-zsh/custom/plugins/ason

# Generate completion script
ason completion zsh > ~/.oh-my-zsh/custom/plugins/ason/_ason

# Add to plugins list in ~/.zshrc
plugins=(... ason)

# Reload zsh
source ~/.zshrc
```

**Alternative Oh My Zsh method:**
```bash
# If completions directory exists
ason completion zsh > ~/.oh-my-zsh/completions/_ason

# Rebuild completion cache
rm -f ~/.zcompdump; compinit
```

**For standard Zsh setup:**
```bash
# Create function directory
mkdir -p ~/.zfunc

# Generate completion script
ason completion zsh > ~/.zfunc/_ason

# Add to ~/.zshrc (if not already present)
echo 'fpath=(~/.zfunc $fpath)' >> ~/.zshrc
echo 'autoload -U compinit' >> ~/.zshrc
echo 'compinit' >> ~/.zshrc

# Reload configuration
source ~/.zshrc
```

#### Fish Completion

```bash
# Create completions directory if it doesn't exist
mkdir -p ~/.config/fish/completions

# Generate completion script
ason completion fish > ~/.config/fish/completions/ason.fish

# Reload fish
fish -c "source ~/.config/fish/completions/ason.fish"
```

#### PowerShell Completion (Windows)

```powershell
# Generate completion script
ason completion powershell > ason.ps1

# Add to PowerShell profile
Add-Content $PROFILE ". $(Join-Path (Split-Path $PROFILE) 'ason.ps1')"

# Or manually add to profile
echo ". $(Join-Path (Split-Path $PROFILE) 'ason.ps1')" >> $PROFILE
```

## Completion Features

### Command Completion

Tab completion works for all Ason commands:

```bash
ason <TAB>
# Suggests: new, list, add, remove, validate, completion

ason n<TAB>
# Completes to: ason new

ason --<TAB>
# Suggests: --help, --version
```

### Template Name Completion

Intelligent completion for template names from your registry:

```bash
ason new <TAB>
# Shows available templates: react-app, go-service, docs-site

ason remove <TAB>
# Shows removable templates from registry

ason validate <TAB>
# Shows templates and directories for validation
```

### File Path Completion

Smart file and directory path completion:

```bash
ason new my-template <TAB>
# Shows directories for output path

ason add my-template <TAB>
# Shows directories for template source path

ason validate <TAB>
# Shows template files and directories
```

### Variable Name Completion

Completion for common variable patterns:

```bash
ason new my-template output --var <TAB>
# Suggests common patterns: name=, version=, author=, description=

ason new react-app my-app --var n<TAB>
# Completes to: --var name=

ason new go-service api --var name=MyAPI --var <TAB>
# Suggests remaining variables: version=, port=, database=
```

### Flag Completion

Comprehensive flag completion for all commands:

```bash
ason new --<TAB>
# Suggests: --dry-run, --var, --help

ason list --<TAB>
# Suggests: --format, --filter, --sort, --reverse, --help

ason validate --<TAB>
# Suggests: --strict, --format, --fix, --check, --ignore-warnings, --help
```

## Completion Behavior

### Context-Aware Suggestions

The completion system provides context-aware suggestions:

- **After `ason new`**: Shows templates and directories
- **After `ason remove`**: Shows only templates in registry
- **After `ason validate`**: Shows templates, files, and directories
- **After `--var`**: Suggests variable patterns
- **After template name**: Shows appropriate next arguments

### Smart Filtering

Completions are filtered based on what you've already typed:

```bash
ason new re<TAB>
# Shows only templates starting with "re": react-app, react-native

ason list --f<TAB>
# Completes to: ason list --format

ason new react-app output --var na<TAB>
# Completes to: ason new react-app output --var name=
```

### Performance Optimization

- Template names are cached for fast completion
- File system access is minimized
- Completion scripts are optimized for speed
- Use `--no-descriptions` for faster completion in large registries

## Testing Completion

### Verify Installation

```bash
# Test basic command completion
ason <TAB><TAB>

# Test template completion (requires templates in registry)
ason new <TAB><TAB>

# Test flag completion
ason new --<TAB><TAB>

# Test variable completion
ason new template output --var <TAB><TAB>
```

### Troubleshooting Completion

#### Bash Issues
```bash
# Check if bash-completion is loaded
echo $BASH_COMPLETION_VERSINFO

# Manually source completion
source ~/.local/share/bash-completion/completions/ason

# Check completion function exists
complete -p ason
```

#### Zsh Issues
```bash
# Check if completion function is loaded
which _ason

# Rebuild completion cache
rm -f ~/.zcompdump*
compinit

# Check fpath includes completion directory
echo $fpath | grep -o ~/.zfunc
```

#### Fish Issues
```bash
# Check if completion file exists
ls ~/.config/fish/completions/ason.fish

# Test completion manually
complete -C'ason ' ''

# Reload fish configuration
source ~/.config/fish/config.fish
```

## Advanced Configuration

### Custom Completion Functions

You can extend completion with custom functions:

```bash
# ~/.bashrc - Custom template path completion
_ason_custom_templates() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  local templates=$(find ~/my-templates -maxdepth 1 -type d -printf '%f\n')
  COMPREPLY=($(compgen -W "$templates" -- "$cur"))
}

# Override default completion for specific use cases
complete -F _ason_custom_templates ason-dev
```

### Performance Tuning

For large template registries:

```bash
# Generate completion without descriptions
ason completion zsh --no-descriptions > ~/.zfunc/_ason

# Cache template names for faster completion
ason list --format json | jq -r '.templates[].name' > ~/.ason/template-cache

# Use completion cache in custom functions
_cached_templates() {
  if [ -f ~/.ason/template-cache ]; then
    cat ~/.ason/template-cache
  else
    ason list --format json | jq -r '.templates[].name'
  fi
}
```

## Integration Examples

### Development Workflow

```bash
# Quick project creation with completion
ason new <TAB>                    # Select template
ason new react-app <TAB>          # Select output directory
ason new react-app my-app --var <TAB>  # Set variables

# Template management with completion
ason add <TAB>                    # Add new template
ason remove <TAB>                 # Remove old template
ason validate <TAB>               # Validate specific template
```

### Scripting with Completion

```bash
#!/bin/bash
# create-project.sh with completion support

# Source completion for script use
source ~/.local/share/bash-completion/completions/ason

# Interactive template selection
echo "Available templates:"
ason list --format table

read -p "Template name: " template
read -p "Project name: " project

# Use completion data for validation
if ason list --format json | jq -e ".templates[] | select(.name == \"$template\")" > /dev/null; then
  ason new "$template" "$project"
else
  echo "Template not found. Available templates:"
  ason list --format json | jq -r '.templates[].name'
fi
```

### IDE Integration

For editors that support shell completion:

```bash
# VS Code tasks.json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Create Ason Project",
      "type": "shell",
      "command": "ason",
      "args": ["new"],
      "group": "build",
      "presentation": {
        "echo": true,
        "reveal": "always",
        "focus": false,
        "panel": "shared"
      },
      "options": {
        "shell": {
          "executable": "/bin/bash",
          "args": ["-c", "source ~/.bashrc && ${command} ${args}"]
        }
      }
    }
  ]
}
```

## Related Commands

- [`ason new`](new.md) - Create projects (enhanced with completion)
- [`ason list`](list.md) - List templates (used by completion)
- [`ason add`](add.md) - Add templates (enhanced with completion)
- [`ason remove`](remove.md) - Remove templates (enhanced with completion)
- [`ason validate`](validate.md) - Validate templates (enhanced with completion)

## See Also

- [Getting Started Guide](../getting-started/quick-start.md)
- [Template Registry Guide](../guides/registry.md)
- [Troubleshooting Guide](../troubleshooting/common-issues.md)

---

*Let the shell be wise with the knowledge of templates. Tab your way to transformation! 🪇*