# ※ Configuration Guide

> *Customize the sacred rattle to harmonize with your workflow*

This guide covers all aspects of configuring Ason to work perfectly with your development environment, from basic settings to advanced customizations.

## Configuration Features

- **XDG Base Directory Specification**: Configuration files use XDG-compliant paths for better system integration
- **TOML Format**: Uses TOML configuration format for clear syntax and reduced formatting errors

## Configuration Overview

Ason can be configured through multiple methods, with the following priority order:

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Configuration files**
4. **Default values** (lowest priority)

## Configuration File

### Default Location

Ason looks for configuration files in these locations (in order):

```bash
# XDG-compliant user configuration
~/.config/ason/config.toml
~/.local/share/ason/config.toml

# Project-specific configuration (if in a project directory)
./ason.toml
./.ason/config.toml
```

### Basic Configuration File

Create your configuration file:

```bash
# Create Ason config directory (XDG-compliant)
mkdir -p ~/.config/ason

# Create basic configuration
cat > ~/.config/ason/config.toml << 'EOF'
# Ason Configuration File

# Registry settings
[registry]
path = "~/.local/share/ason/templates"
auto_backup = true
backup_count = 5

# Default template settings
[defaults]
author = "Your Name"
email = "your.email@example.com"
license = "MIT"
version = "1.0.0"
organization = "Your Organization"

# Template directories
template_dirs = [
  "~/templates",
  "~/projects/templates",
  "/shared/templates"
]

# Output settings
[output]
default_dir = "~/projects"
create_parent_dirs = true
overwrite_confirm = true

# UI preferences
[ui]
colors = true
emoji = true
progress_bars = true
confirm_destructive = true

# Completion settings
[completion]
cache_enabled = true
cache_ttl = 3600  # 1 hour in seconds
template_suggestions = true

# Logging
[logging]
level = "info"  # debug, info, warn, error
file = "~/.local/share/ason/ason.log"
max_size = "10MB"
max_backups = 3
EOF
```

### Advanced Configuration

> **Note**: The following example shows advanced configuration options. Convert YAML syntax to TOML format as needed.

```bash
cat > ~/.config/ason/config.yaml << 'EOF'
# Advanced Ason Configuration

# Global settings
global:
  debug: false
  verbose: false
  quiet: false

# Registry configuration
registry:
  path: ~/.ason/templates
  auto_backup: true
  backup_dir: ~/.ason/backups
  backup_count: 10
  compression: gzip
  index_cache: true

# Template discovery
discovery:
  auto_scan: true
  scan_paths:
    - ~/templates
    - ~/projects/templates
    - ~/.local/share/ason/templates
  scan_depth: 3
  ignore_patterns:
    - ".*"
    - "node_modules"
    - "__pycache__"
    - "*.tmp"

# Default variables
defaults:
  # Personal info
  author: "Your Name"
  email: "your.email@example.com"
  github_username: "yourusername"
  organization: "Your Organization"

  # Project defaults
  license: "MIT"
  version: "1.0.0"
  go_version: "1.25"
  node_version: "18"
  python_version: "3.11"

  # Paths and URLs
  homepage: "https://github.com/yourusername"
  repository_base: "github.com/yourusername"

# Template processing
processing:
  # Variable resolution
  variable_prefix: "{{"
  variable_suffix: "}}"
  strict_variables: false

  # File processing
  binary_extensions:
    - ".png"
    - ".jpg"
    - ".jpeg"
    - ".gif"
    - ".ico"
    - ".pdf"
    - ".zip"
    - ".tar.gz"

  ignore_files:
    - ".DS_Store"
    - "Thumbs.db"
    - "*.tmp"
    - "*.swp"
    - ".git"
    - ".svn"

  # Size limits
  max_file_size: 10MB
  max_template_size: 100MB

# Output configuration
output:
  default_directory: ~/projects
  create_parent_directories: true
  overwrite_existing: confirm  # confirm, always, never
  permissions:
    files: 644
    directories: 755
    executables: 755

  # Backup before overwrite
  backup_on_overwrite: true
  backup_suffix: ".ason-backup"

# User interface
ui:
  colors: true
  emoji: true
  progress_bars: true
  animations: true
  confirm_destructive: true

  # Theme settings
  theme:
    primary: blue
    secondary: green
    warning: yellow
    error: red
    success: green

# Shell completion
completion:
  enabled: true
  cache_enabled: true
  cache_ttl: 3600
  template_suggestions: true
  variable_suggestions: true
  path_completion: true

# Plugin system (future feature)
plugins:
  enabled: true
  auto_install: false
  directory: ~/.ason/plugins
  repositories:
    - "github.com/madstone-tech/ason-plugins"

# External integrations
integrations:
  git:
    auto_init: true
    initial_commit: true
    commit_message: "Initial commit from Ason template"

  editors:
    preferred: "code"  # code, vim, emacs, etc.
    open_after_generation: ask  # always, never, ask

# Logging configuration
logging:
  level: info
  file: ~/.ason/ason.log
  max_size: 10MB
  max_backups: 5
  compress: true

  # Log formats
  format: "json"  # json, text
  timestamp_format: "2006-01-02T15:04:05Z07:00"

# Network settings
network:
  timeout: 30s
  retries: 3
  proxy: ""  # HTTP proxy URL if needed

# Security settings
security:
  trusted_sources:
    - "github.com/madstone-tech"
    - "gitlab.com/yourusername"

  # Template validation
  validate_on_add: true
  strict_validation: false

  # Execution restrictions
  allow_hooks: false
  sandbox_execution: true
EOF
```

