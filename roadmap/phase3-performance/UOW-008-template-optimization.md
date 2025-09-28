# UOW-008: Optimize Template Detection with Compiled Regex

## Overview
**Phase**: 3 - Performance
**Priority**: Medium
**Estimated Effort**: 4-6 hours
**Dependencies**: UOW-007 (file streaming)

## Problem Description
Current template detection uses inefficient string operations that impact performance:

```go
// Current inefficient approach in generator.go:206-212
if !strings.Contains(input, "{{") {
    return input, nil
}
```

This approach:
- Performs string search on every file
- Doesn't cache compiled patterns
- Misses complex template syntax patterns
- Repeats work for similar files

## Acceptance Criteria
- [ ] Replace string operations with compiled regex patterns
- [ ] Cache compiled patterns for reuse
- [ ] Support multiple template syntax patterns
- [ ] Improve detection accuracy for edge cases
- [ ] Achieve >30% performance improvement in template detection
- [ ] Maintain backward compatibility
- [ ] Add configurable template syntax support

## Technical Approach

### Implementation Strategy
1. Create template syntax detection system
2. Implement regex pattern compilation and caching
3. Add support for multiple template syntaxes
4. Optimize pattern matching for common cases
5. Add benchmarking and performance validation

### Code Changes

**New File**: `internal/engine/detection.go`
```go
package engine

import (
    "regexp"
    "sync"
)

// TemplateSyntax defines different template syntax patterns
type TemplateSyntax int

const (
    // Pongo2Syntax for Django/Jinja2-style templates
    Pongo2Syntax TemplateSyntax = iota
    // GoTemplateSyntax for Go text/template
    GoTemplateSyntax
    // MustacheSyntax for Mustache templates
    MustacheSyntax
    // CustomSyntax for user-defined patterns
    CustomSyntax
)

// SyntaxPattern defines a template syntax detection pattern
type SyntaxPattern struct {
    Name        string
    VariableRe  *regexp.Regexp
    BlockRe     *regexp.Regexp
    CommentRe   *regexp.Regexp
    QuickCheck  []byte // Fast byte sequence check
}

// TemplateDetector handles efficient template syntax detection
type TemplateDetector struct {
    patterns    map[TemplateSyntax]*SyntaxPattern
    enabledSyntax []TemplateSyntax
    mu          sync.RWMutex
}

// NewTemplateDetector creates a new template detector with default patterns
func NewTemplateDetector() *TemplateDetector {
    detector := &TemplateDetector{
        patterns: make(map[TemplateSyntax]*SyntaxPattern),
        enabledSyntax: []TemplateSyntax{Pongo2Syntax}, // Default to Pongo2
    }

    detector.initializeDefaultPatterns()
    return detector
}

// initializeDefaultPatterns sets up built-in template syntax patterns
func (td *TemplateDetector) initializeDefaultPatterns() {
    // Pongo2/Django/Jinja2 syntax
    td.patterns[Pongo2Syntax] = &SyntaxPattern{
        Name:       "Pongo2",
        VariableRe: regexp.MustCompile(`\{\{\s*[^}]+\s*\}\}`),
        BlockRe:    regexp.MustCompile(`\{%\s*[^%]+\s*%\}`),
        CommentRe:  regexp.MustCompile(`\{#\s*[^#]*\s*#\}`),
        QuickCheck: []byte("{{"),
    }

    // Go template syntax
    td.patterns[GoTemplateSyntax] = &SyntaxPattern{
        Name:       "GoTemplate",
        VariableRe: regexp.MustCompile(`\{\{\s*[^}]+\s*\}\}`),
        BlockRe:    regexp.MustCompile(`\{\{\s*(if|range|with|define|template|block)\s+[^}]+\s*\}\}`),
        CommentRe:  regexp.MustCompile(`\{\{/\*.*?\*/\}\}`),
        QuickCheck: []byte("{{"),
    }

    // Mustache syntax
    td.patterns[MustacheSyntax] = &SyntaxPattern{
        Name:       "Mustache",
        VariableRe: regexp.MustCompile(`\{\{\s*[^}]+\s*\}\}`),
        BlockRe:    regexp.MustCompile(`\{\{[#^/]\s*[^}]+\s*\}\}`),
        CommentRe:  regexp.MustCompile(`\{\{!\s*[^}]*\s*\}\}`),
        QuickCheck: []byte("{{"),
    }
}

// SetEnabledSyntax configures which template syntaxes to detect
func (td *TemplateDetector) SetEnabledSyntax(syntaxes []TemplateSyntax) {
    td.mu.Lock()
    defer td.mu.Unlock()
    td.enabledSyntax = make([]TemplateSyntax, len(syntaxes))
    copy(td.enabledSyntax, syntaxes)
}

