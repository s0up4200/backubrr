package backup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"backubrr/config"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func CreateBackup(config *config.Config, sourceDir string, passphrase string) error {
	var encryptionKey string
	if passphrase == "" {
		encryptionKey = config.EncryptionKey
	} else {
		encryptionKey = passphrase
	}

	//fmt.Println("Encryption key:", encryptionKey)

	// Print source directory being backed up
	color.Blue("Backing up %s...\n", sourceDir)

	// Define archive name
	sourceDirBase := filepath.Base(sourceDir)
	archiveName := fmt.Sprintf("%s_%s.tar.gz", sourceDirBase, time.Now().Format("2006-01-02_15-04-05"))

	// Update the output directory to include the source directory name
	outputDir := filepath.Join(config.OutputDir, sourceDirBase)

	// Create the output directory if it does not exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}
	}

	// Create destination file for writing
	destFile, err := os.Create(filepath.Join(outputDir, archiveName))
	if err != nil {
		return err
	}

	// Create gzip writer
	gzipWriter := gzip.NewWriter(destFile)

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)

	// Create a new spinner with rotating character set
	spin := spinner.New(spinner.CharSets[50], 100*time.Millisecond)

	// Start the spinner
	spin.Start()

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

	// Stop the spinner
	spin.Stop()

	// Close writers and files
	tarWriter.Close()
	gzipWriter.Close()
	destFile.Close()

	// Encrypt archive using GPG, if encryption key is set
	if encryptionKey != "" {
		encryptedArchiveName := fmt.Sprintf("%s.gpg", archiveName)
		cmd := exec.Command("gpg", "--batch", "--symmetric", "--cipher-algo", "AES256", "--passphrase", encryptionKey, "--output", filepath.Join(outputDir, encryptedArchiveName), filepath.Join(outputDir, archiveName))

		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running GPG command:", err)
			fmt.Println("GPG output:", stderr.String())
			return err
		}

		// Remove unencrypted archive file
		if err := os.Remove(filepath.Join(outputDir, archiveName)); err != nil {
			return err
		}

		// Print success message
		message := fmt.Sprintf("Backup created successfully! Archive saved to %s\n\n", filepath.Join(sourceDir, encryptedArchiveName))
		color.Green(message)

	}

	return nil
}
