# ※ Quick Start Guide

> *Master the sacred rattle's rhythm in just 5 minutes*

This guide gets you up and running with Ason quickly, from installation to creating your first project from a template.

## Prerequisites

- **Ason installed** - See [Installation Guide](installation.md) if not already installed
- **Basic command line knowledge**
- **5 minutes of your time**

## Step 1: Verify Installation

First, let's make sure Ason is properly installed:

```bash
# Check Ason version
ason --version
# Should output: ※ Ason v1.0.0

# See available commands
ason --help
```

If you see the version output, you're ready to proceed!

## Step 2: Check Your Registry

Let's see what templates are available:

```bash
# List available templates
ason list
```

If this is your first time, you'll see:
```
※ The registry echoes with silence...

No templates ready for invocation.

💡 Prepare templates for transformation:
   ason register my-template /path/to/template
```

This is expected! Let's create your first template.

## Step 3: Create Your First Template

Let's create a simple project template:

```bash
# Create a template directory
mkdir -p my-first-template

# Create some template files
cat > my-first-template/README.md << 'EOF'
# {{ project_name }}

{{ description | default:"A project created with Ason" }}

## Getting Started

Welcome to {{ project_name }}!

Author: {{ author | default:"Unknown" }}
Version: {{ version | default:"1.0.0" }}
EOF

cat > my-first-template/package.json << 'EOF'
{
  "name": "{{ project_name | lower | replace(" ", "-") }}",
  "version": "{{ version | default:"1.0.0" }}",
  "description": "{{ description }}",
  "author": "{{ author }}",
  "license": "MIT"
}
EOF

# Create a configuration file (optional)
cat > my-first-template/ason.toml << 'EOF'
name = "My First Template"
description = "A simple template for learning Ason"
version = "1.0.0"
type = "example"

[[variables]]
name = "project_name"
description = "Name of the project"
required = true

[[variables]]
name = "description"
description = "Project description"
required = false
default = "A project created with Ason"

[[variables]]
name = "author"
description = "Project author"
required = false
default = "Unknown"

[[variables]]
name = "version"
description = "Initial version"
required = false
default = "1.0.0"
EOF
```

## Step 4: Add Template to Registry

Now let's add this template to your Ason registry:

```bash
# Add the template to your registry
ason register my-first-template ./my-first-template \
  --description "My first Ason template for learning"

# Verify it was added
ason list
```

You should now see:
```
※ Templates ready for invocation:

┌─────────────────┬──────────────────────────────────┬──────────┬─────────────┐
│ Name            │ Description                      │ Size     │ Added       │
├─────────────────┼──────────────────────────────────┼──────────┼─────────────┤
│ my-first-template│ My first Ason template for learning │ 1.2 KB   │ just now    │
└─────────────────┴──────────────────────────────────┴──────────┴─────────────┘

💡 Use 'ason new TEMPLATE OUTPUT_DIR' to create a project
```

## Step 5: Generate Your First Project

Now comes the exciting part - let's generate a project from your template:

```bash
# Test with dry run first
ason new my-first-template my-awesome-project --dry-run

# Generate the actual project
ason new my-first-template my-awesome-project \
  --var project_name="My Awesome Project" \
  --var description="A fantastic project built with Ason" \
  --var author="Your Name" \
  --var version="0.1.0"
```

You should see output like:
```
※ The ason shakes, preparing transformation...
✨ Catalyst activated for template: my-first-template
🎭 Variables ready for invocation
📿 Processing template files...
💫 Transforming: README.md
💫 Transforming: package.json
💫 Transforming: ason.toml
🔮 The rhythm is complete! Project created at: my-awesome-project
```

## Step 6: Examine the Result

Let's look at what was created:

```bash
# Check the generated project structure
ls -la my-awesome-project/

# Look at the transformed files
cat my-awesome-project/README.md
cat my-awesome-project/package.json
```

You should see your variables have been substituted:

**README.md:**
```markdown
# My Awesome Project

A fantastic project built with Ason

## Getting Started

Welcome to My Awesome Project!

Author: Your Name
Version: 0.1.0
```

**package.json:**
```json
{
  "name": "my-awesome-project",
  "version": "0.1.0",
  "description": "A fantastic project built with Ason",
  "author": "Your Name",
  "license": "MIT"
}
```

## Step 7: Shell Completion (Optional)

Set up shell completion for a better experience:

```bash
# Install completion for your shell
ason completion bash > ~/.local/share/bash-completion/completions/ason

# Or for zsh (Oh My Zsh)
mkdir -p ~/.oh-my-zsh/custom/plugins/ason
ason completion zsh > ~/.oh-my-zsh/custom/plugins/ason/_ason
# Add 'ason' to plugins in ~/.zshrc

# Or for fish
ason completion fish > ~/.config/fish/completions/ason.fish

# Test completion (after restarting shell)
ason <TAB>
ason new <TAB>
```

## Next Steps: Real-World Examples

Now that you understand the basics, let's create some more practical templates:

### Web Application Template

