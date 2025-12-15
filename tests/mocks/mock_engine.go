package mocks

import (
	"fmt"
	"sync"
)

// MockEngine is a test implementation of the Engine interface.
// It tracks calls and allows configuring responses.
type MockEngine struct {
	mu                 sync.Mutex
	renderCalls        []RenderCall
	renderFilesCalls   []RenderFileCall
	renderResponse     string
	renderErr          error
	renderFileResponse string
	renderFileErr      error
	shouldFail         bool
}

// RenderCall represents a call to Render()
type RenderCall struct {
	Template string
	Context  map[string]interface{}
}

// RenderFileCall represents a call to RenderFile()
type RenderFileCall struct {
	FilePath string
	Context  map[string]interface{}
}

// NewMockEngine creates a new mock engine with default responses.
func NewMockEngine() *MockEngine {
	return &MockEngine{
		renderResponse:     "rendered_output",
		renderFileResponse: "rendered_file_output",
	}
}

// Render implements the Engine interface for testing.
func (m *MockEngine) Render(template string, context map[string]interface{}) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.renderCalls = append(m.renderCalls, RenderCall{
		Template: template,
		Context:  context,
	})

	if m.renderErr != nil {
		return "", m.renderErr
	}

	if m.shouldFail {
		return "", fmt.Errorf("mock render error")
	}

	return m.renderResponse, nil
}

// RenderFile implements the Engine interface for testing.
func (m *MockEngine) RenderFile(filePath string, context map[string]interface{}) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.renderFilesCalls = append(m.renderFilesCalls, RenderFileCall{
		FilePath: filePath,
		Context:  context,
	})

	if m.renderFileErr != nil {
		return "", m.renderFileErr
	}

	if m.shouldFail {
		return "", fmt.Errorf("mock render file error")
	}

	return m.renderFileResponse, nil
}

// SetRenderResponse sets the response for Render() calls.
func (m *MockEngine) SetRenderResponse(response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.renderResponse = response
}

// SetRenderError sets the error response for Render() calls.
func (m *MockEngine) SetRenderError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.renderErr = err
}

// SetRenderFileResponse sets the response for RenderFile() calls.
func (m *MockEngine) SetRenderFileResponse(response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.renderFileResponse = response
}

// SetRenderFileError sets the error response for RenderFile() calls.
func (m *MockEngine) SetRenderFileError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.renderFileErr = err
}

// SetShouldFail configures the mock to return errors for all calls.
func (m *MockEngine) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// GetRenderCalls returns all recorded Render() calls.
func (m *MockEngine) GetRenderCalls() []RenderCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	calls := make([]RenderCall, len(m.renderCalls))
	copy(calls, m.renderCalls)
	return calls
}

// GetRenderFileCalls returns all recorded RenderFile() calls.
func (m *MockEngine) GetRenderFileCalls() []RenderFileCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	calls := make([]RenderFileCall, len(m.renderFilesCalls))
	copy(calls, m.renderFilesCalls)
	return calls
}

// Reset clears all recorded calls and resets to default state.
func (m *MockEngine) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.renderCalls = nil
	m.renderFilesCalls = nil
	m.renderResponse = "rendered_output"
	m.renderErr = nil
	m.renderFileResponse = "rendered_file_output"
	m.renderFileErr = nil
	m.shouldFail = false
}
