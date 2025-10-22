# ※ ason validate

> *Ensure the sacred templates are pure and ready for transformation*

The `ason validate` command verifies template integrity, syntax, and configuration to ensure templates are ready for project generation.

## Synopsis

```bash
ason validate [TEMPLATE_OR_PATH] [flags]
```

## Description

The `validate` command performs comprehensive verification of template structure, configuration files, variable definitions, and template syntax to ensure everything is properly configured for successful project generation.

## Arguments

### TEMPLATE_OR_PATH (Optional)
The template to validate. Can be:
- **Template name**: `react-app` (from registry)
- **Local path**: `./my-template` or `/path/to/template`
- **Configuration file**: `ason.toml`
- **No argument**: Validates all templates in registry

## Flags

### --strict
Enable strict validation with additional checks.

```bash
# Strict validation
ason validate my-template --strict
```

### --format FORMAT
Specify output format for validation results.

**Available formats:**
- `text` (default) - Human-readable output
- `json` - JSON format for scripting
- `junit` - JUnit XML for CI integration

```bash
# JSON output
ason validate --format json

# JUnit XML for CI
ason validate --format junit > validation-results.xml
```

### --fix
Automatically fix common issues where possible.

```bash
# Fix issues automatically
ason validate my-template --fix
```

### --check CATEGORY
Validate specific categories only.

**Available categories:**
- `structure` - Directory and file structure
- `config` - Configuration file syntax
- `variables` - Variable definitions and usage
- `templates` - Template file syntax
- `metadata` - Template metadata

```bash
# Check only template syntax
ason validate my-template --check templates

# Check multiple categories
ason validate my-template --check structure,config,variables
```

### --ignore-warnings
Show only errors, ignore warnings.

```bash
# Errors only
ason validate my-template --ignore-warnings
```

### Global Flags
- `-h, --help` - Show help for the command
- `-v, --version` - Show Ason version

## Examples

### Basic Validation

```bash
# Validate specific template
ason validate react-app

# Validate local template
ason validate ./my-template

# Validate configuration file
ason validate ason.toml

# Validate all registry templates
ason validate
```

### Advanced Validation

```bash
# Strict validation with fixes
ason validate my-template --strict --fix

# JSON output for scripting
ason validate --format json > validation.json

# Check specific aspects
ason validate complex-template --check variables,templates

# CI-friendly validation
ason validate --format junit --ignore-warnings
```

### Batch Validation

```bash
# Validate all templates
ason validate

# Validate multiple specific templates
for template in react-app vue-app go-service; do
  ason validate "$template"
done

# Validate with error collection
ason validate --format json | jq '.templates[] | select(.status == "error")'
```

## Validation Categories

### 1. Structure Validation
Checks template directory structure and required files.

**Validates:**
- Template directory exists and is readable
- Contains at least one template file
- Directory structure is logical
- No circular references or invalid links

**Common issues:**
```
❌ Template directory is empty
❌ Template contains no processable files
❌ Broken symbolic links detected
❌ Invalid file permissions
```

### 2. Configuration Validation
Verifies `ason.toml` syntax and content.

**Validates:**
- TOML syntax is correct
- Required fields are present
- Field types are correct
- Version compatibility

**Example issues:**
```
❌ ason.toml syntax error on line 15
❌ Missing required field: 'name'
❌ Invalid variable type: expected string, got number
❌ Unsupported ason version: 2.0.0
```

### 3. Variable Validation
Checks variable definitions and usage consistency.

**Validates:**
- Variables defined in config are used in templates
- Variables used in templates are defined
- Variable types are consistent
- Default values are valid

**Example issues:**
```
❌ Undefined variable 'project_title' used in template.html
❌ Variable 'unused_var' defined but never used
❌ Variable 'port' expects number but default is string
❌ Required variable 'author' has no default value
```

### 4. Template Syntax Validation
Verifies Pongo2 template syntax in all template files.

**Validates:**
- Template syntax is correct
- Filters are valid and available
- Control structures are properly closed
- Variable references are syntactically correct

**Example issues:**
```
❌ Template syntax error in src/main.js:15
❌ Unknown filter 'invalid_filter' in config.json:8
❌ Unclosed if statement in README.md:25
❌ Invalid variable reference: {{ project-name }}
```

### 5. Metadata Validation
Checks template metadata and documentation.

**Validates:**
- README.md exists and is well-formed
- Template description is present
- Author information is provided
- Version format is valid

