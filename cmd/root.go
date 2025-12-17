package cmd

import (
	"github.com/spf13/cobra"
)

var version = "0.1.0"

// SetVersionInfo sets the version information (called from main)
func SetVersionInfo(v, c, d, b string) {
	_ = c // commit info - may be used in future
	_ = d // date info - may be used in future
	_ = b // built by info - may be used in future
	version = v
	// Update rootCmd.Version directly to override the initialization value
	rootCmd.Version = v
}

var rootCmd = &cobra.Command{
	Use:   "ason",
	Short: "※ Shake your projects into being",
	Long: `※ Ason - The Sacred Rattle of Code Generation

Ason shakes templates into living projects, catalyzing the transformation
from idea to implementation. Like the sacred rattle used in ceremonies,
this tool invokes change and brings forth new creations.

Named after the ason, the ritual rattle that activates spiritual work
in Haitian Vodou, this tool activates your templates, transforming them
into ready-to-use projects with rhythm and purpose.`,
	Version: version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(`※ Ason {{.Version}}
`)

	// Add commands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(validateCmd)

	// Setup autocompletion
	setupCompletions()
}
