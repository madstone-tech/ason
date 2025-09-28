package cmd

import (
	"fmt"
	"os"

	"github.com/madstone-tech/ason/internal/engine"
	"github.com/madstone-tech/ason/internal/generator"
	"github.com/madstone-tech/ason/internal/registry"
	"github.com/spf13/cobra"
)

var (
	outputDir  string
	noInput    bool
	extraVars  map[string]string
	configFile string
	skipHooks  bool
	dryRun     bool
)

var newCmd = &cobra.Command{
	Use:   "new [template] [output]",
	Short: "Create a new project from a template",
	Long: `Create a new project from a template.

Examples:
  ason new golang-service my-service
  ason new ./my-template ./output`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runNew,
}

func init() {
	newCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	newCmd.Flags().BoolVar(&noInput, "no-input", false, "Don't prompt for variables")
	newCmd.Flags().StringToStringVar(&extraVars, "var", nil, "Set variables (key=value)")
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

	// Generate with placeholder context
	context := make(map[string]interface{})
	for k, v := range extraVars {
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
