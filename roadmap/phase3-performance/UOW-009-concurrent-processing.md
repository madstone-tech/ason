# UOW-009: Add Worker Pools for Concurrent File Processing

## Overview
**Phase**: 3 - Performance
**Priority**: Medium
**Estimated Effort**: 8-10 hours
**Dependencies**: UOW-007 (file streaming), UOW-008 (template optimization)

## Problem Description
Current implementation processes files sequentially, which is inefficient for large templates with many files:
- No parallelization of file operations
- CPU cores underutilized during I/O operations
- Long processing times for large templates
- No control over concurrency levels

Sequential processing in `walkTemplateFiles` function limits performance for templates with hundreds of files.

## Acceptance Criteria
- [ ] Concurrent file processing with configurable worker pools
- [ ] Intelligent work distribution based on file size and type
- [ ] Dependency-aware processing order
- [ ] Progress tracking for concurrent operations
- [ ] Resource limits to prevent system overload
- [ ] Error handling that doesn't fail entire generation
- [ ] Performance improvement of >60% for large templates

## Technical Approach

### Implementation Strategy
1. Create worker pool system for file processing
2. Implement dependency-aware job scheduling
3. Add concurrent-safe progress tracking
4. Create resource management and throttling
5. Implement error aggregation and recovery

### Code Changes

