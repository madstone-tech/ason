# Variable Files

Ason supports loading template variables from external files instead of passing them all via command-line flags. This is especially useful when you have many variables or want to maintain different configuration sets for different environments.

## Supported Formats

- **TOML** (`.toml`) - Recommended, matches `ason.toml` template format
- **YAML** (`.yaml`, `.yml`) - Common in cloud/DevOps workflows
- **JSON** (`.json`) - Universal format

## Usage

### Basic Usage

```bash
# Load all variables from a file
ason new my-template ./output --var-file prod.toml
```

### Override File Variables with CLI

CLI variables take precedence over file variables:

```bash
# Use base.toml but override environment
ason new my-template ./output --var-file base.toml --var environment=prod
```

### Short Flag

You can use the short flag `-f`:

```bash
ason new my-template ./output -f prod.toml
```

## File Formats

### Simple TOML (simple.toml)

```toml
environment = "production"
aws_region = "us-west-2"
organization = "acme"
project_name = "my-app"
```

### Template-Style TOML (template-style.toml)

This format is compatible with `ason.toml` template configuration and includes default values:

```toml
[variables]
organization = { type = "string", description = "Organization name", default = "acme" }
aws_region = { type = "string", description = "AWS region", default = "us-east-1" }
environment = { type = "string", description = "Environment", default = "dev" }
```

### YAML (simple.yaml)

```yaml
environment: production
aws_region: us-west-2
organization: acme
project_name: my-app
```

### JSON (simple.json)

```json
{
  "environment": "production",
  "aws_region": "us-west-2",
  "organization": "acme",
  "project_name": "my-app"
}
```

## Real-World Example

For a complex AWS Lambda deployment with many variables:

**Before** (with individual --var flags):
```bash
ason new lambda-waf-ipset ./deployments/prod \
  --var environment=prod \
  --var aws_region=us-west-2 \
  --var organization=acme \
  --var project_name=waf-updater \
  --var lambda_name=lambda_waf_ipset_updater \
  --var lambda_description="Lambda to update WAFv2 IPSET" \
  --var wafv_ipset_arn=arn:aws:wafv2:us-west-2:123:regional/ipset/prod/xyz
```

**After** (with --var-file):
```bash
ason new lambda-waf-ipset ./deployments/prod --var-file prod.toml
```

Where `prod.toml` contains:
```toml
environment = "prod"
aws_region = "us-west-2"
organization = "acme"
project_name = "waf-updater"
lambda_name = "lambda_waf_ipset_updater"
lambda_description = "Lambda to update WAFv2 IPSET"
wafv_ipset_arn = "arn:aws:wafv2:us-west-2:123:regional/ipset/prod/xyz"
```

## Environment-Specific Configurations

Create separate variable files for each environment:

```bash
# Development
ason new my-service ./dev --var-file environments/dev.toml

# Staging
ason new my-service ./staging --var-file environments/staging.toml

# Production
ason new my-service ./prod --var-file environments/prod.toml
```

## Best Practices

1. **Version Control**: Commit variable files to git for team collaboration
2. **Secrets Management**: Don't store sensitive values in variable files; use CLI overrides or environment variables
3. **Documentation**: Add comments in TOML/YAML files to explain variable purposes
4. **Naming**: Use descriptive filenames like `prod.toml`, `staging.yaml`, `dev.toml`
5. **Validation**: Use `--dry-run` flag to preview output before generating

## Variable Precedence

When variables are defined in multiple places, this is the precedence (highest to lowest):

1. CLI flags (`--var key=value`)
2. Variable file (`--var-file file.toml`)
3. Template defaults (in `ason.toml`)

## Tips

- Use TOML for most cases - it's readable and supports comments
- Use YAML when integrating with existing YAML-based tooling
- Use JSON when generating variable files programmatically
- Combine `--var-file` with selective `--var` overrides for maximum flexibility
