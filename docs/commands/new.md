# ※ ason new

> *Invoke the sacred rhythm to birth new projects from templates*

The `ason new` command is the heart of Ason - it transforms templates into living, breathing projects through the power of the sacred rattle.

## Synopsis

```bash
ason new TEMPLATE OUTPUT_DIR [flags]
```

## Description

The `new` command catalyzes the transformation of a template into a new project. It takes a template (either a local directory path or a template name from your registry) and generates a complete project structure in the specified output directory.

## Arguments

### TEMPLATE
The template to use for generation. Can be:
- **Local path**: `./my-template` or `/path/to/template`
- **Registry name**: `web-app` (when template is in your registry)
- **Relative path**: `../templates/golang-service`

### OUTPUT_DIR
The directory where the new project will be created. If the directory doesn't exist, it will be created automatically.

## Flags

### --dry-run
Perform a dry run without creating any files.

```bash
ason new my-template my-project --dry-run
```

**Use cases:**
- Preview what files will be generated
- Test template structure before actual generation
- Validate template syntax and variables

### --var name=value
Set template variables for substitution.

```bash
ason new web-template my-app --var project_name=MyAwesomeApp --var version=1.0.0
```

**Variable examples:**
- `--var name=MyProject` - Project name
- `--var version=1.2.3` - Version number
- `--var author="John Doe"` - Author name (use quotes for spaces)
- `--var description="A cool project"` - Project description

### Global Flags
- `-h, --help` - Show help for the command
- `-v, --version` - Show Ason version

## Examples

### Basic Project Generation

```bash
# Generate from local template
ason new ./templates/react-app my-new-app

# Generate from registry template (once implemented)
ason new golang-service user-service
```

### Using Variables

```bash
# Single variable
ason new go-template my-service --var service_name=UserService

# Multiple variables
ason new web-template my-blog \
  --var project_name="Personal Blog" \
  --var author="Jane Smith" \
  --var version=0.1.0 \
  --var port=3000
```

### Dry Run Testing

```bash
# Test template before generation
ason new complex-template test-output --dry-run

# Test with variables
ason new api-template test-api \
  --var service=payments \
  --var database=postgres \
  --dry-run
```

### Advanced Examples

```bash
# Generate infrastructure project
ason new terraform-aws infrastructure/production \
  --var environment=prod \
  --var region=us-west-2 \
  --var instance_type=t3.medium

# Generate documentation site
ason new docs-template project-docs \
  --var project_name="Ason Documentation" \
  --var version=1.0.0 \
  --var theme=material
```

## Template Structure

Ason templates use the Pongo2 templating engine (Django/Jinja2-like syntax). Your template directory should contain:

### Template Files
Files with Pongo2 syntax will be processed:

```
my-template/
├── README.md              # Can contain {{ project_name }}
├── package.json           # Can contain {{ version }}
├── src/
│   └── main.js           # Can contain {{ author }}
└── ason.toml             # Template configuration (optional)
```

### Variable Substitution Examples

**README.md template:**
```markdown
# {{ project_name }}

Version: {{ version }}
Author: {{ author }}

## Description

{{ description | default:"A new project created with Ason" }}
```

**package.json template:**
```json
{
  "name": "{{ project_name | lower | replace(" ", "-") }}",
  "version": "{{ version }}",
  "description": "{{ description }}",
  "author": "{{ author }}"
}
```

## Output

### Success Output
```
※ The ason shakes, preparing transformation...
✨ Catalyst activated for template: my-template
🎭 Variables ready for invocation
📿 Processing template files...
💫 Transforming: README.md
💫 Transforming: package.json
💫 Transforming: src/main.js
🔮 The rhythm is complete! Project created at: my-project
```

### Dry Run Output
```
※ The ason shakes, preparing transformation...
[DRY RUN] Would create: my-project/
[DRY RUN] Would process: README.md → my-project/README.md
[DRY RUN] Would process: package.json → my-project/package.json
[DRY RUN] Would process: src/main.js → my-project/src/main.js
🔮 [DRY RUN] The rhythm is prepared! No files were created.
```

## Common Use Cases

### 1. Web Applications
```bash
# React application
ason new react-template my-react-app \
  --var app_name="My React App" \
  --var port=3000

# Vue.js application
ason new vue-template my-vue-app \
  --var app_name="My Vue App" \
  --var router=true
```

### 2. Backend Services
```bash
# Go microservice
ason new go-service user-service \
  --var service_name=UserService \
  --var database=postgresql \
  --var port=8080

# Node.js API
ason new node-api payment-api \
  --var service_name="Payment API" \
  --var database=mongodb
```

### 3. Infrastructure
```bash
# Terraform AWS setup
ason new terraform-aws production-infra \
  --var environment=production \
  --var region=us-east-1

# Docker compose stack
ason new docker-stack monitoring \
  --var stack_name=monitoring \
  --var grafana_port=3000
```

### 4. Documentation
```bash
# Project documentation
ason new docs-template project-docs \
  --var project_name="My Project" \
  --var docs_theme=mkdocs-material

# API documentation
ason new api-docs user-api-docs \
  --var api_name="User API" \
  --var version=v1
```

## Error Handling

### Template Not Found
```
❌ Template not found: nonexistent-template
💡 Available templates: web-app, go-service, docs-template
💡 Use 'ason list' to see all available templates
```

### Output Directory Exists
```
❌ Output directory already exists: my-project
💡 Choose a different name or remove the existing directory
```

### Variable Errors
```
❌ Required variable 'project_name' not provided
💡 Use: --var project_name=YourProjectName
```

### Permission Errors
```
❌ Permission denied creating directory: /root/my-project
💡 Check directory permissions or choose a different location
```

## Tips and Best Practices

### 1. Variable Naming
- Use consistent variable names across templates
- Use snake_case or camelCase consistently
- Provide sensible defaults in templates

### 2. Template Organization
- Keep templates focused and single-purpose
- Use descriptive template names
- Include example variable files

### 3. Testing Templates
- Always test with `--dry-run` first
- Test with different variable combinations
- Validate generated output

### 4. Performance
- Large templates may take longer to process
- Use specific paths rather than wildcards when possible
- Consider breaking large templates into smaller ones

## Related Commands

- [`ason list`](list.md) - List available templates
- [`ason add`](add.md) - Add templates to registry
- [`ason validate`](validate.md) - Validate template configuration

## See Also

- [Template Creation Guide](../guides/template-creation.md)
- [Variable Systems Guide](../guides/variables.md)
- [Advanced Templating Guide](../guides/advanced-templating.md)
- [Web Application Examples](../examples/web-app.md)
- [Go Service Examples](../examples/go-service.md)

---

*The sacred rattle transforms templates into living projects. Let the rhythm guide your creation! 🪇*