package main

import (
	"archive/tar"
	"backubrr/cleaner"
	"backubrr/config"
	"backubrr/notifications"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
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
  interval           Run every X hours.
  discord            Send notifications to Discord after a backup run.
  `)
}

func main() {
	// Parse command-line arguments
	flag.Usage = printHelp
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.yaml", "path to config file")
	flag.Parse()

	var backupMessage string

	// Load configuration from file
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create destination directory if it doesn't exist
	os.MkdirAll(config.OutputDir, 0755)

	for {
		for _, sourceDir := range config.SourceDirs {
			// Print source directory being backed up
			color.Blue("Backing up %s...\n", sourceDir)

			// Define archive name
			archiveName := fmt.Sprintf("%s_%s.tar.gz", filepath.Base(sourceDir), time.Now().Format("2006-01-02_15-04-05"))

			// Create destination file for writing
			destFile, err := os.Create(filepath.Join(config.OutputDir, archiveName))
			if err != nil {
				log.Fatal(err)
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

			// Append success message to backup message variable
			backupMessage += fmt.Sprintf("Backup of **`%s`** created successfully! Archive saved to **`%s`**\n", filepath.Base(sourceDir), filepath.Join(config.OutputDir, archiveName))

			if err != nil {
				log.Fatal(err)
			}
		}

		// Send backup message to Discord webhook
		if config.DiscordWebhookURL != "" {
			if err := notifications.SendToDiscordWebhook(config.DiscordWebhookURL, backupMessage); err != nil {
				fmt.Println("Error sending message to Discord:", err)
			} else {
				fmt.Println("Message sent to Discord successfully!")
			}
		}

		// Clean up old backups
		if err := cleaner.Cleaner(configFilePath); err != nil {
			log.Fatal("Error cleaning up old backups: ", err)
		}

		// Sleep until the next backup time, if configured
		if config.Interval > 0 {
			duration := time.Duration(config.Interval) * time.Hour
			nextBackupTime := time.Now().Add(duration)
			color.Cyan("Next backup will run at %s\n", nextBackupTime.Format("2006-01-02 15:04:05"))
			time.Sleep(duration)
		} else {
			break
		}
	}
}
