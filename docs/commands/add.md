# ※ ason add

> *Prepare templates for transformation by adding them to your sacred registry*

The `ason add` command prepares templates for invocation by adding them to your local template registry, making them available for use with `ason new`.

## Synopsis

```bash
ason add TEMPLATE_NAME SOURCE_PATH [flags]
```

## Description

The `add` command catalyzes the preparation of templates by copying them into your local registry. Once added, templates become available for project generation and can be referenced by name rather than full path.

## Arguments

### TEMPLATE_NAME
The name to assign to the template in your registry. This will be the identifier used with `ason new`.

**Naming conventions:**
- Use lowercase with hyphens: `react-app`, `go-service`
- Be descriptive but concise: `terraform-aws`, `docs-site`
- Avoid spaces and special characters

### SOURCE_PATH
The path to the template directory or archive to add.

**Supported sources:**
- **Local directory**: `/path/to/my-template`
- **Relative path**: `./templates/web-app`
- **Git repository**: `https://github.com/user/template.git` *(future)*
- **Archive file**: `template.tar.gz` *(future)*

## Flags

### --description DESC
Provide a description for the template.

```bash
ason add react-app ./templates/react \
  --description "Modern React application with TypeScript"
```

### --type TYPE
Specify the template type/category.

**Common types:**
- `web` - Web applications
- `backend` - Backend services
- `infrastructure` - Infrastructure as code
- `docs` - Documentation
- `mobile` - Mobile applications
- `desktop` - Desktop applications

```bash
ason add go-service ./templates/go-microservice \
  --type backend \
  --description "Go microservice with gRPC"
```

### --force
Overwrite existing template with the same name.

```bash
# Replace existing template
ason add react-app ./new-react-template --force
```

### --validate
Validate template structure before adding.

```bash
# Validate before adding
ason add my-template ./path/to/template --validate
```

### --dry-run
Show what would be added without actually adding.

```bash
# Preview the add operation
ason add test-template ./my-template --dry-run
```

### Global Flags
- `-h, --help` - Show help for the command
- `-v, --version` - Show Ason version

## Examples

### Basic Template Addition

```bash
# Add local template
ason add react-app ./templates/react-app

# Add with description
ason add go-service ./templates/golang-service \
  --description "Production-ready Go microservice"

# Add with type and description
ason add terraform-aws ./infra/aws-template \
  --type infrastructure \
  --description "AWS infrastructure with Terraform"
```

### Advanced Usage

```bash
# Replace existing template
ason add react-app ./new-react-template --force

# Validate before adding
ason add complex-template ./templates/complex --validate

# Preview addition
ason add test-template ./experimental --dry-run

# Full specification
ason add node-api ./templates/nodejs-api \
  --type backend \
  --description "Node.js REST API with Express and PostgreSQL" \
  --validate
```

### Batch Template Addition

```bash
# Add multiple templates
for template in ./templates/*/; do
  name=$(basename "$template")
  ason add "$name" "$template"
done

# Add with consistent typing
ason add react-spa ./templates/react --type web
ason add vue-app ./templates/vue --type web
ason add angular-app ./templates/angular --type web
```

## Template Structure Requirements

### Basic Structure
A valid template should contain:

```
my-template/
├── template-files...        # Any files/directories
├── ason.yaml               # Optional configuration
└── README.md               # Optional documentation
```

### Optional Configuration (ason.yaml)

```yaml
# Template metadata
name: "My Template"
description: "A comprehensive template for..."
version: "1.0.0"
author: "Your Name"
type: "web"

# Template variables
variables:
  - name: project_name
    description: "Name of the project"
    required: true
    default: "my-project"

  - name: version
    description: "Initial version"
    required: false
    default: "0.1.0"

  - name: author
    description: "Project author"
    required: true

# Files to ignore during processing
ignore:
  - "*.tmp"
  - ".git/"
  - "node_modules/"

# Template hooks (future feature)
hooks:
  pre_generate: "scripts/pre-generate.sh"
  post_generate: "scripts/post-generate.sh"
```

### Template Variables
Use Pongo2 syntax for variables:

```
{{ project_name }}           # Simple variable
{{ project_name | lower }}   # With filter
{{ description | default:"No description" }}  # With default
```

## Output

### Success Output
```
※ The ason prepares to embrace new wisdom...
✨ Analyzing template: ./templates/react-app
📿 Validating template structure...
💫 Template structure confirmed
🎭 Copying template to registry...
🔮 Template 'react-app' added to registry successfully!

💡 Use it with: ason new react-app my-project
```

### Dry Run Output
```
※ The ason prepares to embrace new wisdom...
[DRY RUN] Would analyze: ./templates/react-app
[DRY RUN] Would validate template structure
[DRY RUN] Would copy to: ~/.ason/templates/react-app
[DRY RUN] Would register as: react-app
🔮 [DRY RUN] Template ready for addition. Use without --dry-run to add.
```

### Validation Output
```
※ Validating template structure...
✅ Template directory exists
✅ Contains template files
✅ ason.yaml is valid
✅ Variables are properly defined
✅ No syntax errors in templates
🔮 Template validation passed!
```

## Error Handling

### Template Name Already Exists
```
❌ Template 'react-app' already exists in registry
💡 Use --force to overwrite: ason add react-app ./new-template --force
💡 Or choose a different name: ason add react-app-v2 ./new-template
```

