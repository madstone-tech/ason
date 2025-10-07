package cmd

import (
	"fmt"
	"os"

	"github.com/madstone-tech/ason/internal/engine"
	"github.com/madstone-tech/ason/internal/generator"
	"github.com/madstone-tech/ason/internal/registry"
	"github.com/madstone-tech/ason/internal/varfile"
	"github.com/spf13/cobra"
)

var (
	outputDir  string
	noInput    bool
	extraVars  map[string]string
	varFile    string
	configFile string
	skipHooks  bool
	dryRun     bool
)

var newCmd = &cobra.Command{
	Use:   "new [template] [output]",
	Short: "Create a new project from a template",
	Long: `Create a new project from a template.

Examples:
  # Create from registry template
  ason new golang-service my-service

  # Create from local template
  ason new ./my-template ./output

  # Use variables from file
  ason new lambda-waf-ipset ./output --var-file prod.toml

  # Mix file variables with CLI overrides
  ason new lambda-waf-ipset ./output --var-file base.toml --var environment=prod`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runNew,
}

func init() {
	newCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	newCmd.Flags().BoolVar(&noInput, "no-input", false, "Don't prompt for variables")
	newCmd.Flags().StringToStringVar(&extraVars, "var", nil, "Set variables (key=value)")
	newCmd.Flags().StringVarP(&varFile, "var-file", "f", "", "Load variables from file (TOML, YAML, or JSON)")
	newCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be generated")
}

func runNew(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	if len(args) > 1 {
		outputDir = args[1]
	}

	fmt.Println("※ The ason shakes, preparing transformation...")

	// Get template path
	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	templatePath, err := reg.Get(templateName)
	if err != nil {
		// Try as direct path
		if info, err := os.Stat(templateName); err == nil && info.IsDir() {
			templatePath = templateName
		} else {
			return fmt.Errorf("template not found: %s", templateName)
		}
	}

	// Create a simple template object
	tmpl := &generator.Template{
		Path: templatePath,
	}

	// Create generator
	gen := generator.New(tmpl, engine.NewPongo2Engine())

	// Load variables from file if specified
	var fileVars map[string]string
	if varFile != "" {
		var err error
		fileVars, err = varfile.Load(varFile)
		if err != nil {
			return fmt.Errorf("failed to load variables from file: %w", err)
		}
	}

	// Merge variables (CLI vars override file vars)
	mergedVars := varfile.Merge(fileVars, extraVars)

	// Generate with context
	context := make(map[string]interface{})
	for k, v := range mergedVars {
		context[k] = v
	}

	if err := gen.Generate(outputDir, context, generator.Options{
		DryRun: dryRun,
	}); err != nil {
		return err
	}

	if !dryRun {
		fmt.Println("※ The rhythm is complete! Project manifested successfully!")
	}

	return nil
}
