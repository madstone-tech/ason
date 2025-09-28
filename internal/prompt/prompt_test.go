package prompt

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewTextPrompt(t *testing.T) {
	tests := []struct {
		name         string
		prompt       string
		defaultValue interface{}
		wantValue    string
	}{
		{
			name:         "with string default",
			prompt:       "Enter name:",
			defaultValue: "test",
			wantValue:    "test",
		},
		{
			name:         "with nil default",
			prompt:       "Enter value:",
			defaultValue: nil,
			wantValue:    "",
		},
		{
			name:         "with integer default",
			prompt:       "Enter number:",
			defaultValue: 42,
			wantValue:    "42",
		},
		{
			name:         "with boolean default",
			prompt:       "Enable feature:",
			defaultValue: true,
			wantValue:    "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := NewTextPrompt(tt.prompt, tt.defaultValue)

			if prompt.prompt != tt.prompt {
				t.Errorf("TextPrompt.prompt = %v, want %v", prompt.prompt, tt.prompt)
			}

			if prompt.Value != tt.wantValue {
				t.Errorf("TextPrompt.Value = %v, want %v", prompt.Value, tt.wantValue)
			}

			if prompt.Default != tt.defaultValue {
				t.Errorf("TextPrompt.Default = %v, want %v", prompt.Default, tt.defaultValue)
			}

			if prompt.done {
				t.Error("TextPrompt.done should be false initially")
			}
		})
	}
}

func TestTextPrompt_Init(t *testing.T) {
	prompt := NewTextPrompt("Test:", "default")
	cmd := prompt.Init()
	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestTextPrompt_Update_Enter(t *testing.T) {
	tests := []struct {
		name         string
		initialValue string
		defaultValue interface{}
		wantValue    string
		wantDone     bool
	}{
		{
			name:         "enter with custom value",
			initialValue: "custom",
			defaultValue: "default",
			wantValue:    "custom",
			wantDone:     true,
		},
		{
			name:         "enter with empty value uses default",
			initialValue: "",
			defaultValue: "default",
			wantValue:    "default",
			wantDone:     true,
		},
		{
			name:         "enter with empty value and nil default",
			initialValue: "",
			defaultValue: nil,
			wantValue:    "",
			wantDone:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := NewTextPrompt("Test:", tt.defaultValue)
			prompt.Value = tt.initialValue

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			model, cmd := prompt.Update(msg)

			updatedPrompt := model.(TextPrompt)

			if updatedPrompt.Value != tt.wantValue {
				t.Errorf("After Enter, Value = %v, want %v", updatedPrompt.Value, tt.wantValue)
			}

			if updatedPrompt.done != tt.wantDone {
				t.Errorf("After Enter, done = %v, want %v", updatedPrompt.done, tt.wantDone)
			}

			if cmd == nil {
				t.Error("Enter should return tea.Quit command, got nil")
			}
		})
	}
}

func TestTextPrompt_Update_CtrlC(t *testing.T) {
	prompt := NewTextPrompt("Test:", "default")

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	model, cmd := prompt.Update(msg)

	if cmd == nil {
		t.Error("Ctrl+C should return tea.Quit command, got nil")
	}

	// Should not change other fields
	updatedPrompt := model.(TextPrompt)
	if updatedPrompt.done {
		t.Error("Ctrl+C should not mark as done")
	}
}

func TestTextPrompt_Update_Esc(t *testing.T) {
	prompt := NewTextPrompt("Test:", "default")

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	model, cmd := prompt.Update(msg)

	if cmd == nil {
		t.Error("Esc should return tea.Quit command, got nil")
	}

	// Should not change other fields
	updatedPrompt := model.(TextPrompt)
	if updatedPrompt.done {
		t.Error("Esc should not mark as done")
	}
}

