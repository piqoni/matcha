package main

import (
	"fmt"
	"os"
	"os/exec"
)

func buildBinaries() {
	// Set the build targets
	targets := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"darwin", "amd64"},
		{"windows", "amd64"},
	}

	// Create the bin folder if it doesn't already exist
	if _, err := os.Stat("bin"); os.IsNotExist(err) {
		os.Mkdir("bin", 0755)
	}

	// Build the binaries
	for _, target := range targets {
		fmt.Printf("Building for %s/%s...\n", target.os, target.arch)
		cmd := exec.Command("go", "build", "-o", "bin/matcha-"+target.os+"-"+target.arch)
		cmd.Env = append(os.Environ(),
			"GOOS="+target.os,
			"GOARCH="+target.arch,
		)
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
