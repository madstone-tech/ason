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

The `context7.json` file configures how Context7 processes Ason:

### Included Folders
```json
"folders": ["docs", "pkg", "examples"]
```
- **docs/** - User guides and API documentation
- **pkg/** - Public API source code with godoc comments
- **examples/** - Usage examples

### Excluded Folders
```json
"excludeFolders": [
  "src", "internal", "cmd", "tests", ".github", "roadmap", "node_modules"
]
```
- Avoids indexing internal implementation details
- Reduces parsing complexity
- Focuses on user-facing documentation

### Excluded Files
```json
"excludeFiles": ["CHANGELOG.md", "SECURITY.md", "LICENSE", "CODE_OF_CONDUCT.md"]
```
- CHANGELOG is noisy (Context7 has version history)
- Legal/security docs not needed for API reference

### Best Practice Rules
```json
"rules": [
  "Always use the public API from pkg/ (Generator, Registry, Engine)",
  "Use NewDefaultEngine() for Pongo2 engine",
  "Registry uses XDG Base Directory specification",
  ...
]
```
These rules appear in the documentation provided to coding agents, ensuring they follow best practices.

### Version History
```json
"previousVersions": [
  {"tag": "v0.3.0", "title": "v0.3.0 - Library Export Release"},
  {"tag": "v0.2.2", "title": "v0.2.2"},
  ...
]
```
Users can access documentation for multiple versions - important for users still on v0.2.x.

## After Publishing

### Developers Can Use Ason Documentation

Once published, developers can access Ason documentation in their AI editor:

**In Cursor or Claude:**
```
Create a Go package that uses the Ason library to scaffold a new TypeScript project. use context7
```

The AI will have access to:
- Current API signatures
- Usage examples
- Best practices rules
- Multiple version support

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
