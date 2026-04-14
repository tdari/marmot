package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marmotdata/marmot/internal/docgen"
)

func main() {
	connDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Connection providers are at: internal/core/connection/providers/<name>/
	// Five levels up reaches the project root, then web/docs/docs
	docsPath := filepath.Join(connDir, "..", "..", "..", "..", "..", "web", "docs", "docs")
	fmt.Printf("Generating docs for connection in: %s\n", connDir)

	if err := docgen.GenerateConnectionDocs(connDir, docsPath); err != nil {
		fmt.Printf("Error generating docs: %v\n", err)
		os.Exit(1)
	}
}
