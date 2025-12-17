# Publishing Ason to Context7

This guide explains how to publish Ason's documentation to [Context7](https://context7.com) - an up-to-date documentation service for LLMs and AI code editors.

## Overview

Context7 allows developers to access current, version-specific documentation directly within their AI coding assistants (Cursor, Claude, etc.) by adding `use context7` to their prompts.

## What's Been Prepared

The Ason repository is now configured for Context7 publishing with:

1. **context7.json** - Configuration file that tells Context7 how to parse and present Ason
2. **Public API Documentation** - Complete documentation in `pkg/` for:
   - Generator API (`pkg/generator.go`)
   - Registry API (`pkg/registry.go`)
   - Engine Interface (`pkg/engine.go`)
3. **Example Documentation** - Usage examples in `examples/` directory
4. **Guides** - Comprehensive guides in `docs/`

## Step-by-Step: Publishing to Context7

### Step 1: Submit to Context7 (One-Time Setup)

Navigate to: **https://context7.com/add-library**

1. Click **"Submit a Library"** or use the **GitHub tab**
2. Paste the Ason repository URL:
   ```
   https://github.com/madstone-tech/ason
   ```
3. Click **Submit**
4. Context7 will:
   - Read the `context7.json` configuration
   - Parse documentation from the configured folders
   - Index all versions listed in `previousVersions`
   - Make Ason available in their library catalog

**Expected Time**: 1-5 minutes for initial parsing

### Step 2: Verify Submission (Check Status)

After submission, you'll see:
- Parsing progress
- Any warnings or errors during indexing
- Preview of how the documentation will appear
- List of indexed files and versions

### Step 3: Claim Ownership (Optional but Recommended)

To manage Ason's configuration on Context7:

1. Go to: **https://context7.com/dashboard**
2. Sign in with GitHub
3. Find Ason in your library list
4. Click **"Claim"** to verify ownership
5. You'll now have access to:
   - Web-based configuration editor
   - Team member management
   - Usage statistics
   - Higher refresh rate limits

## Configuration Details

The `context7.json` file configures how Context7 processes Ason.

### Understanding Ason's Dual Nature

Ason serves **two distinct audiences**, and Context7 is configured to support both:

#### 1. **CLI Users** (DevOps, Template Creators)
```bash
# Using Ason as a command-line tool
ason register my-template ./templates/my-template
ason list
ason new my-template my-project --var name=value
ason validate my-template
```

**Indexed Documentation**:
- `docs/commands/` - Complete CLI reference for all commands
- `docs/getting-started/` - Tutorial guides for CLI usage
- `docs/README.md` - Overview of Ason capabilities

#### 2. **Library Users** (Go Developers, Tool Builders)
```go
// Using Ason as a Go library in your code
gen, _ := pkg.NewGenerator(pkg.NewDefaultEngine())
gen.Generate(ctx, templatePath, variables, outputPath)
```

**Indexed Documentation**:
- `pkg/generator.go` - Generator API with godoc comments
- `pkg/registry.go` - Registry API with godoc comments
- `pkg/engine.go` - Engine interface with godoc comments
- `examples/test_library.go` - Working library usage examples
- `docs/api/engine_interface.md` - Engine interface deep dive

### Included Folders

```json
"folders": ["docs", "pkg", "examples"]
```

