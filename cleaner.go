package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func cleaner() {
	// Get output directory from environment variable
	outputDir := os.Getenv("BACKUP_OUTPUT_DIR")
	if outputDir == "" {
		log.Fatal("BACKUP_OUTPUT_DIR environment variable is not set")
	}

	// Get current time
	now := time.Now()

	// Walk through output directory
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if file is a backup file older than 7 days
		if !info.IsDir() && filepath.Ext(path) == ".tar.gz" && now.Sub(info.ModTime()) > 7*24*time.Hour {
			// Delete file
			if err := os.Remove(path); err != nil {
				return err
			}

			// Print success message
			fmt.Printf("Deleted backup file %s\n", path)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
