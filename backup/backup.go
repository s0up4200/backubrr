package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"backubrr/config"

	"github.com/fatih/color"
)

func CreateBackup(config *config.Config, sourceDir string) error {
	// Print source directory being backed up
	color.Blue("Backing up %s...\n", sourceDir)

	// Define archive name
	archiveName := fmt.Sprintf("%s_%s.tar.gz", filepath.Base(sourceDir), time.Now().Format("2006-01-02_15-04-05"))

	// Create destination file for writing
	destFile, err := os.Create(filepath.Join(config.OutputDir, archiveName))
	if err != nil {
		return err
	}

	// Create gzip writer
	gzipWriter := gzip.NewWriter(destFile)

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)

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

	// Close writers and files
	tarWriter.Close()
	gzipWriter.Close()
	destFile.Close()

	// Print success message
	message := fmt.Sprintf("Backup created successfully! Archive saved to %s\n\n", filepath.Join(config.OutputDir, archiveName))
	color.Green(message)

	return nil
}
