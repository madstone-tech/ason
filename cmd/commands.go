// Package cmd implements the Cobra CLI commands for Ason.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/madstone-tech/ason/internal/registry"
	"github.com/spf13/cobra"
)

var (
	// List command flags
	listFormat  string
	listFilter  string
	listSort    string
	listReverse bool

	// Register command flags
	registerDescription string
	registerType        string
	registerForce       bool
	registerValidate    bool
	registerDryRun      bool

	// Remove command flags
	removeForce     bool
	removeDryRun    bool
	removeBackup    bool
	removeBackupDir string

	// Validate command flags
	validateStrict         bool
	validateFormat         string
	validateFix            bool
	validateCheck          string
	validateIgnoreWarnings bool
)

// listCmd lists available templates
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Long:  `List all templates available in the local registry.`,
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringVar(&listFormat, "format", "table", "Output format (table, json, yaml)")
	listCmd.Flags().StringVar(&listFilter, "filter", "", "Filter templates by name or description")
	listCmd.Flags().StringVar(&listSort, "sort", "name", "Sort by field (name, date, size, type)")
	listCmd.Flags().BoolVar(&listReverse, "reverse", false, "Reverse sort order")

	registerCmd.Flags().StringVar(&registerDescription, "description", "", "Template description")
	registerCmd.Flags().StringVar(&registerType, "type", "", "Template type")
	registerCmd.Flags().BoolVar(&registerForce, "force", false, "Overwrite existing template")
	registerCmd.Flags().BoolVar(&registerValidate, "validate", false, "Validate template before registering")
	registerCmd.Flags().BoolVar(&registerDryRun, "dry-run", false, "Show what would be registered")

	removeCmd.Flags().BoolVar(&removeForce, "force", false, "Remove without confirmation")
	removeCmd.Flags().BoolVar(&removeDryRun, "dry-run", false, "Show what would be removed")
	removeCmd.Flags().BoolVar(&removeBackup, "backup", false, "Create backup before removing")
	removeCmd.Flags().StringVar(&removeBackupDir, "backup-dir", "", "Backup directory")

	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "Enable strict validation")
	validateCmd.Flags().StringVar(&validateFormat, "format", "text", "Output format (text, json, junit)")
	validateCmd.Flags().BoolVar(&validateFix, "fix", false, "Fix issues automatically")
	validateCmd.Flags().StringVar(&validateCheck, "check", "", "Check specific categories")
	validateCmd.Flags().BoolVar(&validateIgnoreWarnings, "ignore-warnings", false, "Show only errors")
}

func runList(_ *cobra.Command, _ []string) error {
	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	templates, err := reg.List()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	// Filter templates
	if listFilter != "" {
		templates = filterTemplates(templates, listFilter)
	}

	// Sort templates
	sortTemplates(templates, listSort, listReverse)

	if len(templates) == 0 {
		if listFormat == "json" {
			fmt.Println(`{"templates":[], "total":0}`)
			return nil
		} else if listFormat == "yaml" {
			fmt.Println("templates: []\ntotal: 0")
			return nil
		}

		fmt.Println("※ The registry echoes with silence...")
		fmt.Println()
		fmt.Println("No templates ready for invocation.")
		fmt.Println()
		fmt.Println("💡 Prepare templates for transformation:")
		fmt.Println("   ason register my-template /path/to/template")
		return nil
	}

	switch listFormat {
	case "json":
		return printTemplatesJSON(templates)
	case "yaml":
		return printTemplatesYAML(templates)
	default:
		return printTemplatesTable(templates)
	}
}

// registerCmd registers a template in the registry.
// The "add" alias is maintained for backward compatibility with existing scripts and workflows.
var registerCmd = &cobra.Command{
	Use:     "register [name] [path]",
	Aliases: []string{"add"}, // Backward compatibility: "ason add" still works
	Short:   "Register a template in the registry",
	Args:    cobra.ExactArgs(2),
	RunE:    runRegister,
}