```bash
# Create a React-like template structure
mkdir -p web-app-template/{src,public,tests}

cat > web-app-template/package.json << 'EOF'
{
  "name": "{{ project_name | lower | replace(" ", "-") }}",
  "version": "{{ version | default:"1.0.0" }}",
  "description": "{{ description }}",
  "main": "src/index.js",
  "scripts": {
    "start": "{{ start_command | default:"npm start" }}",
    "build": "{{ build_command | default:"npm run build" }}",
    "test": "{{ test_command | default:"npm test" }}"
  },
  "author": "{{ author }}",
  "license": "{{ license | default:"MIT" }}"
}
EOF

cat > web-app-template/src/index.js << 'EOF'
// {{ project_name }}
// Created by {{ author }}

console.log('Welcome to {{ project_name }}!');

// TODO: Add your application logic here
EOF

cat > web-app-template/README.md << 'EOF'
# {{ project_name }}

{{ description }}

## Quick Start

```bash
# Install dependencies
npm install

# Start development server
{{ start_command | default:"npm start" }}

# Build for production
{{ build_command | default:"npm run build" }}
```

## About

- **Author**: {{ author }}
- **Version**: {{ version }}
- **License**: {{ license | default:"MIT" }}
EOF

# Add to registry
ason register web-app ./web-app-template --type web --description "Modern web application template"
```

### Go Service Template

```bash
# Create a Go service template
mkdir -p go-service-template/{cmd,internal,pkg}

cat > go-service-template/go.mod << 'EOF'
module {{ module_name | default:"github.com/user/project" }}

go {{ go_version | default:"1.25" }}

require (
    // Add your dependencies here
)
EOF

cat > go-service-template/cmd/main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Println("Starting {{ service_name }}...")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from {{ service_name }}!")
    })

    port := "{{ port | default:"8080" }}"
    fmt.Printf("Server listening on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
EOF

cat > go-service-template/README.md << 'EOF'
# {{ service_name }}

{{ description }}

## Running the Service

```bash
# Run locally
go run cmd/main.go

# Build binary
go build -o {{ service_name | lower }} cmd/main.go

# Run binary
./{{ service_name | lower }}
```

The service will be available at: http://localhost:{{ port | default:"8080" }}

## Configuration

- **Service Name**: {{ service_name }}
- **Port**: {{ port | default:"8080" }}
- **Go Version**: {{ go_version | default:"1.25" }}
EOF

# Add to registry
ason register go-service ./go-service-template --type backend --description "Go HTTP service template"
```

### Test the New Templates

```bash
# Generate a web app
ason new web-app my-web-app \
  --var project_name="My Web App" \
  --var description="A cool web application" \
  --var author="Your Name"

# Generate a Go service
ason new go-service user-api \
  --var service_name="UserAPI" \
  --var description="User management API" \
  --var port="8080" \
  --var module_name="github.com/yourname/user-api"

# Check what was created
ls -la my-web-app/
ls -la user-api/
```

## Common Patterns and Tips

### 1. Variable Naming Conventions
```bash
# Use snake_case for consistency
project_name, service_name, database_type

# Use meaningful defaults
{{ port | default:"8080" }}
{{ author | default:"Unknown" }}
{{ license | default:"MIT" }}
```

### 2. Template Organization
```bash
# Organize templates by type
ason register react-app ./templates/web/react --type web
ason register vue-app ./templates/web/vue --type web
ason register go-api ./templates/backend/go --type backend
ason register python-api ./templates/backend/python --type backend
```

### 3. Using Filters
```bash
# Common Pongo2 filters
{{ project_name | lower }}                    # Convert to lowercase
{{ project_name | replace(" ", "-") }}        # Replace spaces with dashes
{{ description | default:"No description" }}  # Provide default value
{{ author | title }}                          # Title case
```

### 4. Template Validation
```bash
# Always validate templates before adding
ason validate ./my-template --strict

# Fix common issues automatically
ason validate ./my-template --fix
```

## Summary

In just 5 minutes, you've learned to:

1. ✅ **Verify Ason installation**
2. ✅ **Create a simple template** with variables
3. ✅ **Add templates to your registry**
4. ✅ **Generate projects** from templates
5. ✅ **Use template variables** and filters
6. ✅ **Set up shell completion** (optional)
7. ✅ **Create real-world templates** for web apps and services

## What's Next?

Now that you've mastered the basics, explore these advanced topics:

- **[Your First Template](first-template.md)** - Deep dive into template creation
- **[Configuration Guide](configuration.md)** - Customize Ason for your workflow
- **[Template Creation Guide](../guides/template-creation.md)** - Advanced templating techniques
- **[Variable Systems Guide](../guides/variables.md)** - Master variable usage
- **[Examples](../examples/)** - Real-world template examples

## Common Commands Reference

```bash
# Template management
ason list                           # List available templates
ason register NAME PATH                  # Add template to registry
ason remove NAME                    # Remove template from registry
ason validate TEMPLATE              # Validate template

# Project generation
ason new TEMPLATE OUTPUT            # Generate project
ason new TEMPLATE OUTPUT --dry-run  # Preview generation
ason new TEMPLATE OUTPUT --var key=value  # Set variables

# Help and information
ason --help                         # Show help
ason COMMAND --help                 # Show command help
ason --version                      # Show version
```

## Cleanup

If you want to clean up the examples created in this guide:

```bash
# Remove generated projects
rm -rf my-awesome-project my-web-app user-api

# Remove example templates from registry
ason remove my-first-template
ason remove web-app
ason remove go-service

# Remove template source directories
rm -rf my-first-template web-app-template go-service-template
```

---

*You've mastered the basic rhythm of the sacred rattle! The templates await your creative transformation! 🪇*