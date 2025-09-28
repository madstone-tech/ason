package engine

import (
	"fmt"

	"github.com/flosch/pongo2/v6"
)

// Engine defines the interface for template engines
type Engine interface {
	Render(template string, context map[string]interface{}) (string, error)
	RenderFile(filepath string, context map[string]interface{}) (string, error)
}

// Pongo2Engine implements Engine using Pongo2
type Pongo2Engine struct{}

// NewPongo2Engine creates a new Pongo2 templating engine
func NewPongo2Engine() *Pongo2Engine {
	return &Pongo2Engine{}
}

// Render renders a template string with the given context
func (e *Pongo2Engine) Render(template string, context map[string]interface{}) (string, error) {
	tpl, err := pongo2.FromString(template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.Execute(pongo2.Context(context))
}

// RenderFile renders a template file with the given context
func (e *Pongo2Engine) RenderFile(filepath string, context map[string]interface{}) (string, error) {
	tpl, err := pongo2.FromFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to load template file: %w", err)
	}

	return tpl.Execute(pongo2.Context(context))
}
