# UOW-012: Add Context Support for Operation Timeouts and Cancellation

## Overview
**Phase**: 4 - Architecture
**Priority**: Medium
**Estimated Effort**: 6-8 hours
**Dependencies**: UOW-009 (concurrent processing), UOW-011 (config management)

## Problem Description
Current implementation lacks context support for operation control:
- No way to cancel long-running operations
- Missing timeout handling for template processing
- No graceful shutdown mechanism
- Resource leaks during interrupted operations
- Poor user experience with unresponsive operations

Functions don't accept `context.Context` parameters, making it impossible to implement proper cancellation and timeout handling.

## Acceptance Criteria
- [ ] All major operations accept `context.Context` parameters
- [ ] Timeout handling for template processing operations
- [ ] Graceful cancellation with proper cleanup
- [ ] Resource management with context-aware goroutines
- [ ] Progress indication that respects cancellation
- [ ] Configurable timeouts through configuration system
- [ ] Signal handling for CLI operations
- [ ] Backward compatibility maintained

## Technical Approach

### Implementation Strategy
1. Add context parameters to all major functions
2. Implement timeout and cancellation handling
3. Create context-aware resource management
4. Add signal handling for CLI operations
5. Integrate with configuration system for timeouts

### Code Changes

**Update**: `internal/generator/generator.go`
```go
package generator

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"
)

// Generate processes a template with context support for cancellation and timeouts
func (g *Generator) Generate(ctx context.Context, templatePath, outputPath string, context map[string]interface{}) error {
    // Create a derived context with timeout if configured
    if g.options.ProcessingTimeout > 0 {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, g.options.ProcessingTimeout)
        defer cancel()
    }

    if g.options.Verbose {
        fmt.Printf("🚀 Generating project from template: %s\n", templatePath)
    }

    // Check for early cancellation
    select {
    case <-ctx.Done():
        return fmt.Errorf("generation cancelled: %w", ctx.Err())
    default:
    }

    // Collect all files to process
    jobs, err := g.collectProcessingJobsWithContext(ctx, templatePath, outputPath, context)
    if err != nil {
        return fmt.Errorf("failed to collect processing jobs: %w", err)
    }

    if len(jobs) == 0 {
        fmt.Println("ℹ️  No files to process")
        return nil
    }

    // Use concurrent processing if enabled and beneficial
    if g.options.ConcurrentProcessing && len(jobs) > g.options.ConcurrentThreshold {
        return g.concurrent.ProcessFiles(ctx, jobs)
    }

    // Fall back to sequential processing with context
    return g.processSequentiallyWithContext(ctx, jobs)
}

// collectProcessingJobsWithContext walks the template with context awareness
func (g *Generator) collectProcessingJobsWithContext(ctx context.Context, templatePath, outputPath string, context map[string]interface{}) ([]*ProcessingJob, error) {
    var jobs []*ProcessingJob
    jobID := 0

    // Create a channel to signal walking completion
    done := make(chan error, 1)

    go func() {
        err := filepath.Walk(templatePath, func(srcPath string, info os.FileInfo, err error) error {
            // Check for cancellation during walk
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
            }

            if err != nil {
                return err
            }

            if info.IsDir() {
                return nil // Directories are created on-demand
            }

            srcRelPath, err := filepath.Rel(templatePath, srcPath)
            if err != nil {
                return fmt.Errorf("failed to get relative path: %w", err)
            }

            destRelPath := g.processPath(srcRelPath, context)
            destPath, err := validatePath(outputPath, destRelPath)
            if err != nil {
                return fmt.Errorf("path validation failed: %w", err)
            }

            // Create processing job
            job := &ProcessingJob{
                ID:       fmt.Sprintf("job_%d", jobID),
                SrcPath:  srcPath,
                DestPath: destPath,
                Context:  context,
                Mode:     info.Mode(),
                Size:     info.Size(),
                Priority: g.determineJobPriority(srcPath, info.Size()),
            }

            jobs = append(jobs, job)
            jobID++

            return nil
        })
        done <- err
    }()

    // Wait for completion or cancellation
    select {
    case err := <-done:
        if err != nil {
            return nil, err
        }
    case <-ctx.Done():
        return nil, fmt.Errorf("file collection cancelled: %w", ctx.Err())
    }

    return jobs, nil
}

// processSequentiallyWithContext handles files one by one with context support
func (g *Generator) processSequentiallyWithContext(ctx context.Context, jobs []*ProcessingJob) error {
    for i, job := range jobs {
        // Check for cancellation before processing each file
        select {
        case <-ctx.Done():
            return fmt.Errorf("processing cancelled after %d/%d files: %w", i, len(jobs), ctx.Err())
        default:
        }

        if g.options.Verbose {
            fmt.Printf("📄 Processing file %d/%d: %s\n", i+1, len(jobs), filepath.Base(job.SrcPath))
        }

        // Create a context with timeout for individual file processing
        fileCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
        err := g.streaming.ProcessFileWithContext(fileCtx, job.SrcPath, job.DestPath, g.engine, job.Context, job.Mode)
        cancel()

        if err != nil {
            if fileCtx.Err() == context.DeadlineExceeded {
                return fmt.Errorf("file processing timeout for %s: %w", job.SrcPath, err)
            }
            return fmt.Errorf("failed to process %s: %w", job.SrcPath, err)
        }
    }

    return nil
}
```

