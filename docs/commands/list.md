# ※ ason list

> *Reveal the templates ready for invocation in your sacred registry*

The `ason list` command displays all templates available in your local registry, showing what templates are ready to be used for project generation.

## Synopsis

```bash
ason list [flags]
```

## Description

The `list` command reveals all templates that have been prepared for invocation in your local template registry. It provides an overview of your available templates, their sources, and basic information to help you choose the right template for your project.

## Flags

### --format FORMAT
Specify the output format for the template list.

**Available formats:**
- `table` (default) - Formatted table output
- `json` - JSON format for scripting
- `yaml` - YAML format for configuration

```bash
# Default table format
ason list

# JSON output for scripts
ason list --format json

# YAML output
ason list --format yaml
```

### --filter PATTERN
Filter templates by name or description pattern.

```bash
# Show only web-related templates
ason list --filter web

# Show templates containing "service"
ason list --filter service

# Case-insensitive pattern matching
ason list --filter API
```

### --sort FIELD
Sort templates by specific field.

**Available sort fields:**
- `name` (default) - Sort by template name
- `date` - Sort by date added
- `size` - Sort by template size
- `type` - Sort by template type

```bash
# Sort by name (default)
ason list --sort name

# Sort by most recently added
ason list --sort date

# Sort by template size
ason list --sort size
```

### --reverse
Reverse the sort order.

```bash
# Newest templates first
ason list --sort date --reverse

# Largest templates first
ason list --sort size --reverse
```

### Global Flags
- `-h, --help` - Show help for the command
- `-v, --version` - Show Ason version

## Examples

### Basic Usage

```bash
# List all templates
ason list

# List with JSON output
ason list --format json

# Filter web templates
ason list --filter web
```

### Advanced Filtering and Sorting

```bash
# Find API-related templates
ason list --filter api

# Show newest templates first
ason list --sort date --reverse

# Find large templates
ason list --sort size --reverse

# Get machine-readable output
ason list --format json | jq '.templates[].name'
```

## Output Formats

### Table Format (Default)

```
※ Templates ready for invocation:

┌─────────────────┬──────────────────────────────────┬──────────┬─────────────┐
│ Name            │ Description                      │ Size     │ Added       │
├─────────────────┼──────────────────────────────────┼──────────┼─────────────┤
│ react-app       │ Modern React application         │ 45.2 KB  │ 2 days ago  │
│ go-service      │ Go microservice with gRPC        │ 23.1 KB  │ 1 week ago  │
│ terraform-aws   │ AWS infrastructure template      │ 78.5 KB  │ 3 days ago  │
│ docs-site       │ Documentation site with MkDocs   │ 12.8 KB  │ 5 days ago  │
│ node-api        │ Node.js REST API                 │ 34.7 KB  │ 1 day ago   │
└─────────────────┴──────────────────────────────────┴──────────┴─────────────┘

💡 Use 'ason new TEMPLATE OUTPUT_DIR' to create a project
💡 Use 'ason add' to prepare more templates for invocation
```

### Empty Registry

```
※ The registry echoes with silence...

No templates ready for invocation.

💡 Prepare templates for transformation:
   ason add my-template /path/to/template

💡 Find community templates:
   Visit https://github.com/madstone-tech/ason-templates
```

### JSON Format

```json
{
  "templates": [
    {
      "name": "react-app",
      "description": "Modern React application",
      "path": "/Users/user/.ason/templates/react-app",
      "size": 46285,
      "files": 23,
      "added": "2023-12-01T10:30:00Z",
      "type": "web",
      "variables": [
        "project_name",
        "port",
        "typescript"
      ]
    },
    {
      "name": "go-service",
      "description": "Go microservice with gRPC",
      "path": "/Users/user/.ason/templates/go-service",
      "size": 23654,
      "files": 15,
      "added": "2023-11-24T15:45:00Z",
      "type": "backend",
      "variables": [
        "service_name",
        "port",
        "database"
      ]
    }
  ],
  "total": 2,
  "registry_path": "/Users/user/.ason/templates"
}
```

### YAML Format

