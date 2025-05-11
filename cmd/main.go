package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/C0d3-5t3w/go-ssg/inc/cli"
	"github.com/C0d3-5t3w/go-ssg/inc/gen"
)

const outputDir = "output" // Should match the one in gen/gen.go
const serverPort = "8080"

func main() {
	// Ensure the output directory exists, or create it if gen.GEN() might not run
	// if it's empty or due to other logic.
	// gen.GEN() already creates it, but this is a good practice for the server.
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory for server: %v", err)
		}
	}

	if err := gen.GEN(); err != nil {
		log.Fatalf("Failed to generate site: %v", err)
	}
	cli.CLI() // Assuming CLI() might have other setup or logging

	fmt.Printf("Successfully generated site in '%s' directory.\n", outputDir)
	fmt.Printf("Serving files from '%s' on http://localhost:%s\n", outputDir, serverPort)

	fs := http.FileServer(http.Dir(outputDir))
	http.Handle("/", fs)

	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