## Environment Variables

Ason respects these environment variables:

### Core Settings

```bash
# Registry location
export ASON_REGISTRY_PATH="~/.local/share/ason/templates"

# Default template directory
export ASON_TEMPLATE_DIR="~/templates"

# Configuration file location
export ASON_CONFIG="~/.config/ason/config.toml"

# Debug mode
export ASON_DEBUG="true"

# Quiet mode (minimal output)
export ASON_QUIET="true"

# Verbose mode (detailed output)
export ASON_VERBOSE="true"
```

### Default Variables

```bash
# Personal information
export ASON_AUTHOR="Your Name"
export ASON_EMAIL="your.email@example.com"
export ASON_GITHUB_USERNAME="yourusername"
export ASON_ORGANIZATION="Your Organization"

# Project defaults
export ASON_LICENSE="MIT"
export ASON_VERSION="1.0.0"
export ASON_GO_VERSION="1.25"
export ASON_NODE_VERSION="18"
```

### Output Settings

```bash
# Default output directory
export ASON_OUTPUT_DIR="~/projects"

# Overwrite behavior
export ASON_OVERWRITE="confirm"  # confirm, always, never

# Create parent directories
export ASON_CREATE_PARENT_DIRS="true"
```

### UI and Completion

```bash
# Disable colors
export ASON_NO_COLOR="true"

# Disable emoji
export ASON_NO_EMOJI="true"

# Disable completion cache
export ASON_NO_COMPLETION_CACHE="true"

# Completion cache TTL (seconds)
export ASON_COMPLETION_TTL="3600"
```

## Project-Specific Configuration

### Local Configuration

Create project-specific configuration:

```bash
# In your project directory
cat > .ason/config.yaml << 'EOF'
# Project-specific Ason configuration

defaults:
  author: "Project Team"
  organization: "MyCompany"
  license: "Apache-2.0"
  go_version: "1.25"

template_dirs:
  - ./templates
  - ../shared-templates

output:
  default_dir: ./generated

# Project-specific variables
variables:
  company_name: "MyCompany Inc."
  project_prefix: "mycompany"
  api_version: "v1"
EOF
```

### Template-Specific Configuration

Override settings for specific templates:

```bash
cat > ~/.ason/config.yaml << 'EOF'
# Template-specific overrides
template_overrides:
  go-service:
    defaults:
      go_version: "1.25"
      port: "8080"
      database: "postgresql"

  react-app:
    defaults:
      node_version: "18"
      package_manager: "npm"
      typescript: "true"

  python-service:
    defaults:
      python_version: "3.11"
      framework: "fastapi"
      testing: "pytest"
EOF
```

## User-Specific Customization

### Personal Defaults

Set up your personal defaults:

```bash
cat > ~/.ason/defaults.yaml << 'EOF'
# Personal default values
author: "Your Full Name"
email: "your.email@example.com"
github_username: "yourusername"
gitlab_username: "yourusername"
organization: "Your Organization"
company: "Your Company"
website: "https://yourwebsite.com"

# Coding preferences
preferred_license: "MIT"
default_go_version: "1.25"
default_node_version: "18"
default_python_version: "3.11"

# Project structure preferences
use_src_dir: true
include_tests: true
include_docs: true
include_ci: true

# Git preferences
git_auto_init: true
git_initial_commit: true
default_branch: "main"
EOF
```