// AddCustomSyntax adds a user-defined template syntax pattern
func (td *TemplateDetector) AddCustomSyntax(name string, variablePattern, blockPattern, commentPattern string) error {
    td.mu.Lock()
    defer td.mu.Unlock()

    variableRe, err := regexp.Compile(variablePattern)
    if err != nil {
        return fmt.Errorf("invalid variable pattern: %w", err)
    }

    blockRe, err := regexp.Compile(blockPattern)
    if err != nil {
        return fmt.Errorf("invalid block pattern: %w", err)
    }

    commentRe, err := regexp.Compile(commentPattern)
    if err != nil {
        return fmt.Errorf("invalid comment pattern: %w", err)
    }

    td.patterns[CustomSyntax] = &SyntaxPattern{
        Name:       name,
        VariableRe: variableRe,
        BlockRe:    blockRe,
        CommentRe:  commentRe,
        QuickCheck: []byte(variablePattern[:2]), // Use first 2 chars for quick check
    }

    return nil
}

// HasTemplateContent efficiently detects if content contains template syntax
func (td *TemplateDetector) HasTemplateContent(content []byte) bool {
    td.mu.RLock()
    defer td.mu.RUnlock()

    // Quick check for common cases
    for _, syntax := range td.enabledSyntax {
        pattern, exists := td.patterns[syntax]
        if !exists {
            continue
        }

        // Fast byte sequence check first
        if len(pattern.QuickCheck) > 0 && !containsBytes(content, pattern.QuickCheck) {
            continue
        }

        // More thorough regex check
        if pattern.VariableRe.Match(content) ||
           pattern.BlockRe.Match(content) ||
           pattern.CommentRe.Match(content) {
            return true
        }
    }

    return false
}

// DetectSyntaxType identifies which template syntax is used in content
func (td *TemplateDetector) DetectSyntaxType(content []byte) (TemplateSyntax, bool) {
    td.mu.RLock()
    defer td.mu.RUnlock()

    for _, syntax := range td.enabledSyntax {
        pattern, exists := td.patterns[syntax]
        if !exists {
            continue
        }

        // Quick check first
        if len(pattern.QuickCheck) > 0 && !containsBytes(content, pattern.QuickCheck) {
            continue
        }

        // Check for syntax-specific patterns
        if pattern.VariableRe.Match(content) ||
           pattern.BlockRe.Match(content) ||
           pattern.CommentRe.Match(content) {
            return syntax, true
        }
    }

    return Pongo2Syntax, false // Default fallback
}

// AnalyzeTemplate provides detailed analysis of template content
func (td *TemplateDetector) AnalyzeTemplate(content []byte) *TemplateAnalysis {
    td.mu.RLock()
    defer td.mu.RUnlock()

    analysis := &TemplateAnalysis{
        HasTemplates: false,
        SyntaxTypes:  make(map[TemplateSyntax]int),
        Variables:    make([]string, 0),
        Blocks:       make([]string, 0),
    }

    for _, syntax := range td.enabledSyntax {
        pattern, exists := td.patterns[syntax]
        if !exists {
            continue
        }

        // Count variables
        varMatches := pattern.VariableRe.FindAllString(string(content), -1)
        if len(varMatches) > 0 {
            analysis.HasTemplates = true
            analysis.SyntaxTypes[syntax] += len(varMatches)
            analysis.Variables = append(analysis.Variables, varMatches...)
        }

        // Count blocks
        blockMatches := pattern.BlockRe.FindAllString(string(content), -1)
        if len(blockMatches) > 0 {
            analysis.HasTemplates = true
            analysis.SyntaxTypes[syntax] += len(blockMatches)
            analysis.Blocks = append(analysis.Blocks, blockMatches...)
        }
    }

    return analysis
}

// TemplateAnalysis provides detailed information about template content
type TemplateAnalysis struct {
    HasTemplates bool
    SyntaxTypes  map[TemplateSyntax]int
    Variables    []string
    Blocks       []string
}

// containsBytes efficiently checks if haystack contains needle
func containsBytes(haystack, needle []byte) bool {
    if len(needle) == 0 {
        return true
    }
    if len(haystack) < len(needle) {
        return false
    }

    for i := 0; i <= len(haystack)-len(needle); i++ {
        if haystack[i] == needle[0] {
            match := true
            for j := 1; j < len(needle); j++ {
                if haystack[i+j] != needle[j] {
                    match = false
                    break
                }
            }
            if match {
                return true
            }
        }
    }
    return false
}
```

**Update**: `internal/generator/generator.go`
```go
// Add template detector to Generator
type Generator struct {
    engine    Engine
    options   *Options
    streaming *StreamingProcessor
    detector  *engine.TemplateDetector
}

// Update NewGenerator to include detector
func NewGenerator(engine Engine, options *Options) *Generator {
    if options == nil {
        options = defaultOptions()
    }

    // Initialize template detector
    detector := engine.NewTemplateDetector()
    if options.TemplateSyntax != nil {
        detector.SetEnabledSyntax(options.TemplateSyntax)
    }

    // Configure streaming with detector
    streamingOpts := DefaultStreamingOptions()
    streamingOpts.TemplateDetector = detector

    // ... rest of initialization
}