| Folder | Purpose | Content |
|--------|---------|---------|
| **docs/** | User-facing documentation | CLI commands, getting started, API guides, configuration |
| **pkg/** | Public Go APIs | Generator, Registry, Engine interface with godoc |
| **examples/** | Working examples | test_library.go (library usage), variable examples |

**What Gets Indexed from `docs/`**:
```
docs/
├── commands/              ← CLI reference (7 command docs)
│   ├── new.md            ← ason new command
│   ├── list.md           ← ason list command
│   ├── register.md       ← ason register command
│   ├── remove.md         ← ason remove command
│   ├── validate.md       ← ason validate command
│   └── completion.md     ← ason completion command
├── getting-started/      ← Learning guides
│   ├── quick-start.md    ← 5-minute tutorial
│   ├── first-template.md ← First template creation
│   ├── installation.md   ← Setup instructions
│   └── configuration.md  ← Configuration guide
├── api/
│   └── engine_interface.md ← Engine API documentation
├── TESTING_GUIDE.md      ← Testing (excluded - internal)
└── README.md             ← Overview
```

**What Gets Indexed from `pkg/`**:
```
pkg/
├── generator.go   ← Generator struct, NewGenerator(), Generate()
├── registry.go    ← Registry struct, operations (Register, List, Remove)
└── engine.go      ← Engine interface, NewDefaultEngine()
```

### Excluded Folders

```json
"excludeFolders": [
  "src", "internal", "cmd", "tests", ".github", "roadmap", "node_modules", ".test", ".git"
]
```

**Why Each Folder is Excluded**:

| Folder | Reason | Impact |
|--------|--------|--------|
| **cmd/** | Implementation source code (CLI command handlers) | Users have CLI docs in `docs/commands/` instead |
| **internal/** | Internal packages (engine, generator impl, registry impl) | Users have public API in `pkg/` instead |
| **tests/** | Test code (not relevant for API users) | Test fixtures and test implementations excluded |
| **src/** | Legacy source directory (if exists) | Avoid redundant indexing |
| **.github/** | CI/CD workflows and templates | Not relevant for users |
| **roadmap/** | Future planning documents | Not part of current API |
| **node_modules/** | NPM dependencies (if any) | Unnecessary dependencies |
| **.test/, .git/** | Build and VCS directories | Auto-cleanup |

**Key Principle**: Exclude *implementation*, keep *user documentation* ✅

The distinction matters:
- ❌ **cmd/new.go** (implementation) → excluded
- ✅ **docs/commands/new.md** (user documentation) → indexed
- ❌ **internal/generator/generator.go** (internal) → excluded
- ✅ **pkg/generator.go** (public API) → indexed

### Excluded Files

```json
"excludeFiles": [
  "CHANGELOG.md", "SECURITY.md", "LICENSE", "CODE_OF_CONDUCT.md", "TESTING_GUIDE.md"
]
```

**Why Each File is Excluded**:

| File | Reason |
|------|--------|
| **CHANGELOG.md** | Version history is handled by `previousVersions` in context7.json |
| **SECURITY.md** | Not relevant for API usage documentation |
| **LICENSE** | Legal document (GitHub already displays) |
| **CODE_OF_CONDUCT.md** | Community guidelines (not API documentation) |
| **TESTING_GUIDE.md** | Internal testing documentation (9000+ words for contributors, not API users) |

**Impact**: Cleaner, more focused documentation without administrative overhead

### Best Practice Rules

```json
"rules": [
  "Always use the public API from pkg/ (Generator, Registry, Engine) for programmatic access",
  "CLI users: Use 'ason new TEMPLATE OUTPUT --var key=value' for project generation",
  "Use NewDefaultEngine() for Pongo2 template engine, or implement custom Engine interface",
  "Registry uses XDG Base Directory specification (~/.local/share/ason/registry.toml)",
  "All paths are validated for path traversal attacks",
  "Template variables must be valid identifiers (alphanumeric + underscore)",
  "Binary files are automatically detected and preserved during generation",
  "Use context.Context for cancellation and timeout support in Generate()",
  "Thread-safe operations: Generator uses RWMutex, Registry supports concurrent reads",
  "Always handle errors returned by Generate() and Registry methods",
  "For template development: Use 'ason validate' CLI command to check template syntax",
  "For large files (>10MB), consider streaming approach for better performance"
]
```

**How Rules Support Both Audiences**:

| Rule | CLI Users | Library Users |
|------|-----------|--------------|
| "CLI users: Use 'ason new TEMPLATE OUTPUT...'" | ✅ Direct guidance | - |
| "Always use the public API from pkg/" | - | ✅ Direct guidance |
| "Use NewDefaultEngine()" | - | ✅ Engine selection |
| "Registry uses XDG Base Directory..." | ✅ File location reference | ✅ Registry behavior |
| "Use context.Context for cancellation..." | - | ✅ Concurrency pattern |
| "Use 'ason validate' CLI command" | ✅ Template validation tool | ✅ Reference for validation |

These rules appear in the documentation provided to coding agents, ensuring both CLI and library users follow best practices.

### Version History
```json
"previousVersions": [
  {"tag": "v0.3.0", "title": "v0.3.0 - Library Export Release"},
  {"tag": "v0.2.2", "title": "v0.2.2"},
  ...
]
```
Users can access documentation for multiple versions - important for users still on v0.2.x.

## What Gets Indexed: Complete Breakdown

### Files Being Indexed (18 total)

**Documentation Files (13)**:
```
docs/api/engine_interface.md
docs/commands/completion.md
docs/commands/list.md
docs/commands/new.md
docs/commands/register.md
docs/commands/remove.md
docs/commands/validate.md
docs/getting-started/configuration.md
docs/getting-started/first-template.md
docs/getting-started/installation.md
docs/getting-started/quick-start.md
docs/README.md
examples/variables/README.md
```

**Source Code Files (5)**:
```
pkg/generator.go          (with godoc comments)
pkg/registry.go           (with godoc comments)
pkg/engine.go             (with godoc comments)
examples/test_library.go  (working example)
examples/variables/simple.json  (example variables)
```

### Files Being Excluded (40+ files)

**Implementation Code** (29 Go files):
- cmd/*.go - CLI command implementations
- internal/**/*.go - Internal packages
- tests/**/*.go - Test code

**Documentation** (5 files):
- CHANGELOG.md - Version history
- SECURITY.md - Security policy
- LICENSE - Legal
- CODE_OF_CONDUCT.md - Community
- TESTING_GUIDE.md - Internal testing