### Shell Integration

Add to your shell configuration:

```bash
# ~/.bashrc or ~/.zshrc

# Ason configuration
export ASON_AUTHOR="Your Name"
export ASON_EMAIL="your.email@example.com"
export ASON_GITHUB_USERNAME="yourusername"

# Ason aliases
alias anew='ason new'
alias alist='ason list'
alias aadd='ason add'
alias avalidate='ason validate'

# Quick project creation functions
function create-go-service() {
  ason new go-service "$1" \
    --var service_name="$1" \
    --var author="$ASON_AUTHOR" \
    --var module_name="github.com/$ASON_GITHUB_USERNAME/$1"
}

function create-react-app() {
  ason new react-app "$1" \
    --var app_name="$1" \
    --var author="$ASON_AUTHOR" \
    --var typescript="true"
}

# Completion for custom functions
complete -F _ason_completion create-go-service
complete -F _ason_completion create-react-app
```

## Advanced Configuration

### Custom Template Processing

```yaml
# Advanced template processing configuration
processing:
  # Custom variable resolvers
  variable_resolvers:
    - type: "git"
      config:
        resolve_author_from_git: true
        resolve_email_from_git: true

    - type: "environment"
      config:
        prefix: "PROJECT_"
        fallback_to_defaults: true

    - type: "prompt"
      config:
        interactive_mode: true
        save_responses: true

  # Custom filters
  custom_filters:
    kebab_case: "{{ . | lower | replace(' ', '-') | replace('_', '-') }}"
    pascal_case: "{{ . | title | replace(' ', '') | replace('-', '') | replace('_', '') }}"
    snake_case: "{{ . | lower | replace(' ', '_') | replace('-', '_') }}"

  # Template hooks
  hooks:
    pre_process:
      - command: "validate-template"
        args: ["--strict"]

    post_process:
      - command: "format-code"
        args: ["--language", "auto"]
      - command: "git-init"
        condition: "{{ git_init | default(true) }}"
```

### Integration with External Tools

```yaml
# External tool integration
integrations:
  # Version control
  git:
    auto_init: true
    initial_commit: true
    commit_message: "🎭 Initial commit from Ason template: {{ template_name }}"
    branch: "{{ default_branch | default('main') }}"

  # Package managers
  package_managers:
    npm:
      auto_install: ask
      scripts:
        - "npm install"
        - "npm audit fix"

    go:
      auto_tidy: true
      scripts:
        - "go mod tidy"
        - "go mod download"

  # IDEs and editors
  editors:
    vscode:
      open_after_generation: ask
      workspace_settings: true
      recommended_extensions: true

    jetbrains:
      open_after_generation: false
      project_files: true

  # CI/CD platforms
  cicd:
    github_actions:
      auto_setup: true
      workflows: ["test", "build", "deploy"]

    gitlab_ci:
      auto_setup: false
      template: "default"
```

### Conditional Configuration

```yaml
# Conditional configuration based on context
conditional:
  # Work vs personal projects
  - condition: "{{ output_path | contains('/work/') }}"
    config:
      defaults:
        author: "Work Name"
        email: "work.email@company.com"
        license: "Proprietary"
        organization: "Company Name"

  - condition: "{{ output_path | contains('/personal/') }}"
    config:
      defaults:
        author: "Personal Name"
        email: "personal.email@domain.com"
        license: "MIT"
        organization: "Personal"

  # Language-specific settings
  - condition: "{{ template_type == 'go' }}"
    config:
      processing:
        format_on_save: true
        run_gofmt: true
        run_goimports: true

  - condition: "{{ template_type == 'node' }}"
    config:
      processing:
        format_on_save: true
        run_prettier: true
        run_eslint: true
```

## Configuration Validation

### Validate Your Configuration

```bash
# Check configuration syntax
ason config validate

# Show current configuration
ason config show

# Show configuration sources
ason config sources

# Test variable resolution
ason config test-vars

# Check template overrides
ason config check-overrides
```

### Configuration File Schema

```yaml
# Configuration file schema (for validation)
$schema: "https://schemas.ason.dev/config/v1.yaml"

# Your configuration here...
```

## Performance Tuning

### Cache Configuration

```yaml
# Performance optimization
caching:
  template_cache:
    enabled: true
    max_size: 100MB
    ttl: 24h

  completion_cache:
    enabled: true
    max_entries: 1000
    ttl: 1h

  variable_cache:
    enabled: true
    max_entries: 500
    ttl: 30m

# Memory limits
limits:
  max_template_size: 100MB
  max_file_size: 10MB
  max_memory_usage: 256MB
  concurrent_jobs: 4
```

