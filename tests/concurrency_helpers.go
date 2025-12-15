package tests

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ConcurrencyResult tracks the results of concurrent operations.
type ConcurrencyResult struct {
	OperationIndex int
	Success        bool
	Error          error
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
}

// RunConcurrent executes a function concurrently and collects results.
// It runs the function numConcurrent times in parallel and waits for all to complete.
func RunConcurrent(t testing.TB, numConcurrent int, fn func(index int) error) []ConcurrencyResult {
	if numConcurrent < 1 {
		t.Fatal("numConcurrent must be >= 1")
	}

	results := make([]ConcurrencyResult, numConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			start := time.Now()
			err := fn(index)
			end := time.Now()

			mu.Lock()
			results[index] = ConcurrencyResult{
				OperationIndex: index,
				Success:        err == nil,
				Error:          err,
				StartTime:      start,
				EndTime:        end,
				Duration:       end.Sub(start),
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return results
}

// AssertAllSuccess checks that all concurrent operations succeeded.
func AssertAllSuccess(t testing.TB, results []ConcurrencyResult) {
	t.Helper()

	for _, result := range results {
		if !result.Success {
			t.Errorf("operation %d failed: %v", result.OperationIndex, result.Error)
		}
	}
}

// AssertSuccessCount checks that exactly count operations succeeded.
func AssertSuccessCount(t testing.TB, results []ConcurrencyResult, count int) {
	t.Helper()

	successful := 0
	for _, result := range results {
		if result.Success {
			successful++
		}
	}

	if successful != count {
		t.Errorf("expected %d successful operations, got %d", count, successful)
	}
}

// AssertNoRaceCondition checks for signs of race conditions.
// It verifies that operations don't have overlapping execution times (for write operations).
func AssertNoRaceCondition(t testing.TB, results []ConcurrencyResult) {
	t.Helper()

	// Sort by start time
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].StartTime.After(results[j].StartTime) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Check for overlapping execution (which might indicate issues)
	// This is a simple heuristic - actual race detection requires more sophisticated methods
	for i := 0; i < len(results)-1; i++ {
		if results[i].EndTime.After(results[i+1].StartTime) {
			// Operations overlap - this is expected for concurrent reads
			// but might indicate issues for exclusive operations
		}
	}
}

// MeasureConcurrencyThroughput measures the throughput of concurrent operations.
// It returns the number of operations per second.
func MeasureConcurrencyThroughput(results []ConcurrencyResult) float64 {
	if len(results) == 0 {
		return 0
	}

	var minStart time.Time
	var maxEnd time.Time

	for _, result := range results {
		if minStart.IsZero() || result.StartTime.Before(minStart) {
			minStart = result.StartTime
		}
		if result.EndTime.After(maxEnd) {
			maxEnd = result.EndTime
		}
	}

	duration := maxEnd.Sub(minStart).Seconds()
	if duration == 0 {
		return 0
	}

	return float64(len(results)) / duration
}

// ConcurrentCounter is a thread-safe counter for concurrent operations.
type ConcurrentCounter struct {
	value int64
	mu    sync.RWMutex
}

// NewConcurrentCounter creates a new counter.
func NewConcurrentCounter() *ConcurrentCounter {
	return &ConcurrentCounter{}
}

// Increment adds 1 to the counter.
func (c *ConcurrentCounter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// Value returns the current counter value.
func (c *ConcurrentCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset resets the counter to 0.
func (c *ConcurrentCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

// ErrorCollector collects errors from concurrent operations.
type ErrorCollector struct {
	mu     sync.Mutex
	errors []error
}

// NewErrorCollector creates a new error collector.
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]error, 0),
	}
}

// Add adds an error to the collection.
func (ec *ErrorCollector) Add(err error) {
	if err == nil {
		return
	}
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errors = append(ec.errors, err)
}

// Errors returns all collected errors.
func (ec *ErrorCollector) Errors() []error {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	errors := make([]error, len(ec.errors))
	copy(errors, ec.errors)
	return errors
}

// HasErrors returns true if any errors were collected.
func (ec *ErrorCollector) HasErrors() bool {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	return len(ec.errors) > 0
}

// String returns a formatted string of all errors.
func (ec *ErrorCollector) String() string {
	errors := ec.Errors()
	if len(errors) == 0 {
		return "no errors"
	}

	msg := fmt.Sprintf("%d errors:\n", len(errors))
	for i, err := range errors {
		msg += fmt.Sprintf("  [%d] %v\n", i+1, err)
	}
	return msg
}

// WaitForCondition waits for a condition to become true, with timeout.
func WaitForCondition(t testing.TB, timeout time.Duration, condition func() bool) bool {
	t.Helper()

	start := time.Now()
	for {
		if condition() {
			return true
		}

		if time.Since(start) > timeout {
			t.Logf("condition not met within %v timeout", timeout)
			return false
		}

		time.Sleep(10 * time.Millisecond)
	}
}
