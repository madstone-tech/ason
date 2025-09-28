# ※ ason remove

> *Release templates from your registry when their service is complete*

The `ason remove` command removes templates from your local registry, freeing space and cleaning up templates that are no longer needed.

## Synopsis

```bash
ason remove TEMPLATE_NAME [flags]
```

## Description

The `remove` command cleanses your registry by removing templates that are no longer needed. This operation removes the template files from your local registry while preserving the original source files.

## Arguments

### TEMPLATE_NAME
The name of the template to remove from the registry. Use `ason list` to see available templates.

## Flags

### --force
Remove template without confirmation prompt.

```bash
# Remove without confirmation
ason remove old-template --force
```

### --dry-run
Show what would be removed without actually removing.

```bash
# Preview removal
ason remove test-template --dry-run
```

### --backup
Create a backup before removing.

```bash
# Create backup before removal
ason remove important-template --backup
```

### --backup-dir DIR
Specify backup directory location.

```bash
# Backup to specific directory
ason remove template-name --backup --backup-dir ./template-backups
```

### Global Flags
- `-h, --help` - Show help for the command
- `-v, --version` - Show Ason version

## Examples

### Basic Template Removal

```bash
# Remove template (with confirmation)
ason remove old-template

# Remove without confirmation
ason remove test-template --force

# Preview removal
ason remove experimental --dry-run
```

### Safe Removal with Backup

```bash
# Create backup before removing
ason remove important-template --backup

# Backup to specific location
ason remove legacy-template --backup --backup-dir ./backups

# Force remove with backup
ason remove old-version --force --backup
```

### Batch Removal

```bash
# Remove multiple templates
for template in old-react old-vue deprecated-api; do
  ason remove "$template" --force
done

# Remove all templates matching pattern
ason list --format json | \
  jq -r '.templates[] | select(.name | startswith("test-")) | .name' | \
  xargs -I {} ason remove {} --force
```

## Output

### Interactive Removal
```
※ The ason prepares to release template from registry...

Template: react-old-version
Description: Outdated React template
Size: 45.2 KB
Files: 23
Added: 2 weeks ago

⚠️  This action cannot be undone.
🔮 Remove template 'react-old-version' from registry? [y/N]: y

✨ Removing template from registry...
💫 Cleaning up template files...
🎭 Template 'react-old-version' has been released from registry

💡 Original source files remain unchanged
```

### Forced Removal
```
※ The ason prepares to release template from registry...
✨ Removing template 'test-template'...
💫 Cleaning up template files...
🔮 Template 'test-template' removed successfully!
```

### Dry Run Output
```
※ The ason prepares to release template from registry...
[DRY RUN] Would remove template: experimental-template
[DRY RUN] Would delete: ~/.ason/templates/experimental-template/
[DRY RUN] Would clean registry metadata
[DRY RUN] Size to be freed: 23.4 KB
🔮 [DRY RUN] Template ready for removal. Use without --dry-run to remove.
```

### Backup Creation
```
※ Creating backup before removal...
✨ Backing up template to: ~/.ason/backups/react-app-2023-12-01-143022.tar.gz
💫 Backup created successfully (45.2 KB)
🎭 Proceeding with template removal...
✨ Template 'react-app' removed successfully!

💡 Backup available at: ~/.ason/backups/react-app-2023-12-01-143022.tar.gz
💡 Restore with: ason add react-app ~/.ason/backups/react-app-2023-12-01-143022.tar.gz
```

## Error Handling

### Template Not Found
```
❌ Template 'nonexistent' not found in registry
💡 Use 'ason list' to see available templates
💡 Check template name spelling
```

### Registry Permission Issues
```
❌ Permission denied removing template 'template-name'
💡 Check permissions for ~/.ason/templates/
💡 Try running with appropriate permissions
```

### Backup Creation Failed
```
❌ Failed to create backup for 'template-name'
💡 Check available disk space
💡 Check permissions for backup directory
💡 Use --force to skip backup and proceed
```

### Template in Use
```
⚠️  Template 'active-template' was recently used
💡 Consider backing up before removal
💡 Use --force to proceed anyway
```

## Registry Cleanup

### Space Management
```bash
# Check registry size
du -sh ~/.ason/templates/

# Remove old test templates
ason list --format json | \
  jq -r '.templates[] | select(.name | startswith("test-")) | .name' | \
  xargs -I {} ason remove {} --force

# Remove templates older than 30 days
ason list --format json | \
  jq -r --arg cutoff "$(date -d '30 days ago' -Iseconds)" \
  '.templates[] | select(.added < $cutoff) | .name' | \
  xargs -I {} ason remove {} --backup --force
```

### Selective Cleanup
```bash
# Remove development templates
ason remove dev-react --force
ason remove dev-vue --force
ason remove experimental-api --force

# Remove deprecated versions
ason remove react-v1 --backup
ason remove vue-old --backup
ason remove api-legacy --backup
```

## Backup Management

### Backup Locations
Backups are stored in:
- **Default**: `~/.ason/backups/`
- **Custom**: Specified with `--backup-dir`

### Backup Format
```
~/.ason/backups/
├── react-app-2023-12-01-143022.tar.gz
├── vue-app-2023-11-28-091530.tar.gz
└── go-service-2023-11-25-164512.tar.gz
```

### Backup Naming Convention
```
{template-name}-{YYYY-MM-DD-HHMMSS}.tar.gz
```