func TestTextPrompt_Update_Backspace(t *testing.T) {
	tests := []struct {
		name         string
		initialValue string
		wantValue    string
	}{
		{
			name:         "backspace with text",
			initialValue: "hello",
			wantValue:    "hell",
		},
		{
			name:         "backspace with single char",
			initialValue: "a",
			wantValue:    "",
		},
		{
			name:         "backspace with empty",
			initialValue: "",
			wantValue:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := NewTextPrompt("Test:", nil)
			prompt.Value = tt.initialValue

			msg := tea.KeyMsg{Type: tea.KeyBackspace}
			model, cmd := prompt.Update(msg)

			updatedPrompt := model.(TextPrompt)

			if updatedPrompt.Value != tt.wantValue {
				t.Errorf("After Backspace, Value = %v, want %v", updatedPrompt.Value, tt.wantValue)
			}

			if cmd != nil {
				t.Error("Backspace should return nil command")
			}

			if updatedPrompt.done {
				t.Error("Backspace should not mark as done")
			}
		})
	}
}

func TestTextPrompt_Update_RegularKey(t *testing.T) {
	prompt := NewTextPrompt("Test:", nil)
	prompt.Value = "hello"

	// Simulate typing 'a'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	model, cmd := prompt.Update(msg)

	updatedPrompt := model.(TextPrompt)

	if updatedPrompt.Value != "helloa" {
		t.Errorf("After typing 'a', Value = %v, want %v", updatedPrompt.Value, "helloa")
	}

	if cmd != nil {
		t.Error("Regular key should return nil command")
	}

	if updatedPrompt.done {
		t.Error("Regular key should not mark as done")
	}
}

func TestTextPrompt_Update_OtherMessage(t *testing.T) {
	prompt := NewTextPrompt("Test:", "default")
	originalValue := prompt.Value

	// Send a non-key message
	msg := "some other message"
	model, cmd := prompt.Update(msg)

	updatedPrompt := model.(TextPrompt)

	if updatedPrompt.Value != originalValue {
		t.Errorf("Non-key message should not change Value")
	}

	if cmd != nil {
		t.Error("Non-key message should return nil command")
	}

	if updatedPrompt.done {
		t.Error("Non-key message should not mark as done")
	}
}

func TestTextPrompt_View(t *testing.T) {
	tests := []struct {
		name         string
		prompt       string
		value        string
		defaultValue interface{}
		done         bool
		wantContains []string
		wantEmpty    bool
	}{
		{
			name:         "normal view with default",
			prompt:       "Enter name",
			value:        "test",
			defaultValue: "default",
			done:         false,
			wantContains: []string{"Enter name", "(default: default)", "test"},
			wantEmpty:    false,
		},
		{
			name:         "view with nil default",
			prompt:       "Enter value",
			value:        "input",
			defaultValue: nil,
			done:         false,
			wantContains: []string{"Enter value", "input"},
			wantEmpty:    false,
		},
		{
			name:         "view with empty default",
			prompt:       "Enter something",
			value:        "text",
			defaultValue: "",
			done:         false,
			wantContains: []string{"Enter something", "text"},
			wantEmpty:    false,
		},
		{
			name:         "view when done",
			prompt:       "Enter name",
			value:        "test",
			defaultValue: "default",
			done:         true,
			wantContains: nil,
			wantEmpty:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := TextPrompt{
				prompt:  tt.prompt,
				Value:   tt.value,
				Default: tt.defaultValue,
				done:    tt.done,
			}

			view := prompt.View()

			if tt.wantEmpty {
				if view != "" {
					t.Errorf("View() should be empty when done, got %v", view)
				}
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(view, want) {
					t.Errorf("View() = %v, should contain %v", view, want)
				}
			}
		})
	}
}

func TestTextPrompt_Struct(t *testing.T) {
	// Test that TextPrompt struct has all expected fields
	prompt := TextPrompt{
		prompt:  "test",
		Value:   "value",
		Default: "default",
		done:    true,
	}

	if prompt.prompt != "test" {
		t.Error("prompt field not set correctly")
	}

	if prompt.Value != "value" {
		t.Error("Value field not set correctly")
	}

	if prompt.Default != "default" {
		t.Error("Default field not set correctly")
	}

	if !prompt.done {
		t.Error("done field not set correctly")
	}
}