**Update**: `internal/generator/streaming.go`
```go
// ProcessFileWithContext handles file processing with context support
func (sp *StreamingProcessor) ProcessFileWithContext(ctx context.Context, srcPath, destPath string, engine Engine, context map[string]interface{}, mode os.FileMode) error {
    // Check for cancellation before starting
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    info, err := os.Stat(srcPath)
    if err != nil {
        return fmt.Errorf("failed to stat file %s: %w", srcPath, err)
    }

    // Use streaming for large files
    if info.Size() > sp.options.LargeFileThreshold {
        return sp.streamLargeFileWithContext(ctx, srcPath, destPath, engine, context, mode, info.Size())
    }

    // Use in-memory processing for small files
    return sp.processSmallFileWithContext(ctx, srcPath, destPath, engine, context, mode)
}

// streamLargeFileWithContext processes large files with context awareness
func (sp *StreamingProcessor) streamLargeFileWithContext(ctx context.Context, srcPath, destPath string, engine Engine, context map[string]interface{}, mode os.FileMode, fileSize int64) error {
    // Check if file contains template syntax before processing
    hasTemplates, err := sp.hasTemplateSyntaxWithContext(ctx, srcPath)
    if err != nil {
        return fmt.Errorf("failed to check template syntax: %w", err)
    }

    if !hasTemplates {
        // Stream copy binary files directly
        return sp.streamCopyFileWithContext(ctx, srcPath, destPath, mode, fileSize)
    }

    // Stream process template files
    return sp.streamProcessTemplateWithContext(ctx, srcPath, destPath, engine, context, mode, fileSize)
}

// copyWithProgressAndContext copies data with progress reporting and context awareness
func (sp *StreamingProcessor) copyWithProgressAndContext(ctx context.Context, src io.Reader, dest io.Writer, totalSize int64) error {
    buffer := make([]byte, sp.options.BufferSize)
    var written int64

    for {
        // Check for cancellation
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        n, err := src.Read(buffer)
        if n > 0 {
            if _, writeErr := dest.Write(buffer[:n]); writeErr != nil {
                return fmt.Errorf("write failed: %w", writeErr)
            }
            written += int64(n)

            // Report progress
            if sp.options.EnableProgress && sp.options.ProgressCallback != nil {
                if written%ProgressReportInterval == 0 || written == totalSize {
                    sp.options.ProgressCallback(written, totalSize)
                }
            }
        }

        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("read failed: %w", err)
        }
    }

    return nil
}
```

**New File**: `internal/context/context.go`
```go
package context

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// WithSignalHandling creates a context that cancels on common signals
func WithSignalHandling(parent context.Context) (context.Context, context.CancelFunc) {
    ctx, cancel := context.WithCancel(parent)

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

    go func() {
        select {
        case sig := <-sigChan:
            fmt.Printf("\n🛑 Received signal %v, shutting down gracefully...\n", sig)
            cancel()
        case <-ctx.Done():
        }
        signal.Stop(sigChan)
    }()

    return ctx, cancel
}

// WithTimeout creates a context with configurable timeout
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
    if timeout <= 0 {
        // Return a context that never times out
        return context.WithCancel(parent)
    }
    return context.WithTimeout(parent, timeout)
}

// WithDeadline creates a context with a specific deadline
func WithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
    return context.WithDeadline(parent, deadline)
}

// BackgroundWithTimeout creates a background context with timeout
func BackgroundWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
    return WithTimeout(context.Background(), timeout)
}

// TODO creates a context for TODO operations
func TODO() context.Context {
    return context.TODO()
}

// WithValue creates a context with a key-value pair
func WithValue(parent context.Context, key, value interface{}) context.Context {
    return context.WithValue(parent, key, value)
}

// OperationType represents different operation types for context keys
type OperationType string

const (
    // Generation represents project generation operations
    Generation OperationType = "generation"
    // TemplateProcessing represents template processing operations
    TemplateProcessing OperationType = "template_processing"
    // PluginLoading represents plugin loading operations
    PluginLoading OperationType = "plugin_loading"
    // RegistryOperation represents registry operations
    RegistryOperation OperationType = "registry_operation"
)

// WithOperation adds operation type to context
func WithOperation(parent context.Context, operation OperationType) context.Context {
    return context.WithValue(parent, "operation", operation)
}

// GetOperation retrieves operation type from context
func GetOperation(ctx context.Context) (OperationType, bool) {
    op, ok := ctx.Value("operation").(OperationType)
    return op, ok
}

// ProgressTracker represents progress tracking in context
type ProgressTracker struct {
    Current int64
    Total   int64
    Message string
}

// WithProgress adds progress tracking to context
func WithProgress(parent context.Context, tracker *ProgressTracker) context.Context {
    return context.WithValue(parent, "progress", tracker)
}

// GetProgress retrieves progress tracker from context
func GetProgress(ctx context.Context) (*ProgressTracker, bool) {
    tracker, ok := ctx.Value("progress").(*ProgressTracker)
    return tracker, ok
}

// UpdateProgress updates progress in context if available
func UpdateProgress(ctx context.Context, current, total int64, message string) {
    if tracker, ok := GetProgress(ctx); ok {
        tracker.Current = current
        tracker.Total = total
        tracker.Message = message
    }
}
```

