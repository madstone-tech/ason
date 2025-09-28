package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/madstone-tech/ason/internal/registry"
	"github.com/spf13/cobra"
)

// completeTemplateNames provides completion for template names from the registry
func completeTemplateNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	reg, err := registry.NewRegistry()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	templates, err := reg.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for _, template := range templates {
		if strings.HasPrefix(template.Name, toComplete) {
			completions = append(completions, template.Name)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeTemplateNamesOrPaths provides completion for template names or local paths
func completeTemplateNamesOrPaths(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var completions []string

	// First, try to complete template names from registry
	reg, err := registry.NewRegistry()
	if err == nil {
		templates, err := reg.List()
		if err == nil {
			for _, template := range templates {
				if strings.HasPrefix(template.Name, toComplete) {
					completions = append(completions, template.Name)
				}
			}
		}
	}

	// If we have registry completions, don't show files
	if len(completions) > 0 {
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	// Otherwise, allow directory completion for local templates
	return nil, cobra.ShellCompDirectiveFilterDirs
}

// completeOutputPaths provides completion for output directory paths
func completeOutputPaths(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveFilterDirs
}

// completeTemplatePaths provides completion for template file or directory paths
func completeTemplatePaths(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Look for template files and directories
	var completions []string

	// If toComplete is empty or doesn't contain a path separator, look in current directory
	searchDir := "."
	prefix := toComplete

	if strings.Contains(toComplete, string(filepath.Separator)) {
		searchDir = filepath.Dir(toComplete)
		prefix = filepath.Base(toComplete)
	}

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, prefix) {
			fullPath := name
			if searchDir != "." {
				fullPath = filepath.Join(searchDir, name)
			}

			if entry.IsDir() {
				// Add trailing slash for directories
				completions = append(completions, fullPath+string(filepath.Separator))
			} else if isTemplateFile(name) {
				completions = append(completions, fullPath)
			}
		}
	}

	return completions, cobra.ShellCompDirectiveDefault
}

// isTemplateFile checks if a file could be a template file
func isTemplateFile(filename string) bool {
	templateExts := []string{".yaml", ".yml", ".json", ".tmpl", ".template"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, templateExt := range templateExts {
		if ext == templateExt {
			return true
		}
	}

	return false
}

// completeVariableKeys provides completion for variable keys
func completeVariableKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Common variable names for completion
	commonVars := []string{
		"name=",
		"version=",
		"author=",
		"description=",
		"license=",
		"project_name=",
		"package_name=",
		"app_name=",
		"module_name=",
		"namespace=",
		"use_docker=",
		"use_tests=",
		"use_ci=",
	}

	var completions []string
	for _, varName := range commonVars {
		if strings.HasPrefix(varName, toComplete) {
			completions = append(completions, varName)
		}
	}

	return completions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

// completeAddCommand provides completion for the add command
func completeAddCommand(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// First argument is template name (no completion needed, it's user-defined)
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Second argument is template path - complete directories
	if len(args) == 1 {
		return nil, cobra.ShellCompDirectiveFilterDirs
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

// setupCompletions configures completion for all commands
func setupCompletions() {
	// Set up completion for the new command
	newCmd.ValidArgsFunction = completeTemplateNamesOrPaths

	// Set up completion for the add command
	addCmd.ValidArgsFunction = completeAddCommand

	// Set up completion for remove command
	removeCmd.ValidArgsFunction = completeTemplateNames

	// Set up completion for validate command
	validateCmd.ValidArgsFunction = completeTemplatePaths

	// Add completion for flags
	newCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	})

	newCmd.RegisterFlagCompletionFunc("var", completeVariableKeys)
}