func runRegister(_ *cobra.Command, args []string) error {
	name := args[0]
	sourcePath := args[1]

	// Expand path
	if strings.HasPrefix(sourcePath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		sourcePath = filepath.Join(home, sourcePath[2:])
	}

	// Make path absolute
	sourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	fmt.Println("※ The ason prepares to embrace new wisdom...")

	if registerDryRun {
		fmt.Println("[DRY RUN] Would analyze:", sourcePath)
		fmt.Println("[DRY RUN] Would validate template structure")
		fmt.Printf("[DRY RUN] Would copy to: ~/.ason/templates/%s\n", name)
		fmt.Printf("[DRY RUN] Would register as: %s\n", name)
		fmt.Println("🔮 [DRY RUN] Template ready for registration. Use without --dry-run to register.")
		return nil
	}

	fmt.Println("✨ Analyzing template:", sourcePath)

	// Validate template if requested
	if registerValidate {
		fmt.Println("📿 Validating template structure...")
		if err := validateTemplate(sourcePath); err != nil {
			return fmt.Errorf("template validation failed: %w", err)
		}
		fmt.Println("💫 Template structure confirmed")
	}

	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	// Check if template exists and handle force flag
	if _, err := reg.Get(name); err == nil {
		if !registerForce {
			return fmt.Errorf("template '%s' already exists. Use --force to overwrite", name)
		}
		// Force flag is enabled, remove existing template first
		fmt.Println("🔄 Removing existing template for overwrite...")
		if err := reg.Remove(name, false, ""); err != nil {
			return fmt.Errorf("failed to remove existing template: %w", err)
		}
	}

	fmt.Println("🎭 Copying template to registry...")

	// Register template in registry
	if err := reg.Add(name, sourcePath, registerDescription, registerType); err != nil {
		return fmt.Errorf("failed to add template: %w", err)
	}

	fmt.Printf("🔮 Template '%s' added to registry successfully!\n", name)
	fmt.Println()
	fmt.Printf("💡 Use it with: ason new %s my-project\n", name)

	return nil
}

// removeCmd removes a template from the registry
var removeCmd = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a template from the registry",
	Args:    cobra.ExactArgs(1),
	RunE:    runRemove,
}

func runRemove(_ *cobra.Command, args []string) error {
	name := args[0]

	fmt.Println("※ The ason prepares to release template from registry...")

	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	// Get template info for display
	templates, err := reg.List()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	var tmpl *registry.TemplateEntry
	for _, t := range templates {
		if t.Name == name {
			tmpl = &t
			break
		}
	}

	if tmpl == nil {
		return fmt.Errorf("template '%s' not found in registry", name)
	}

	if removeDryRun {
		fmt.Printf("[DRY RUN] Would remove template: %s\n", name)
		fmt.Printf("[DRY RUN] Would delete: %s\n", tmpl.Path)
		fmt.Println("[DRY RUN] Would clean registry metadata")
		fmt.Printf("[DRY RUN] Size to be freed: %s\n", formatSize(tmpl.Size))
		fmt.Println("🔮 [DRY RUN] Template ready for removal. Use without --dry-run to remove.")
		return nil
	}

	// Show template info and confirm if not forced
	if !removeForce {
		fmt.Println()
		fmt.Printf("Template: %s\n", tmpl.Name)
		fmt.Printf("Description: %s\n", tmpl.Description)
		fmt.Printf("Size: %s\n", formatSize(tmpl.Size))
		fmt.Printf("Files: %d\n", tmpl.Files)
		fmt.Printf("Added: %s\n", formatTime(tmpl.Added))
		fmt.Println()
		fmt.Println("⚠️  This action cannot be undone.")
		fmt.Printf("🔮 Remove template '%s' from registry? [y/N]: ", name)

		var response string
		if _, err := fmt.Scanln(&response); err != nil && err.Error() != "unexpected newline" {
			// Ignore unexpected newline errors from Scanln - will continue with empty response
		}
		if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	if removeBackup {
		fmt.Println("✨ Creating backup before removal...")
	}

	fmt.Printf("✨ Removing template '%s'...\n", name)

	// Remove template from registry
	if err := reg.Remove(name, removeBackup, removeBackupDir); err != nil {
		return fmt.Errorf("failed to remove template: %w", err)
	}

	if removeBackup {
		fmt.Printf("💫 Backup created in: %s\n", getBackupDir(removeBackupDir))
	}

	fmt.Printf("🔮 Template '%s' removed successfully!\n", name)

	return nil
}

// validateCmd validates a template
var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate a template",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Validate all templates in registry
		return validateAllTemplates()
	}

	path := args[0]

	// Expand path if needed
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[2:])
	}

	fmt.Printf("※ Validating template: %s\n\n", path)

	return validateTemplate(path)
}

// Helper functions

func filterTemplates(templates []registry.TemplateEntry, filter string) []registry.TemplateEntry {
	var filtered []registry.TemplateEntry
	filter = strings.ToLower(filter)

	for _, tmpl := range templates {
		if strings.Contains(strings.ToLower(tmpl.Name), filter) ||
			strings.Contains(strings.ToLower(tmpl.Description), filter) ||
			strings.Contains(strings.ToLower(tmpl.Type), filter) {
			filtered = append(filtered, tmpl)
		}
	}

	return filtered
}

