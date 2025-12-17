// Package prompt provides interactive prompt functionality using Bubble Tea.
package prompt

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// TextPrompt is a simple text input prompt
type TextPrompt struct {
	prompt  string
	Value   string
	Default interface{}
	done    bool
}

// NewTextPrompt creates a new text prompt with the given prompt text and optional default value
func NewTextPrompt(prompt string, defaultValue interface{}) TextPrompt {
	defaultStr := ""
	if defaultValue != nil {
		defaultStr = fmt.Sprintf("%v", defaultValue)
	}
	return TextPrompt{
		prompt:  prompt,
		Value:   defaultStr,
		Default: defaultValue,
	}
}

// Init initializes the prompt (returns no command)
func (m TextPrompt) Init() tea.Cmd {
	return nil
}

// Update handles user input and returns the updated model
func (m TextPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Value == "" && m.Default != nil {
				m.Value = fmt.Sprintf("%v", m.Default)
			}
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyBackspace:
			if len(m.Value) > 0 {
				m.Value = m.Value[:len(m.Value)-1]
			}
		default:
			m.Value += msg.String()
		}
	}
	return m, nil
}

// View renders the prompt as a string
func (m TextPrompt) View() string {
	if m.done {
		return ""
	}

	defaultHint := ""
	if m.Default != nil && m.Default != "" {
		defaultHint = fmt.Sprintf(" (default: %v)", m.Default)
	}

	return fmt.Sprintf("%s%s: %s", m.prompt, defaultHint, m.Value)
}