**New File**: `internal/generator/concurrent.go`
```go
package generator

import (
    "context"
    "fmt"
    "path/filepath"
    "runtime"
    "sync"
    "time"
)

// ConcurrentProcessor handles parallel file processing
type ConcurrentProcessor struct {
    workerCount     int
    maxMemoryMB     int
    generator       *Generator
    progressTracker *ProgressTracker
}

// ProcessingJob represents a file processing task
type ProcessingJob struct {
    ID          string
    SrcPath     string
    DestPath    string
    Context     map[string]interface{}
    Mode        os.FileMode
    Size        int64
    Priority    JobPriority
    Dependencies []string // Job IDs this job depends on
}

// JobPriority defines processing priority
type JobPriority int

const (
    // LowPriority for large binary files
    LowPriority JobPriority = iota
    // NormalPriority for regular template files
    NormalPriority
    // HighPriority for small configuration files
    HighPriority
)

// ProcessingResult contains the outcome of a job
type ProcessingResult struct {
    JobID   string
    Success bool
    Error   error
    Duration time.Duration
}

// NewConcurrentProcessor creates a new concurrent processor
func NewConcurrentProcessor(generator *Generator, options *ConcurrentOptions) *ConcurrentProcessor {
    if options == nil {
        options = DefaultConcurrentOptions()
    }

    return &ConcurrentProcessor{
        workerCount:     options.WorkerCount,
        maxMemoryMB:     options.MaxMemoryMB,
        generator:       generator,
        progressTracker: NewProgressTracker(),
    }
}

// ConcurrentOptions configures concurrent processing
type ConcurrentOptions struct {
    WorkerCount      int  // Number of worker goroutines
    MaxMemoryMB      int  // Maximum memory usage in MB
    EnableDependency bool // Enable dependency-aware processing
    BufferSize       int  // Channel buffer size
}

// DefaultConcurrentOptions returns sensible defaults
func DefaultConcurrentOptions() *ConcurrentOptions {
    return &ConcurrentOptions{
        WorkerCount:      runtime.NumCPU(),
        MaxMemoryMB:      1024, // 1GB
        EnableDependency: true,
        BufferSize:       100,
    }
}

// ProcessFiles processes multiple files concurrently
func (cp *ConcurrentProcessor) ProcessFiles(ctx context.Context, jobs []*ProcessingJob) error {
    if len(jobs) == 0 {
        return nil
    }

    // Initialize progress tracking
    cp.progressTracker.SetTotal(len(jobs))

    // Create job scheduler
    scheduler := NewJobScheduler(jobs, cp.workerCount)

    // Create worker pool
    workers := cp.createWorkerPool(ctx, scheduler)

    // Start workers
    var wg sync.WaitGroup
    results := make(chan *ProcessingResult, len(jobs))

    for i := 0; i < cp.workerCount; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            cp.worker(ctx, workerID, workers, results)
        }(i)
    }

    // Close job channel when all jobs are queued
    go func() {
        scheduler.ScheduleAll(ctx)
        close(workers)
    }()

    // Wait for all workers to complete
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results and handle errors
    return cp.collectResults(ctx, results, len(jobs))
}

// worker processes jobs from the worker channel
func (cp *ConcurrentProcessor) worker(ctx context.Context, workerID int, jobs <-chan *ProcessingJob, results chan<- *ProcessingResult) {
    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                return // Channel closed
            }

            result := cp.processJob(ctx, workerID, job)

            select {
            case results <- result:
            case <-ctx.Done():
                return
            }

        case <-ctx.Done():
            return
        }
    }
}

// processJob handles a single file processing job
func (cp *ConcurrentProcessor) processJob(ctx context.Context, workerID int, job *ProcessingJob) *ProcessingResult {
    start := time.Now()

    if cp.generator.options.Verbose {
        fmt.Printf("🔄 Worker %d processing: %s\n", workerID, filepath.Base(job.SrcPath))
    }

    // Check for cancellation
    select {
    case <-ctx.Done():
        return &ProcessingResult{
            JobID:   job.ID,
            Success: false,
            Error:   ctx.Err(),
            Duration: time.Since(start),
        }
    default:
    }

    // Process the file using the streaming processor
    err := cp.generator.streaming.ProcessFile(
        job.SrcPath,
        job.DestPath,
        cp.generator.engine,
        job.Context,
        job.Mode,
    )

    success := err == nil

    // Update progress
    cp.progressTracker.IncrementCompleted()

    if cp.generator.options.Verbose && success {
        fmt.Printf("✅ Worker %d completed: %s (%.2fs)\n",
            workerID, filepath.Base(job.SrcPath), time.Since(start).Seconds())
    }

    return &ProcessingResult{
        JobID:    job.ID,
        Success:  success,
        Error:    err,
        Duration: time.Since(start),
    }
}

// createWorkerPool creates a buffered channel for job distribution
func (cp *ConcurrentProcessor) createWorkerPool(ctx context.Context, scheduler *JobScheduler) chan *ProcessingJob {
    return make(chan *ProcessingJob, cp.workerCount*2) // Buffer for smooth operation
}

// collectResults gathers all processing results and handles errors
func (cp *ConcurrentProcessor) collectResults(ctx context.Context, results <-chan *ProcessingResult, expectedCount int) error {
    var errors []error
    successCount := 0

    for i := 0; i < expectedCount; i++ {
        select {
        case result := <-results:
            if result.Success {
                successCount++
            } else {
                errors = append(errors, fmt.Errorf("job %s failed: %w", result.JobID, result.Error))
            }

        case <-ctx.Done():
            return fmt.Errorf("processing cancelled: %w", ctx.Err())
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("processing completed with %d errors out of %d files: %v",
            len(errors), expectedCount, errors)
    }

    if cp.generator.options.Verbose {
        fmt.Printf("🎉 Concurrent processing completed: %d files processed successfully\n", successCount)
    }

    return nil
}

// JobScheduler handles dependency-aware job scheduling
type JobScheduler struct {
    jobs         []*ProcessingJob
    dependencies map[string][]string
    completed    map[string]bool
    workerCount  int
    mu           sync.Mutex
}

// NewJobScheduler creates a new job scheduler
func NewJobScheduler(jobs []*ProcessingJob, workerCount int) *JobScheduler {
    scheduler := &JobScheduler{
        jobs:        jobs,
        dependencies: make(map[string][]string),
        completed:   make(map[string]bool),
        workerCount: workerCount,
    }

    // Build dependency map
    for _, job := range jobs {
        if len(job.Dependencies) > 0 {
            scheduler.dependencies[job.ID] = job.Dependencies
        }
    }

    return scheduler
}

// ScheduleAll schedules all jobs respecting dependencies
func (js *JobScheduler) ScheduleAll(ctx context.Context) {
    scheduled := make(map[string]bool)

    for len(scheduled) < len(js.jobs) {
        progress := false

        for _, job := range js.jobs {
            if scheduled[job.ID] {
                continue
            }

            if js.canSchedule(job.ID) {
                select {
                case <-ctx.Done():
                    return
                default:
                    scheduled[job.ID] = true
                    progress = true
                }
            }
        }

        if !progress {
            // Deadlock detection - schedule remaining jobs anyway
            for _, job := range js.jobs {
                if !scheduled[job.ID] {
                    scheduled[job.ID] = true
                }
            }
            break
        }
    }
}

// canSchedule checks if a job's dependencies are satisfied
func (js *JobScheduler) canSchedule(jobID string) bool {
    js.mu.Lock()
    defer js.mu.Unlock()

    dependencies, hasDeps := js.dependencies[jobID]
    if !hasDeps {
        return true // No dependencies
    }

    for _, depID := range dependencies {
        if !js.completed[depID] {
            return false
        }
    }

    return true
}

// MarkCompleted marks a job as completed
func (js *JobScheduler) MarkCompleted(jobID string) {
    js.mu.Lock()
    defer js.mu.Unlock()
    js.completed[jobID] = true
}

// ProgressTracker provides thread-safe progress tracking
type ProgressTracker struct {
    total     int
    completed int
    mu        sync.RWMutex
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker() *ProgressTracker {
    return &ProgressTracker{}
}

// SetTotal sets the total number of items to process
func (pt *ProgressTracker) SetTotal(total int) {
    pt.mu.Lock()
    defer pt.mu.Unlock()
    pt.total = total
}

// IncrementCompleted increments the completed count
func (pt *ProgressTracker) IncrementCompleted() {
    pt.mu.Lock()
    defer pt.mu.Unlock()
    pt.completed++
}

// GetProgress returns current progress as (completed, total, percentage)
func (pt *ProgressTracker) GetProgress() (int, int, float64) {
    pt.mu.RLock()
    defer pt.mu.RUnlock()

    if pt.total == 0 {
        return 0, 0, 0
    }

    percentage := float64(pt.completed) / float64(pt.total) * 100
    return pt.completed, pt.total, percentage
}
```