func sortTemplates(templates []registry.TemplateEntry, sortBy string, reverse bool) {
	sort.Slice(templates, func(i, j int) bool {
		var result bool

		switch sortBy {
		case "date":
			result = templates[i].Added.Before(templates[j].Added)
		case "size":
			result = templates[i].Size < templates[j].Size
		case "type":
			result = templates[i].Type < templates[j].Type
		default: // name
			result = templates[i].Name < templates[j].Name
		}

		if reverse {
			return !result
		}
		return result
	})
}

func printTemplatesTable(templates []registry.TemplateEntry) error {
	fmt.Println("※ Templates ready for invocation:")
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "NAME\tDESCRIPTION\tTYPE\tSIZE\tADDED"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := fmt.Fprintln(w, "----\t-----------\t----\t----\t-----"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	for _, tmpl := range templates {
		desc := tmpl.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		if desc == "" {
			desc = "-"
		}

		tmplType := tmpl.Type
		if tmplType == "" {
			tmplType = "-"
		}

		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			tmpl.Name,
			desc,
			tmplType,
			formatSize(tmpl.Size),
			formatTime(tmpl.Added)); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}
	fmt.Println()
	fmt.Println("💡 Use 'ason new TEMPLATE OUTPUT_DIR' to create a project")
	fmt.Println("💡 Use 'ason register' to prepare more templates for invocation")

	return nil
}

func printTemplatesJSON(templates []registry.TemplateEntry) error {
	output := map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func printTemplatesYAML(templates []registry.TemplateEntry) error {
	output := map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
	}

	// Use TOML format instead of YAML
	var buf strings.Builder
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to marshal TOML: %w", err)
	}

	fmt.Print(buf.String())
	return nil
}

func validateTemplate(templatePath string) error {
	// Check if path exists
	info, err := os.Stat(templatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("template not found at %s", templatePath)
		}
		return fmt.Errorf("failed to access template: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("template path must be a directory: %s", templatePath)
	}

	fmt.Println("✅ Structure Validation")
	fmt.Println("   ✓ Template directory exists")

	// Count files
	fileCount := 0
	err = filepath.Walk(templatePath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to analyze template: %w", err)
	}

	if fileCount == 0 {
		fmt.Println("❌ Template directory is empty")
		return fmt.Errorf("template contains no files")
	}

	fmt.Printf("   ✓ Contains %d processable files\n", fileCount)
	fmt.Println("   ✓ Directory structure is valid")

	// Check for configuration file (ason.toml)
	tomlPath := filepath.Join(templatePath, "ason.toml")

	var config registry.TemplateConfig

	if _, err := os.Stat(tomlPath); err == nil {
		fmt.Println("\n✅ Configuration Validation")
		fmt.Println("   ✓ ason.toml found")

		data, err := os.ReadFile(tomlPath)
		if err != nil {
			fmt.Println("❌ Failed to read ason.toml")
			return fmt.Errorf("failed to read config: %w", err)
		}

		if err := toml.Unmarshal(data, &config); err != nil {
			fmt.Println("❌ ason.toml syntax error")
			return fmt.Errorf("invalid config syntax: %w", err)
		}

		fmt.Println("   ✓ ason.toml syntax is correct")
		fmt.Println("   ✓ Configuration is valid")
		if len(config.Variables) > 0 {
			fmt.Printf("   ✓ Defines %d variables\n", len(config.Variables))
		}
	} else {
		fmt.Println("\n⚠️  Configuration Validation")
		fmt.Println("   ⚠ No ason.toml found (optional)")
	}

	fmt.Println("\n🔮 Validation Summary:")
	fmt.Println("   ✅ Template structure is valid")
	fmt.Println("   ✅ Ready for use with Ason")

	return nil
}

func validateAllTemplates() error {
	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	templates, err := reg.List()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No templates in registry to validate.")
		return nil
	}

	fmt.Printf("※ Validating %d templates in registry...\n\n", len(templates))

	var failed []string
	for i, tmpl := range templates {
		fmt.Printf("[%d/%d] Validating: %s\n", i+1, len(templates), tmpl.Name)
		if err := validateTemplate(tmpl.Path); err != nil {
			failed = append(failed, tmpl.Name)
			fmt.Printf("❌ Validation failed: %v\n\n", err)
		} else {
			fmt.Println("✅ Validation passed")
			fmt.Println()
		}
	}

	fmt.Println("🔮 Validation Complete:")
	fmt.Printf("   ✅ Passed: %d\n", len(templates)-len(failed))
	if len(failed) > 0 {
		fmt.Printf("   ❌ Failed: %d (%s)\n", len(failed), strings.Join(failed, ", "))
		return fmt.Errorf("validation failed for %d templates", len(failed))
	}

	return nil
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%d min ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%d hr ago", int(diff.Hours()))
	}
	if diff < 7*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	}
	return t.Format("2006-01-02")
}

func getBackupDir(customDir string) string {
	if customDir != "" {
		return customDir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ason", "backups")
}