**Update**: `cmd/new.go`
```go
// Update new command to use context
func (o *newOptions) run() error {
    // Create context with signal handling
    ctx, cancel := context.WithSignalHandling(context.Background())
    defer cancel()

    // Add operation type to context
    ctx = context.WithOperation(ctx, context.Generation)

    // Add timeout if configured
    if o.timeout > 0 {
        var timeoutCancel context.CancelFunc
        ctx, timeoutCancel = context.WithTimeout(ctx, o.timeout)
        defer timeoutCancel()
    }

    // Create generator options from CLI options
    genOptions := &generator.Options{
        Verbose:            o.verbose,
        DryRun:            o.dryRun,
        Variables:         o.variables,
        ProcessingTimeout: o.processingTimeout,
    }

    // Create template engine
    engine := engine.NewPongo2Engine()

    // Create generator with options
    gen := generator.NewGenerator(engine, genOptions)

    // Get template path
    templatePath, err := o.getTemplatePath()
    if err != nil {
        return err
    }

    // Prepare context variables
    contextVars := make(map[string]interface{})
    for key, value := range o.variables {
        contextVars[key] = value
    }

    // Generate project with context
    return gen.Generate(ctx, templatePath, o.output, contextVars)
}
```

**Update**: `internal/registry/registry.go`
```go
// Add context support to registry operations
func (r *Registry) AddWithContext(ctx context.Context, name, sourcePath string) error {
    // Check for cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Existing implementation with context checks...
    return r.Add(name, sourcePath)
}

func (r *Registry) ListWithContext(ctx context.Context) ([]Template, error) {
    // Check for cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    return r.List()
}
```

## Configuration Integration

**Add to config Options**:
```go
type GeneratorConfig struct {
    // ... existing fields ...

    // Timeout configuration
    ProcessingTimeout   time.Duration `mapstructure:"processing_timeout"`
    FileTimeout        time.Duration `mapstructure:"file_timeout"`
    OperationTimeout   time.Duration `mapstructure:"operation_timeout"`

    // Cancellation configuration
    GracefulShutdown   bool          `mapstructure:"graceful_shutdown"`
    ShutdownTimeout    time.Duration `mapstructure:"shutdown_timeout"`
}
```

## Files to Modify
- `internal/generator/generator.go` (add context parameters)
- `internal/generator/streaming.go` (add context support)
- `internal/generator/concurrent.go` (update worker context handling)
- `internal/context/context.go` (new file)
- `cmd/new.go` (add signal handling and timeouts)
- `cmd/root.go` (add global timeout flags)
- `internal/registry/registry.go` (add context to operations)
- `internal/plugin/plugin.go` (add context to plugin operations)

## CLI Integration

**Add global flags**:
```bash
ason new template output --timeout 5m --processing-timeout 30s
```

**Signal handling**:
- `Ctrl+C` (SIGINT): Graceful cancellation
- `SIGTERM`: Graceful shutdown
- `SIGQUIT`: Force quit with cleanup

## Testing Strategy
- Unit tests for context cancellation scenarios
- Integration tests with timeouts
- Signal handling tests
- Resource cleanup verification tests
- Performance impact measurement

## Error Handling
- Clear error messages for timeout scenarios
- Proper resource cleanup on cancellation
- Graceful degradation when operations are interrupted
- Context error propagation throughout the call stack

## Definition of Done
- All major operations accept and respect context parameters
- Timeout handling works for all configured operations
- Signal handling provides graceful cancellation
- Resource cleanup is properly handled on interruption
- Configuration system supports timeout settings
- Error messages clearly indicate timeout/cancellation reasons
- Performance impact is minimal
- All tests pass including cancellation scenarios