# UOW-007: Implement File Streaming for Large Template Processing

## Overview
**Phase**: 3 - Performance
**Priority**: Medium
**Estimated Effort**: 8-10 hours
**Dependencies**: UOW-004 (global state removal)

## Problem Description
Current implementation loads entire files into memory for template processing, which causes issues with large templates:
- Memory consumption grows linearly with file size
- Risk of out-of-memory errors with large template files
- Poor performance for templates with large binary assets
- No progress indication for large file operations

```go
// Current memory-intensive approach
content, err := os.ReadFile(srcPath)  // Loads entire file
processedContent, err := g.engine.Render(string(content), context)
```

## Acceptance Criteria
- [ ] Large files (>10MB) are processed using streaming
- [ ] Memory usage remains constant regardless of file size
- [ ] Binary files are streamed without loading into memory
- [ ] Progress indication for large operations
- [ ] Configurable buffer sizes for optimization
- [ ] Backward compatibility with existing templates
- [ ] Performance improvement of >50% for large templates

## Technical Approach

### Implementation Strategy
1. Implement streaming template detection
2. Create buffered file processing pipeline
3. Add progress tracking for large operations
4. Implement configurable buffer management
5. Add memory usage optimization

### Code Changes

**New File**: `internal/generator/streaming.go`
```go
package generator

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
)

const (
    // DefaultBufferSize for streaming operations
    DefaultBufferSize = 64 * 1024 // 64KB

    // LargeFileThreshold defines when to use streaming
    LargeFileThreshold = 10 * 1024 * 1024 // 10MB

    // ProgressReportInterval for large operations
    ProgressReportInterval = 5 * 1024 * 1024 // 5MB
)

// StreamingOptions configures streaming behavior
type StreamingOptions struct {
    BufferSize           int
    LargeFileThreshold   int64
    ProgressCallback     func(processed, total int64)
    EnableProgress       bool
}

// DefaultStreamingOptions returns sensible defaults
func DefaultStreamingOptions() *StreamingOptions {
    return &StreamingOptions{
        BufferSize:         DefaultBufferSize,
        LargeFileThreshold: LargeFileThreshold,
        EnableProgress:     false,
    }
}

// StreamingProcessor handles large file operations
type StreamingProcessor struct {
    options *StreamingOptions
}

// NewStreamingProcessor creates a new streaming processor
func NewStreamingProcessor(options *StreamingOptions) *StreamingProcessor {
    if options == nil {
        options = DefaultStreamingOptions()
    }

    return &StreamingProcessor{
        options: options,
    }
}

// ProcessFile handles both small and large files appropriately
func (sp *StreamingProcessor) ProcessFile(srcPath, destPath string, engine Engine, context map[string]interface{}, mode os.FileMode) error {
    info, err := os.Stat(srcPath)
    if err != nil {
        return fmt.Errorf("failed to stat file %s: %w", srcPath, err)
    }

    // Use streaming for large files
    if info.Size() > sp.options.LargeFileThreshold {
        return sp.streamLargeFile(srcPath, destPath, engine, context, mode, info.Size())
    }

    // Use in-memory processing for small files
    return sp.processSmallFile(srcPath, destPath, engine, context, mode)
}

// streamLargeFile processes large files using streaming
func (sp *StreamingProcessor) streamLargeFile(srcPath, destPath string, engine Engine, context map[string]interface{}, mode os.FileMode, fileSize int64) error {
    // Check if file contains template syntax before processing
    hasTemplates, err := sp.hasTemplateSyntax(srcPath)
    if err != nil {
        return fmt.Errorf("failed to check template syntax: %w", err)
    }

    if !hasTemplates {
        // Stream copy binary files directly
        return sp.streamCopyFile(srcPath, destPath, mode, fileSize)
    }

    // Stream process template files
    return sp.streamProcessTemplate(srcPath, destPath, engine, context, mode, fileSize)
}

// hasTemplateSyntax quickly checks if file contains template syntax
func (sp *StreamingProcessor) hasTemplateSyntax(filePath string) (bool, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return false, err
    }
    defer file.Close()

    // Read first chunk to check for template syntax
    buffer := make([]byte, sp.options.BufferSize)
    scanner := bufio.NewScanner(file)
    scanner.Buffer(buffer, sp.options.BufferSize)

    // Check first few chunks for template markers
    chunkCount := 0
    maxChunks := 10 // Check first 640KB

    for scanner.Scan() && chunkCount < maxChunks {
        line := scanner.Text()
        if strings.Contains(line, "{{") || strings.Contains(line, "{%") {
            return true, nil
        }
        chunkCount++
    }

    return false, scanner.Err()
}

// streamCopyFile copies large binary files without loading into memory
func (sp *StreamingProcessor) streamCopyFile(srcPath, destPath string, mode os.FileMode, fileSize int64) error {
    src, err := os.Open(srcPath)
    if err != nil {
        return fmt.Errorf("failed to open source file: %w", err)
    }
    defer src.Close()

    dest, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer dest.Close()

    // Use buffered copying with progress tracking
    return sp.copyWithProgress(src, dest, fileSize)
}

// streamProcessTemplate processes large template files in chunks
func (sp *StreamingProcessor) streamProcessTemplate(srcPath, destPath string, engine Engine, context map[string]interface{}, mode os.FileMode, fileSize int64) error {
    // For very large templates, we need to be careful about template boundaries
    // For now, fall back to loading into memory with size check
    if fileSize > 100*1024*1024 { // 100MB limit
        return fmt.Errorf("template file %s is too large (%d bytes) for processing", srcPath, fileSize)
    }

    // Process as regular template but with progress tracking
    src, err := os.Open(srcPath)
    if err != nil {
        return fmt.Errorf("failed to open template file: %w", err)
    }
    defer src.Close()

    // Read with progress tracking
    content, err := sp.readWithProgress(src, fileSize)
    if err != nil {
        return fmt.Errorf("failed to read template file: %w", err)
    }

    // Process template
    processedContent, err := engine.Render(string(content), context)
    if err != nil {
        return fmt.Errorf("failed to process template: %w", err)
    }

    // Write with progress tracking
    return sp.writeWithProgress(destPath, []byte(processedContent), mode)
}

// processSmallFile handles small files with in-memory processing
func (sp *StreamingProcessor) processSmallFile(srcPath, destPath string, engine Engine, context map[string]interface{}, mode os.FileMode) error {
    content, err := os.ReadFile(srcPath)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }

    // Check if binary file
    if isBinaryFile(content) {
        return copyFile(srcPath, destPath, mode)
    }

    // Process template
    processedContent, err := engine.Render(string(content), context)
    if err != nil {
        return fmt.Errorf("failed to process template: %w", err)
    }

    return os.WriteFile(destPath, []byte(processedContent), mode)
}

// copyWithProgress copies data with progress reporting
func (sp *StreamingProcessor) copyWithProgress(src io.Reader, dest io.Writer, totalSize int64) error {
    buffer := make([]byte, sp.options.BufferSize)
    var written int64

    for {
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

// readWithProgress reads file content with progress tracking
func (sp *StreamingProcessor) readWithProgress(src io.Reader, totalSize int64) ([]byte, error) {
    content := make([]byte, totalSize)
    var totalRead int64

    for totalRead < totalSize {
        n, err := src.Read(content[totalRead:])
        if n > 0 {
            totalRead += int64(n)

            // Report progress
            if sp.options.EnableProgress && sp.options.ProgressCallback != nil {
                if totalRead%ProgressReportInterval == 0 || totalRead == totalSize {
                    sp.options.ProgressCallback(totalRead, totalSize)
                }
            }
        }

        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("read failed: %w", err)
        }
    }

    return content[:totalRead], nil
}

// writeWithProgress writes content with progress tracking
func (sp *StreamingProcessor) writeWithProgress(destPath string, content []byte, mode os.FileMode) error {
    dest, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer dest.Close()

    totalSize := int64(len(content))
    buffer := content
    bufferSize := sp.options.BufferSize
    var written int64

    for written < totalSize {
        chunkSize := bufferSize
        remaining := totalSize - written
        if remaining < int64(chunkSize) {
            chunkSize = int(remaining)
        }

        n, err := dest.Write(buffer[written : written+int64(chunkSize)])
        if err != nil {
            return fmt.Errorf("write failed: %w", err)
        }

        written += int64(n)

        // Report progress
        if sp.options.EnableProgress && sp.options.ProgressCallback != nil {
            if written%ProgressReportInterval == 0 || written == totalSize {
                sp.options.ProgressCallback(written, totalSize)
            }
        }
    }

    return nil
}

// isBinaryFile checks if content appears to be binary
func isBinaryFile(content []byte) bool {
    // Check for null bytes in first 512 bytes
    checkSize := 512
    if len(content) < checkSize {
        checkSize = len(content)
    }

    for i := 0; i < checkSize; i++ {
        if content[i] == 0 {
            return true
        }
    }

    return false
}
```