### Restoring from Backup
```bash
# Extract and re-add from backup
cd /tmp
tar -xzf ~/.ason/backups/react-app-2023-12-01-143022.tar.gz
ason add react-app ./react-app

# Direct restore (future feature)
ason restore ~/.ason/backups/react-app-2023-12-01-143022.tar.gz
```

## Safety Features

### Confirmation Prompts
Interactive confirmation shows:
- Template name and description
- Template size and file count
- When template was added
- Warning about permanent removal

### Dry Run Analysis
```bash
# Analyze removal impact
ason remove large-template --dry-run
```

### Backup Integration
```bash
# Always backup important templates
ason remove production-template --backup

# Backup with custom location
ason remove shared-template --backup --backup-dir ./team-backups
```

## Common Use Cases

### 1. Development Cleanup
```bash
# Remove test templates
ason remove test-react --force
ason remove test-api --force
ason remove experimental --force

# Clean up old versions
ason remove app-v1 --backup
ason remove service-old --backup
```

### 2. Registry Maintenance
```bash
# Remove unused templates
ason list --format json | \
  jq -r '.templates[] | select(.name | contains("unused")) | .name' | \
  xargs -I {} ason remove {} --force

# Archive old templates
mkdir -p ./archive
for template in $(ason list --format json | jq -r '.templates[] | select(.added < "2023-01-01") | .name'); do
  ason remove "$template" --backup --backup-dir ./archive
done
```

### 3. Template Upgrades
```bash
# Remove old version before adding new
ason remove react-app --backup
ason add react-app ./new-react-template

# Batch upgrade
templates=("react-app" "vue-app" "angular-app")
for template in "${templates[@]}"; do
  ason remove "$template" --backup
  ason add "$template" "./new-templates/$template"
done
```

### 4. Space Management
```bash
# Find largest templates
ason list --sort size --reverse | head -5

# Remove largest unused templates
large_templates=$(ason list --format json | \
  jq -r '.templates | sort_by(.size) | reverse | .[0:3] | .[].name')

for template in $large_templates; do
  echo "Remove large template $template? [y/N]"
  read -r response
  [[ "$response" =~ ^[Yy]$ ]] && ason remove "$template" --backup
done
```

## Integration Examples

### Cleanup Script
```bash
#!/bin/bash
# cleanup-templates.sh

# Configuration
DAYS_OLD=30
BACKUP_DIR="./template-backups"
TEST_PREFIX="test-"

echo "※ Starting template registry cleanup..."

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Remove old test templates
echo "Removing test templates..."
ason list --format json | \
  jq -r ".templates[] | select(.name | startswith(\"$TEST_PREFIX\")) | .name" | \
  while read -r template; do
    echo "Removing test template: $template"
    ason remove "$template" --force
  done

# Archive old templates
echo "Archiving templates older than $DAYS_OLD days..."
cutoff=$(date -d "$DAYS_OLD days ago" -Iseconds)
ason list --format json | \
  jq -r --arg cutoff "$cutoff" \
  '.templates[] | select(.added < $cutoff) | .name' | \
  while read -r template; do
    echo "Archiving old template: $template"
    ason remove "$template" --backup --backup-dir "$BACKUP_DIR"
  done

echo "✅ Cleanup complete!"
echo "📊 Registry status:"
ason list | tail -n 5
```

### CI/CD Integration
```yaml
# .github/workflows/cleanup-registry.yml
name: Cleanup Template Registry
on:
  schedule:
    - cron: '0 2 * * 0'  # Weekly on Sunday at 2 AM

jobs:
  cleanup:
    runs-on: ubuntu-latest
    steps:
      - name: Install Ason
        run: |
          curl -sL https://github.com/madstone-tech/ason/releases/latest/download/ason_Linux_x86_64.tar.gz | tar xz
          sudo mv ason /usr/local/bin/

      - name: Remove test templates
        run: |
          ason list --format json | \
            jq -r '.templates[] | select(.name | startswith("test-")) | .name' | \
            xargs -I {} ason remove {} --force

      - name: Archive old templates
        run: |
          mkdir -p backups
          cutoff=$(date -d '60 days ago' -Iseconds)
          ason list --format json | \
            jq -r --arg cutoff "$cutoff" \
            '.templates[] | select(.added < $cutoff) | .name' | \
            xargs -I {} ason remove {} --backup --backup-dir ./backups
```

## Best Practices

### 1. Always Backup Important Templates
```bash
# Backup before removing production templates
ason remove production-template --backup

# Use descriptive backup locations
ason remove team-template --backup --backup-dir ./team-archives
```

### 2. Use Dry Run for Safety
```bash
# Always preview removal of important templates
ason remove critical-template --dry-run
```

### 3. Batch Operations with Care
```bash
# Preview batch removals
for template in old-*; do
  ason remove "$template" --dry-run
done

# Then execute if safe
for template in old-*; do
  ason remove "$template" --backup --force
done
```

### 4. Regular Maintenance
```bash
# Weekly cleanup routine
./scripts/cleanup-templates.sh

# Monthly archive
./scripts/archive-old-templates.sh
```

## Related Commands

- [`ason list`](list.md) - List available templates
- [`ason add`](add.md) - Add templates to registry
- [`ason validate`](validate.md) - Validate template configuration
- [`ason new`](new.md) - Create projects from templates

## See Also

- [Template Registry Guide](../guides/registry.md)
- [Template Creation Guide](../guides/template-creation.md)
- [Getting Started Guide](../getting-started/quick-start.md)

---

*Release what no longer serves, but preserve the wisdom through careful backup! 🪇*