package main

import (
	"os"

	"github.com/uclatall/ckhub/cmd/ckhub/app"
)

var version = "unknown"

func main() {
	cmd := app.NewCommand(version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
