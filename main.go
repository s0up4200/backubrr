package main

import (
	"archive/tar"
	"backubrr/cleaner"
	"backubrr/config"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func printHelp() {
	fmt.Println(`
Backubrr

A command-line tool for backing up files and directories.

Usage:
  backubrr [flags]

Flags:
  --config string    path to config file (default "config.yaml")	Specifies the path to the configuration file.
  -h, --help         show this message					Displays this help message.

Configuration options:
  source_dirs        A list of directories to back up.
  output_dir         The directory where backup files are saved.
  retention_days     The number of days to retain backup files.
  `)
}

func main() {
	// Parse command-line arguments
	flag.Usage = printHelp
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.yaml", "path to config file")
	flag.Parse()

	// Load configuration from file
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create destination directory if it doesn't exist
	os.MkdirAll(config.OutputDir, 0755)

	for _, sourceDir := range config.SourceDirs {
		// Print source directory being backed up
		fmt.Printf("Backing up %s...\n", sourceDir)
		// Define archive name
		archiveName := fmt.Sprintf("%s_%s.tar.gz", filepath.Base(sourceDir), time.Now().Format("2006-01-02_15-04-05"))

		// Create destination file for writing
		destFile, err := os.Create(filepath.Join(config.OutputDir, archiveName))
		if err != nil {
			log.Fatal(err)
		}
		defer destFile.Close()

		// Create gzip writer
		gzipWriter := gzip.NewWriter(destFile)
		defer gzipWriter.Close()

		// Create tar writer
		tarWriter := tar.NewWriter(gzipWriter)
		defer tarWriter.Close()

		// Walk through source directory recursively
		filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || filepath.Base(path)[0] == '.' {
				return err
			}

			// Create tar header
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			header.Name = path[len(sourceDir)+1:]

			// Write header to tar archive
			if err = tarWriter.WriteHeader(header); err != nil {
				return err
			}

			// Open source file for reading
			sourceFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			// Copy source file contents to tar archive
			if _, err = io.Copy(tarWriter, sourceFile); err != nil {
				return err
			}

			return nil
		})

		// Print success message
		fmt.Printf("Backup created successfully! Archive saved to %s\n", filepath.Join(config.OutputDir, archiveName))

		// Clean up old backups
		if err := cleaner.Cleaner(configFilePath); err != nil {
			log.Fatal("Error cleaning up old backups: ", err)
		}
	}
}