**Example issues:**
```
⚠️  No README.md found
⚠️  Template description is empty
⚠️  No author information provided
⚠️  Version format is not semantic (expected x.y.z)
```

## Output Formats

### Text Format (Default)

```
※ Validating template: react-app

✅ Structure Validation
   ✓ Template directory exists
   ✓ Contains 23 processable files
   ✓ Directory structure is valid

✅ Configuration Validation
   ✓ ason.toml syntax is correct
   ✓ All required fields present
   ✓ Variable definitions valid

❌ Variable Validation
   ✗ Undefined variable 'app_title' used in src/App.js:12
   ✗ Variable 'unused_port' defined but never used

✅ Template Syntax Validation
   ✓ All 23 template files valid
   ✓ No syntax errors detected

⚠️  Metadata Validation
   ⚠ No README.md found in template
   ⚠ Template description could be more descriptive

🔮 Validation Summary:
   ✅ Passed: 3 categories
   ❌ Failed: 1 category
   ⚠️  Warnings: 1 category

💡 Template needs fixes before reliable use
💡 Use --fix to automatically resolve some issues
```

### JSON Format

```json
{
  "template": "react-app",
  "path": "/Users/user/.ason/templates/react-app",
  "timestamp": "2023-12-01T15:30:45Z",
  "overall_status": "warnings",
  "categories": {
    "structure": {
      "status": "passed",
      "checks": [
        {"name": "directory_exists", "status": "passed", "message": "Template directory exists"},
        {"name": "has_files", "status": "passed", "message": "Contains 23 processable files"},
        {"name": "structure_valid", "status": "passed", "message": "Directory structure is valid"}
      ]
    },
    "config": {
      "status": "passed",
      "checks": [
        {"name": "toml_syntax", "status": "passed", "message": "ason.toml syntax is correct"},
        {"name": "required_fields", "status": "passed", "message": "All required fields present"}
      ]
    },
    "variables": {
      "status": "failed",
      "checks": [
        {
          "name": "undefined_variable",
          "status": "failed",
          "message": "Undefined variable 'app_title' used in src/App.js:12",
          "file": "src/App.js",
          "line": 12,
          "fixable": false
        },
        {
          "name": "unused_variable",
          "status": "failed",
          "message": "Variable 'unused_port' defined but never used",
          "fixable": true
        }
      ]
    },
    "templates": {
      "status": "passed",
      "checks": [
        {"name": "syntax_validation", "status": "passed", "message": "All 23 template files valid"}
      ]
    },
    "metadata": {
      "status": "warnings",
      "checks": [
        {
          "name": "readme_missing",
          "status": "warning",
          "message": "No README.md found in template",
          "fixable": true
        },
        {
          "name": "description_quality",
          "status": "warning",
          "message": "Template description could be more descriptive",
          "fixable": false
        }
      ]
    }
  },
  "summary": {
    "passed": 3,
    "failed": 1,
    "warnings": 1,
    "total_checks": 12,
    "fixable_issues": 2
  }
}
```

### JUnit XML Format

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="ason-validate" tests="5" failures="1" errors="0" skipped="0" time="0.145">
  <testcase classname="react-app" name="structure" time="0.023"/>
  <testcase classname="react-app" name="config" time="0.015"/>
  <testcase classname="react-app" name="variables" time="0.067">
    <failure message="Variable validation failed" type="ValidationError">
      Undefined variable 'app_title' used in src/App.js:12
      Variable 'unused_port' defined but never used
    </failure>
  </testcase>
  <testcase classname="react-app" name="templates" time="0.034"/>
  <testcase classname="react-app" name="metadata" time="0.006"/>
</testsuite>
```

## Auto-Fix Capabilities

The `--fix` flag can automatically resolve certain issues:

### Fixable Issues

```bash
# Issues that can be auto-fixed
ason validate my-template --fix
```

**Auto-fixable:**
- Remove unused variable definitions
- Generate basic README.md template
- Fix common TOML formatting issues
- Standardize variable naming conventions
- Add missing metadata fields

### Fix Examples

**Before fix:**
```yaml
# ason.yaml
name: my template
variables:
  - name: project_name
  - name: unused_variable  # This will be removed
  - name: Project-Title    # This will be standardized
```

**After fix:**
```yaml
# ason.yaml
name: "my template"
description: "Auto-generated template description"
version: "1.0.0"
variables:
  - name: project_name
    description: "Project name"
    required: true
  - name: project_title    # Standardized naming
    description: "Project title"
    required: false