### Source Path Not Found
```
❌ Source path not found: ./nonexistent-template
💡 Check the path and try again
💡 Use absolute path if relative path fails
```

### Invalid Template Structure
```
❌ Invalid template structure in ./bad-template
💡 Template must contain at least one file
💡 Check ason.yaml syntax if present
```

### Permission Issues
```
❌ Permission denied accessing registry
💡 Check permissions for ~/.ason/templates/
💡 Try running with appropriate permissions
```

### Validation Errors
```
❌ Template validation failed:
  - ason.yaml syntax error on line 5
  - Undefined variable 'missing_var' in template.html
  - Invalid filter 'unknown_filter' in config.json
💡 Fix these issues and try again
```

## Registry Management

### Registry Location
Templates are stored following XDG Base Directory specification:
- **Linux/macOS**: `~/.local/share/ason/templates/` (XDG_DATA_HOME/ason)
- **Windows**: `%LOCALAPPDATA%\ason\templates\`

### Registry Structure
```
~/.local/share/ason/
├── templates/
│   ├── react-app/           # Template directory
│   │   ├── src/
│   │   ├── package.json
│   │   └── ason.toml        # Template configuration
│   ├── go-service/
│   │   ├── cmd/
│   │   ├── internal/
│   │   └── ason.toml
│   └── registry.toml        # Registry metadata
└── config.toml              # Ason configuration (if needed)
```

### Template Metadata
Registry metadata is stored in `registry.toml`:

```toml
updated = 2023-12-01T10:30:00Z

[templates.react-app]
name = "react-app"
path = "/home/user/.local/share/ason/templates/react-app"
description = "Modern React application"
source = "/original/path/to/template"
type = "web"
size = 1024
files = 23
added = 2023-12-01T10:30:00Z
variables = ["project_name", "author", "license"]

[templates.go-service]
name = "go-service"
path = "/home/user/.local/share/ason/templates/go-service"
description = "Go microservice with gRPC"
source = "/path/to/go-template"
type = "backend"
size = 23654
files = 15
added = 2023-11-24T15:45:00Z
variables = ["service_name", "module_path"]
```

## Best Practices

### 1. Template Organization
```bash
# Use consistent naming
ason add react-spa ./templates/react-spa --type web
ason add vue-spa ./templates/vue-spa --type web

# Provide good descriptions
ason add node-api ./templates/nodejs-api \
  --description "Express.js API with PostgreSQL and Docker"

# Specify types for better organization
ason add terraform-aws ./infra/aws --type infrastructure
```

### 2. Template Validation
```bash
# Always validate complex templates
ason add complex-template ./templates/complex --validate

# Test with dry-run first
ason add experimental ./templates/experimental --dry-run
```

### 3. Template Versioning
```bash
# Version your templates
ason add react-app-v1 ./templates/react-v1
ason add react-app-v2 ./templates/react-v2

# Or use descriptive names
ason add react-hooks ./templates/react-with-hooks
ason add react-class ./templates/react-with-classes
```

### 4. Template Documentation
Include comprehensive documentation in your templates:

```
my-template/
├── README.md               # Template usage guide
├── VARIABLES.md           # Variable documentation
├── EXAMPLES.md            # Example generations
└── ason.yaml             # Template configuration
```

## Common Use Cases

### 1. Personal Template Library
```bash
# Add your common templates
ason add my-react ./personal-templates/react
ason add my-go-api ./personal-templates/go-api
ason add my-terraform ./personal-templates/terraform
```

### 2. Team Template Sharing
```bash
# Add team templates
ason add company-frontend ./team-templates/frontend --type web
ason add company-backend ./team-templates/backend --type backend
ason add company-infra ./team-templates/infrastructure --type infrastructure
```

### 3. Template Development
```bash
# Test template during development
ason add test-template ./work-in-progress --dry-run
ason add test-template ./work-in-progress --validate
ason add test-template ./work-in-progress --force
```

### 4. Template Migration
```bash
# Migrate from old structure
ason add new-react ./templates/react-new --force
ason remove old-react
```

## Integration Examples

### CI/CD Template Addition
```yaml
# .github/workflows/add-templates.yml
name: Add Templates to Registry
on:
  push:
    paths: ['templates/**']

jobs:
  add-templates:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install Ason
        run: |
          curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Linux_x86_64.tar.gz | tar xz
          sudo mv ason /usr/local/bin/
      - name: Add Templates
        run: |
          for template in templates/*/; do
            name=$(basename "$template")
            ason add "$name" "$template" --validate
          done
```

### Script Template Management
```bash
#!/bin/bash
# manage-templates.sh

# Add all templates from a directory
add_templates() {
  local template_dir="$1"
  for template in "$template_dir"/*/; do
    if [[ -d "$template" ]]; then
      local name=$(basename "$template")
      echo "Adding template: $name"
      ason add "$name" "$template" --validate || echo "Failed to add $name"
    fi
  done
}

add_templates "./my-templates"
```

## Related Commands

- [`ason new`](new.md) - Create projects from templates
- [`ason list`](list.md) - List available templates
- [`ason remove`](remove.md) - Remove templates from registry
- [`ason validate`](validate.md) - Validate template configuration

## See Also

- [Template Registry Guide](../guides/registry.md)
- [Template Creation Guide](../guides/template-creation.md)
- [Getting Started Guide](../getting-started/quick-start.md)

---

*Add wisdom to your registry, and let the templates serve your creative transformation! 🪇*