**Update**: `internal/generator/generator.go`
```go
// Add concurrent processor to Generator
type Generator struct {
    engine     Engine
    options    *Options
    streaming  *StreamingProcessor
    detector   *engine.TemplateDetector
    concurrent *ConcurrentProcessor
}

// Update Generate method to use concurrent processing
func (g *Generator) Generate(templatePath, outputPath string, context map[string]interface{}) error {
    if g.options.Verbose {
        fmt.Printf("🚀 Generating project from template: %s\n", templatePath)
    }

    // Collect all files to process
    jobs, err := g.collectProcessingJobs(templatePath, outputPath, context)
    if err != nil {
        return fmt.Errorf("failed to collect processing jobs: %w", err)
    }

    if len(jobs) == 0 {
        fmt.Println("ℹ️  No files to process")
        return nil
    }

    // Use concurrent processing if enabled and beneficial
    if g.options.ConcurrentProcessing && len(jobs) > g.options.ConcurrentThreshold {
        ctx, cancel := context.WithTimeout(context.Background(), g.options.ProcessingTimeout)
        defer cancel()

        return g.concurrent.ProcessFiles(ctx, jobs)
    }

    // Fall back to sequential processing
    return g.processSequentially(jobs)
}

// collectProcessingJobs walks the template and creates processing jobs
func (g *Generator) collectProcessingJobs(templatePath, outputPath string, context map[string]interface{}) ([]*ProcessingJob, error) {
    var jobs []*ProcessingJob
    jobID := 0

    err := filepath.Walk(templatePath, func(srcPath string, info os.FileInfo, err error) error {
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

    if err != nil {
        return nil, err
    }

    // Sort jobs by priority (high priority first)
    sort.Slice(jobs, func(i, j int) bool {
        return jobs[i].Priority > jobs[j].Priority
    })

    return jobs, nil
}

// determineJobPriority assigns priority based on file characteristics
func (g *Generator) determineJobPriority(filePath string, size int64) JobPriority {
    ext := filepath.Ext(filePath)

    // High priority for configuration and small template files
    if size < 1024 || isConfigFile(ext) {
        return HighPriority
    }

    // Low priority for large binary files
    if size > 10*1024*1024 || isBinaryFileExt(ext) {
        return LowPriority
    }

    return NormalPriority
}

// processSequentially handles files one by one (fallback)
func (g *Generator) processSequentially(jobs []*ProcessingJob) error {
    for i, job := range jobs {
        if g.options.Verbose {
            fmt.Printf("📄 Processing file %d/%d: %s\n", i+1, len(jobs), filepath.Base(job.SrcPath))
        }

        err := g.streaming.ProcessFile(job.SrcPath, job.DestPath, g.engine, job.Context, job.Mode)
        if err != nil {
            return fmt.Errorf("failed to process %s: %w", job.SrcPath, err)
        }
    }

    return nil
}
```

## Configuration Options

**Add to Options struct**:
```go
type Options struct {
    // ... existing fields ...

    // Concurrent processing configuration
    ConcurrentProcessing bool
    ConcurrentThreshold  int           // Minimum files to enable concurrency
    MaxWorkers          int           // Maximum worker goroutines
    ProcessingTimeout   time.Duration // Timeout for entire operation
    MaxMemoryMB         int           // Memory limit for processing
}
```

## Files to Modify
- `internal/generator/concurrent.go` (new file)
- `internal/generator/concurrent_test.go` (new file)
- `internal/generator/generator.go` (update to use concurrent processing)
- `cmd/new.go` (add concurrency options)

## Testing Strategy
- Unit tests for worker pool functionality
- Integration tests with large template sets
- Stress tests with many small files
- Memory usage tests under load
- Concurrency safety tests
- Performance benchmarks vs sequential processing

## Performance Targets
- **Throughput**: >60% improvement for templates with >100 files
- **Resource Usage**: Configurable memory limits respected
- **Scalability**: Linear improvement with CPU core count
- **Latency**: First files processed quickly (priority scheduling)

## Error Handling
- Individual file failures don't stop entire generation
- Aggregated error reporting at completion
- Graceful handling of context cancellation
- Resource cleanup on errors

## Definition of Done
- Worker pool system processes files concurrently
- Dependency-aware scheduling implemented
- Progress tracking works correctly with concurrent operations
- Resource limits prevent system overload
- Performance improvement >60% verified for large templates
- Error handling is robust and informative
- All tests pass including stress tests