**Other**:
- .github/* - CI/CD workflows
- roadmap/* - Future planning
- Various build artifacts

---

## After Publishing

### Developers Can Use Ason Documentation

Once published, developers can access Ason documentation in their AI editor with both CLI and library examples:

#### Example 1: Library Usage (Go Developer)
**In Cursor or Claude:**
```
Create a Go package that uses the Ason library to scaffold a new TypeScript project. use context7
```

The AI will have access to:
- Generator API signatures with godoc
- Registry API for template management
- Engine interface for custom engines
- Working example from `examples/test_library.go`
- Best practice rules
- Multiple version support

#### Example 2: CLI Usage (DevOps/SRE)
**In Cursor or Claude:**
```
Create a Bash script that uses Ason CLI to generate multiple projects from templates. use context7
```

The AI will have access to:
- All CLI commands (new, list, register, remove, validate, completion)
- Command flags and options
- Example workflows from `docs/commands/`
- Getting started guides
- Best practice rules for CLI usage

#### Example 3: Hybrid Approach
**In Cursor or Claude:**
```
Create a GitHub Actions workflow that both generates projects using the Ason CLI 
and also demonstrates how to use the Ason library in Go. use context7
```

The AI will have access to:
- **CLI section**: How to use `ason new` in shell scripts
- **Library section**: How to use Generator/Registry APIs
- **Both approaches** with complete examples
- Thread-safety considerations
- Error handling patterns

### Monitoring and Updates

#### Check if Library is Published
Search in Context7 chat or use:
```bash
# Context7 CLI (if available)
context7 search ason
```

#### Update Documentation

When you release a new version:

1. **Tag the release** (already done for v0.3.1):
   ```bash
   git tag -a v0.3.2 -m "Release v0.3.2"
   git push origin v0.3.2
   ```

2. **Add to context7.json** if it's important:
   ```json
   "previousVersions": [
     {"tag": "v0.3.2"},
     {"tag": "v0.3.1"},
     ...
   ]
   ```

3. **Trigger refresh** in Context7 dashboard:
   - Go to your library page
   - Click "Refresh" to re-index latest documentation

4. **Monitor parsing**:
   - Check for any errors in indexing
   - Verify new version appears in search results

#### Update Configuration

To change how Context7 parses Ason:

**Option 1: Direct Repository Edit** (Recommended for transparency)
```bash
# Edit context7.json locally
vim context7.json

# Make changes and commit
git add context7.json
git commit -m "docs(context7): update configuration"
git push origin main

# Trigger refresh in Context7 dashboard
```

**Option 2: Context7 Dashboard** (If claimed)
```
1. Go to https://context7.com/dashboard
2. Find Ason library
3. Click "Settings"
4. Edit configuration through UI
5. Changes apply immediately
```

## Troubleshooting

### Library Not Appearing

**Problem**: After submission, Ason doesn't appear in Context7
- **Solution**: Check parsing status in Context7 dashboard for errors
- May take 5-10 minutes for initial indexing
- Trigger manual refresh if available

### Outdated Documentation

**Problem**: Documentation shows old content
- **Solution**: Check the branch configuration in `context7.json`
- Verify branch is set to `"main"`
- Trigger refresh in dashboard

### Version Not Available

**Problem**: User can't access v0.3.0 documentation
- **Solution**: Verify the version tag exists in GitHub:
  ```bash
  git tag | grep v0.3.0  # Should show: v0.3.0
  ```
- Add to `previousVersions` in context7.json if missing
- Trigger refresh to index new version

### Parsing Errors

**Problem**: Context7 reports errors during parsing
- **Solution**: Check `excludeFolders` patterns for correctness
- Remove broken symlinks or circular references
- Ensure documentation is valid markdown
- Review Context7 error messages in dashboard

## Example: What Developers See

When a developer asks with `use context7`:

```
Create a Go program that uses the Ason library to generate a project from a template.
use context7
```

Context7 provides the AI with:
1. **Generator API** - `NewGenerator()`, `Generate()`, `GetEngine()`
2. **Registry API** - `NewRegistry()`, `Register()`, `List()`, `Remove()`
3. **Engine Interface** - `Render()`, `RenderFile()`
4. **Best Practices** - XDG compliance, error handling, threading
5. **Examples** - From `examples/test_library.go`

The AI then generates high-quality code using current APIs without hallucinations.

## Additional Resources

- **Context7 Docs**: https://context7.com/docs/adding-libraries
- **Submit Library**: https://context7.com/add-library
- **Dashboard**: https://context7.com/dashboard
- **GitHub**: https://github.com/upstash/context7

## Notes

- Ason is configured to publish with **main branch** as primary
- Previous versions (v0.2.x, v0.3.0) are available for users on older versions
- Configuration focuses on **public API only** - no internal implementation details
- Best practice rules ensure AI generates correct usage patterns
- All documentation is automatically indexed (no manual step needed)

## Questions?

- Join Context7 Discord: https://upstash.com/discord
- Open issue on Context7 GitHub: https://github.com/upstash/context7/issues
- Contact: https://context7.com/contact