```yaml
templates:
  - name: react-app
    description: Modern React application
    path: /Users/user/.ason/templates/react-app
    size: 46285
    files: 23
    added: 2023-12-01T10:30:00Z
    type: web
    variables:
      - project_name
      - port
      - typescript
  - name: go-service
    description: Go microservice with gRPC
    path: /Users/user/.ason/templates/go-service
    size: 23654
    files: 15
    added: 2023-11-24T15:45:00Z
    type: backend
    variables:
      - service_name
      - port
      - database
total: 2
registry_path: /Users/user/.ason/templates
```

## Template Information

Each template entry shows:

### Basic Information
- **Name**: Template identifier for use with `ason new`
- **Description**: Brief description of the template's purpose
- **Size**: Total size of template files
- **Added**: When the template was added to registry

### Extended Information (JSON/YAML)
- **Path**: Full path to template in registry
- **Files**: Number of files in template
- **Type**: Template category (web, backend, infrastructure, etc.)
- **Variables**: List of template variables

## Filtering Examples

### By Template Type
```bash
# Web applications
ason list --filter "web\|react\|vue\|angular"

# Backend services
ason list --filter "service\|api\|backend"

# Infrastructure
ason list --filter "terraform\|docker\|k8s\|infrastructure"

# Documentation
ason list --filter "docs\|documentation\|readme"
```

### By Technology
```bash
# JavaScript/Node.js templates
ason list --filter "node\|js\|react\|vue"

# Go templates
ason list --filter "go\|golang"

# Python templates
ason list --filter "python\|django\|flask"
```

### By Project Size
```bash
# Small templates (< 20KB)
ason list --sort size | head -5

# Large templates (> 50KB)
ason list --sort size --reverse | head -5
```

## Scripting with ason list

### Extract Template Names
```bash
# Get all template names
ason list --format json | jq -r '.templates[].name'

# Get web templates only
ason list --format json | jq -r '.templates[] | select(.type == "web") | .name'

# Count templates
ason list --format json | jq '.total'
```

### Check Template Existence
```bash
# Check if template exists
if ason list --format json | jq -e '.templates[] | select(.name == "my-template")' > /dev/null; then
  echo "Template exists"
else
  echo "Template not found"
fi
```

### Template Validation Script
```bash
#!/bin/bash
# validate-templates.sh

templates=$(ason list --format json | jq -r '.templates[].name')

for template in $templates; do
  echo "Validating $template..."
  ason validate "$template" || echo "❌ $template failed validation"
done
```

## Common Use Cases

### 1. Browse Available Templates
```bash
# Quick overview
ason list

# Detailed view with descriptions
ason list --format table
```

### 2. Find Specific Templates
```bash
# Find web templates
ason list --filter web

# Find recently added templates
ason list --sort date --reverse | head -5
```

### 3. Template Management
```bash
# Count templates
ason list --format json | jq '.total'

# Find largest templates
ason list --sort size --reverse

# Export template list
ason list --format yaml > my-templates.yaml
```

### 4. Integration with Scripts
```bash
# Generate project from first available template
template=$(ason list --format json | jq -r '.templates[0].name')
ason new "$template" "my-project"

# List templates for shell completion
ason list --format json | jq -r '.templates[].name' | sort
```

## Registry Location

Templates are stored in your local registry:

- **Linux/macOS**: `~/.ason/templates/`
- **Windows**: `%USERPROFILE%\.ason\templates\`

Each template is stored as a directory within the registry.

## Troubleshooting

### Registry Not Found
```
❌ Template registry not found
💡 Initialize with: ason add my-first-template /path/to/template
```

### Permission Issues
```
❌ Cannot access template registry
💡 Check permissions for ~/.ason/templates/
```

### Corrupted Templates
```
❌ Template 'broken-template' appears corrupted
💡 Remove and re-add: ason remove broken-template && ason add broken-template /path/to/source
```

## Related Commands

- [`ason new`](new.md) - Create projects from templates
- [`ason add`](add.md) - Add templates to registry
- [`ason remove`](remove.md) - Remove templates from registry
- [`ason validate`](validate.md) - Validate template configuration

## See Also

- [Template Registry Guide](../guides/registry.md)
- [Template Creation Guide](../guides/template-creation.md)
- [Quick Start Guide](../getting-started/quick-start.md)

---

*The registry holds the wisdom of many templates. Choose wisely for your transformation! 🪇*