package main

import (
	"log"
	"os"

	"github.com/madstone-tech/ason/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