**Update**: `internal/generator/generator.go`
```go
// Add streaming support to Generator
type Generator struct {
    engine    Engine
    options   *Options
    streaming *StreamingProcessor
}

// Update NewGenerator to include streaming
func NewGenerator(engine Engine, options *Options) *Generator {
    if options == nil {
        options = defaultOptions()
    }

    // Initialize streaming processor
    streamingOpts := DefaultStreamingOptions()
    if options.Verbose {
        streamingOpts.EnableProgress = true
        streamingOpts.ProgressCallback = func(processed, total int64) {
            percentage := float64(processed) / float64(total) * 100
            fmt.Printf("\r📊 Progress: %.1f%% (%s/%s)", percentage,
                formatBytes(processed), formatBytes(total))
            if processed == total {
                fmt.Println() // New line when complete
            }
        }
    }

    return &Generator{
        engine:    engine,
        options:   options,
        streaming: NewStreamingProcessor(streamingOpts),
    }
}

// Update processFile to use streaming
func (g *Generator) processFile(srcPath, destPath string, context map[string]interface{}, mode os.FileMode, dryRun bool) error {
    if g.options.Verbose {
        fmt.Printf("📄 Processing file: %s\n", filepath.Base(srcPath))
    }

    if dryRun {
        return nil
    }

    return g.streaming.ProcessFile(srcPath, destPath, g.engine, context, mode)
}

// formatBytes formats byte counts for human reading
func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
```

## Configuration Options

**Add to Options struct**:
```go
type Options struct {
    // ... existing fields ...

    // Streaming configuration
    StreamingEnabled     bool
    BufferSize          int
    LargeFileThreshold  int64
    ShowProgress        bool
}
```

## Files to Modify
- `internal/generator/streaming.go` (new file)
- `internal/generator/generator.go` (update to use streaming)
- `internal/generator/streaming_test.go` (new file)
- `cmd/new.go` (add streaming options)

## Testing Strategy
- Unit tests for streaming functions
- Performance benchmarks comparing old vs new approach
- Memory usage tests with large files
- Integration tests with real large templates
- Error handling tests for edge cases

## Performance Targets
- **Memory Usage**: Constant regardless of file size
- **Processing Speed**: >50% improvement for files >10MB
- **Progress Indication**: Updates every 5MB processed
- **Binary File Handling**: Stream copy without processing

## Backward Compatibility
- Default behavior unchanged for small files
- Streaming automatically enabled for large files
- All existing templates work without modification
- CLI options remain the same

## Definition of Done
- Streaming implementation handles files of any size
- Memory usage remains constant during processing
- Performance improvement verified through benchmarks
- Progress indication works for large operations
- All tests pass including performance tests
- Documentation updated with streaming behavior