```

**Generated README.md:**
```markdown
# My Template

Auto-generated template for project scaffolding.

## Variables

- `project_name` (required): Project name
- `project_title` (optional): Project title

## Usage

```bash
ason new my-template my-project --var project_name=MyProject
```
```

## Validation in CI/CD

### GitHub Actions Integration

```yaml
name: Validate Templates
on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install Ason
        run: |
          curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Linux_x86_64.tar.gz | tar xz
          sudo mv ason /usr/local/bin/

      - name: Validate Templates
        run: |
          ason validate --format junit > validation-results.xml

      - name: Publish Test Results
        uses: EnricoMi/publish-unit-test-result-action@v2
        if: always()
        with:
          files: validation-results.xml

      - name: Check for Failures
        run: |
          if ason validate --format json | jq -e '.summary.failed > 0'; then
            echo "❌ Template validation failed"
            exit 1
          fi
```

### Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

echo "Validating templates..."
if ! ason validate --ignore-warnings; then
  echo "❌ Template validation failed. Fix issues before committing."
  exit 1
fi
echo "✅ Template validation passed"
```

## Validation Scripts

### Development Validation Script

```bash
#!/bin/bash
# validate-templates.sh

set -e

echo "※ Starting comprehensive template validation..."

# Validate all templates
echo "Validating all templates..."
ason validate --format json > validation-results.json

# Check for failures
failures=$(jq '.summary.failed' validation-results.json)
if [ "$failures" -gt 0 ]; then
  echo "❌ $failures validation failures found"
  jq -r '.categories[] | select(.status == "failed") | .checks[] | select(.status == "failed") | .message' validation-results.json
  exit 1
fi

# Check for warnings
warnings=$(jq '.summary.warnings' validation-results.json)
if [ "$warnings" -gt 0 ]; then
  echo "⚠️  $warnings warnings found"
  jq -r '.categories[] | select(.status == "warnings") | .checks[] | select(.status == "warning") | .message' validation-results.json
fi

echo "✅ All templates validated successfully"
```

### Template Quality Check

```bash
#!/bin/bash
# quality-check.sh

echo "Running template quality checks..."

# Strict validation
if ! ason validate --strict --ignore-warnings; then
  echo "❌ Strict validation failed"
  exit 1
fi

# Check template coverage
templates_count=$(ason list --format json | jq '.total')
if [ "$templates_count" -eq 0 ]; then
  echo "⚠️  No templates in registry"
  exit 1
fi

echo "✅ Quality check passed for $templates_count templates"
```

## Best Practices

### 1. Regular Validation
```bash
# Validate after any template changes
ason validate my-template

# Validate before adding to registry
ason validate ./new-template --strict
```

### 2. CI Integration
```bash
# Include validation in CI pipeline
ason validate --format junit > results.xml

# Fail builds on validation errors
ason validate || exit 1
```

### 3. Development Workflow
```bash
# Fix issues during development
ason validate my-template --fix

# Strict validation before release
ason validate my-template --strict
```

### 4. Template Quality
```bash
# Ensure comprehensive templates
ason validate --check metadata,variables

# Validate documentation
ason validate --check structure,metadata
```

## Common Issues and Solutions

### Variable Issues
```bash
# Find undefined variables
ason validate my-template --check variables

# Fix variable definitions
ason validate my-template --fix
```

### Template Syntax
```bash
# Validate template syntax
ason validate my-template --check templates

# Common syntax errors:
# - Unclosed tags: {% if %} without {% endif %}
# - Invalid filters: {{ var | nonexistent_filter }}
# - Wrong variable syntax: {{ project-name }} instead of {{ project_name }}
```

### Configuration Problems
```bash
# Validate configuration
ason validate my-template --check config

# Common config issues:
# - Invalid TOML syntax
# - Missing required fields
# - Incorrect field types
```

## Related Commands

- [`ason new`](new.md) - Create projects from templates
- [`ason register`](add.md) - Add templates to registry
- [`ason list`](list.md) - List available templates
- [`ason remove`](remove.md) - Remove templates from registry

## See Also

- [Template Creation Guide](../guides/template-creation.md)
- [Variable Systems Guide](../guides/variables.md)
- [Advanced Templating Guide](../guides/advanced-templating.md)
- [Troubleshooting Guide](../troubleshooting/common-issues.md)

---

*Validation ensures the sacred templates are pure and ready for transformation. Trust in the process! 🪇*