// Update template processing to use optimized detection
func (g *Generator) processTemplateContent(content []byte, context map[string]interface{}) ([]byte, error) {
    // Fast template detection
    if !g.detector.HasTemplateContent(content) {
        return content, nil // No templates found, return as-is
    }

    // Process with template engine
    processedContent, err := g.engine.Render(string(content), context)
    if err != nil {
        return nil, fmt.Errorf("template rendering failed: %w", err)
    }

    return []byte(processedContent), nil
}
```

**Update**: `internal/generator/streaming.go`
```go
// Update StreamingOptions to include detector
type StreamingOptions struct {
    BufferSize           int
    LargeFileThreshold   int64
    ProgressCallback     func(processed, total int64)
    EnableProgress       bool
    TemplateDetector     *engine.TemplateDetector
}

// Update hasTemplateSyntax to use optimized detection
func (sp *StreamingProcessor) hasTemplateSyntax(filePath string) (bool, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return false, err
    }
    defer file.Close()

    // Read a reasonable chunk for detection
    buffer := make([]byte, sp.options.BufferSize)
    n, err := file.Read(buffer)
    if err != nil && err != io.EOF {
        return false, err
    }

    // Use optimized detector
    if sp.options.TemplateDetector != nil {
        return sp.options.TemplateDetector.HasTemplateContent(buffer[:n]), nil
    }

    // Fallback to simple detection
    return containsBytes(buffer[:n], []byte("{{")), nil
}
```

**New File**: `internal/engine/detection_test.go`
```go
package engine

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestTemplateDetector_HasTemplateContent(t *testing.T) {
    detector := NewTemplateDetector()

    tests := []struct {
        name     string
        content  string
        expected bool
    }{
        {
            name:     "simple variable",
            content:  "Hello {{name}}!",
            expected: true,
        },
        {
            name:     "block syntax",
            content:  "{% if user %}Welcome{% endif %}",
            expected: true,
        },
        {
            name:     "comment syntax",
            content:  "{# This is a comment #}",
            expected: true,
        },
        {
            name:     "no templates",
            content:  "Plain text content",
            expected: false,
        },
        {
            name:     "false positive curly braces",
            content:  "function() { return {}; }",
            expected: false,
        },
        {
            name:     "mixed content",
            content:  "Some text {{variable}} more text",
            expected: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := detector.HasTemplateContent([]byte(tt.content))
            assert.Equal(t, tt.expected, result)
        })
    }
}

func BenchmarkTemplateDetection(b *testing.B) {
    detector := NewTemplateDetector()

    testCases := []struct {
        name    string
        content []byte
    }{
        {
            name:    "no_templates",
            content: []byte("This is plain text content without any template syntax"),
        },
        {
            name:    "simple_template",
            content: []byte("Hello {{name}}, welcome to {{site}}!"),
        },
        {
            name:    "complex_template",
            content: []byte(`
                {% for item in items %}
                    <div>{{item.name}} - {{item.value}}</div>
                {% endfor %}
                {# Comment here #}
            `),
        },
    }

    for _, tc := range testCases {
        b.Run("optimized_"+tc.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                detector.HasTemplateContent(tc.content)
            }
        })

        b.Run("string_contains_"+tc.name, func(b *testing.B) {
            content := string(tc.content)
            for i := 0; i < b.N; i++ {
                _ = strings.Contains(content, "{{")
            }
        })
    }
}
```

## Configuration Options

**Add to Options struct**:
```go
type Options struct {
    // ... existing fields ...

    // Template detection configuration
    TemplateSyntax      []engine.TemplateSyntax
    CustomSyntaxPattern map[string]string
}
```

## Files to Modify
- `internal/engine/detection.go` (new file)
- `internal/engine/detection_test.go` (new file)
- `internal/generator/generator.go` (update to use detector)
- `internal/generator/streaming.go` (update detection logic)
- `cmd/new.go` (add syntax configuration options)

## Performance Targets
- **Detection Speed**: >30% improvement over string operations
- **Memory Usage**: Minimal overhead from compiled patterns
- **Accuracy**: 100% detection of valid template syntax
- **Caching**: Compiled patterns reused across files

## Testing Strategy
- Unit tests for all syntax patterns
- Benchmark tests comparing old vs new detection
- Performance regression tests
- Edge case testing (malformed syntax, mixed content)
- Cross-platform compatibility tests

## Backward Compatibility
- Default behavior identical to current implementation
- Existing templates work without modification
- CLI interface unchanged
- Optional advanced features for power users

## Definition of Done
- Regex-based detection replaces string operations
- Pattern compilation and caching implemented
- Performance improvement >30% verified through benchmarks
- All syntax patterns properly detected
- Tests pass with comprehensive coverage
- Documentation updated with new detection features