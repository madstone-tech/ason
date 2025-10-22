package main

import (
	"log"
	"os"

	"github.com/madstone-tech/ason/cmd"
)

// Version information, set via ldflags during build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "source"
)

func main() {
	// Set version info in cmd package
	cmd.SetVersionInfo(version, commit, date, builtBy)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
