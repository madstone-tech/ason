package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	// Test that Execute function exists and can be called
	err := Execute()
	// Since Execute() uses rootCmd.Execute(), it might return errors
	// based on command line arguments, but we're just testing it doesn't panic
	_ = err // We don't assert on the error since it depends on environment
}

func TestRootCmd(t *testing.T) {
	// Test root command properties
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "ason" {
		t.Errorf("rootCmd.Use = %v, want %v", rootCmd.Use, "ason")
	}

	if rootCmd.Short != "※ Shake your projects into being" {
		t.Errorf("rootCmd.Short = %v, want %v", rootCmd.Short, "※ Shake your projects into being")
	}

	if !strings.Contains(rootCmd.Long, "※ Ason - The Sacred Rattle of Code Generation") {
		t.Error("rootCmd.Long should contain the expected description")
	}

	if rootCmd.Version != version {
		t.Errorf("rootCmd.Version = %v, want %v", rootCmd.Version, version)
	}
}

func TestVersion(t *testing.T) {
	// Test version variable
	if version != "0.1.0" {
		t.Errorf("version = %v, want %v", version, "0.1.0")
	}
}

func TestRootCmdHasSubcommands(t *testing.T) {
	// Test that subcommands are added
	commands := rootCmd.Commands()

	expectedCommands := []string{"new", "list", "add", "remove", "validate"}
	commandNames := make([]string, len(commands))

	for i, cmd := range commands {
		commandNames[i] = cmd.Use
	}

	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if strings.HasPrefix(cmd.Use, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand %v not found in command list", expected)
		}
	}
}

func TestRootCmdVersionTemplate(t *testing.T) {
	// Test version template by capturing output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute version command: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "※ Ason v0.1.0") {
		t.Errorf("Version output should contain '※ Ason v0.1.0', got: %v", output)
	}

	// Reset for other tests
	rootCmd.SetOut(nil)
	rootCmd.SetArgs(nil)
}

func TestRootCmdHelp(t *testing.T) {
	// Test help output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute help command: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Ason shakes templates into living projects") {
		t.Errorf("Help output should contain the description, got: %v", output)
	}

	if !strings.Contains(output, "Available Commands:") {
		t.Error("Help output should contain available commands section")
	}

	// Reset for other tests
	rootCmd.SetOut(nil)
	rootCmd.SetArgs(nil)
}