### Network Optimization

```yaml
# Network performance
network:
  timeout: 30s
  retries: 3
  retry_delay: 1s
  max_retry_delay: 10s

  # Connection pooling
  pool_size: 10
  keep_alive: 30s

  # Compression
  compression: true
  compression_level: 6
```

## Security Configuration

### Security Settings

```yaml
# Security configuration
security:
  # Trusted template sources
  trusted_sources:
    - "github.com/madstone-tech"
    - "gitlab.com/trusted-org"

  # Execution restrictions
  allow_execution: false
  allow_network_access: false
  allow_file_system_write: true

  # Template validation
  strict_validation: true
  require_signature: false

  # Sandbox settings
  sandbox:
    enabled: true
    temp_dir: "/tmp/ason-sandbox"
    memory_limit: "256MB"
    time_limit: "5m"
```

### Access Control

```yaml
# Access control
access:
  # User permissions
  user_permissions:
    read_templates: true
    write_templates: true
    delete_templates: false
    admin_functions: false

  # Path restrictions
  allowed_paths:
    - "~/projects"
    - "~/workspace"
    - "/tmp/ason"

  forbidden_paths:
    - "/etc"
    - "/usr"
    - "/var"
    - "~/.ssh"
```

## Troubleshooting Configuration

### Common Issues

1. **Configuration not loaded**
   ```bash
   # Check config file location
   ason config sources

   # Verify file permissions
   ls -la ~/.ason/config.yaml
   ```

2. **Variables not resolved**
   ```bash
   # Test variable resolution
   ason config test-vars

   # Check default values
   ason config show --section defaults
   ```

3. **Performance issues**
   ```bash
   # Check cache status
   ason config show --section caching

   # Clear caches
   ason cache clear
   ```

### Debug Configuration

```bash
# Enable debug mode
export ASON_DEBUG=true

# Verbose configuration loading
ason --verbose config show

# Test configuration with dry run
ason new template output --dry-run --verbose
```

## Migration and Backup

### Backup Configuration

```bash
# Backup your configuration
cp ~/.ason/config.yaml ~/.ason/config.yaml.backup

# Export configuration
ason config export > ason-config-backup.yaml

# Backup entire Ason directory
tar -czf ason-backup.tar.gz ~/.ason/
```

### Migration Between Systems

```bash
# Export from old system
ason config export --include-templates > migration.yaml

# Import to new system
ason config import migration.yaml

# Verify migration
ason config validate
ason list
```

## Best Practices

### 1. Use Version Control

```bash
# Track your configuration
cd ~/.ason
git init
git add config.yaml
git commit -m "Initial Ason configuration"
```

### 2. Environment-Specific Configs

```bash
# Development environment
~/.ason/config.dev.yaml

# Production environment
~/.ason/config.prod.yaml

# Use environment variable to switch
export ASON_CONFIG="~/.ason/config.${ENVIRONMENT}.yaml"
```

### 3. Modular Configuration

```bash
# Split configuration into modules
~/.ason/
├── config.yaml              # Main config
├── defaults.yaml            # Default variables
├── templates.yaml           # Template overrides
├── integrations.yaml        # External tool config
└── personal.yaml           # Personal preferences
```

### 4. Documentation

```bash
# Document your configuration
cat > ~/.ason/README.md << 'EOF'
# My Ason Configuration

## Overview
This configuration is optimized for:
- Go and Node.js development
- Team collaboration
- CI/CD integration

## Custom Variables
- `company_name`: Set to "MyCompany"
- `default_license`: Set to "MIT"

## Template Overrides
- `go-service`: Uses Go 1.25, PostgreSQL default
- `react-app`: Uses Node 18, TypeScript enabled

## Maintenance
- Backup monthly: `make backup-ason-config`
- Update quarterly: `make update-ason-config`
EOF
```

## Next Steps

Now that you've configured Ason:

1. **[Template Creation Guide](../guides/template-creation.md)** - Create your first templates
2. **[Variable Systems Guide](../guides/variables.md)** - Master variable usage
3. **[Advanced Templating Guide](../guides/advanced-templating.md)** - Explore advanced features
4. **[CI/CD Integration Guide](../guides/cicd.md)** - Automate your workflow

---

*The sacred rattle now resonates perfectly with your workflow. Let the transformation begin! 🪇*