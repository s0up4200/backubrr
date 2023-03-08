package main

import (
	"backubrr/backup"
	"backubrr/cleaner"
	"backubrr/config"
	"backubrr/notifications"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	version string = "unknown"
	commit  string = "unknown"
	date    string = "unknown"
)

func init() {
	if v := os.Getenv("BACKUBRR_VERSION"); v != "" {
		version = v
	}
	if c := os.Getenv("BACKUBRR_COMMIT"); c != "" {
		commit = c
	}
	if d := os.Getenv("BACKUBRR_DATE"); d != "" {
		date = d
	}
}

func PrintVersion() {
	fmt.Printf("backubrr v%s %s %s\n", version, commit[:7], date)
}

func main() {
	// Parse command-line arguments
	flag.Usage = config.PrintHelp
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "config.yaml", "path to config file")
	flag.Parse()

	if len(os.Args) == 2 && (os.Args[1] == "version" || os.Args[1] == "-v" || os.Args[1] == "--version") {
		PrintVersion()
		return
	}

	// Load configuration from file
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create destination directory if it doesn't exist
	os.MkdirAll(config.OutputDir, 0755)

	var backupMessages []string

	for {
		// Create backup for each source directory
		for _, sourceDir := range config.SourceDirs {
			err := backup.CreateBackup(config, sourceDir)
			if err != nil {
				fmt.Printf("Error creating backup of %s: %s\n", sourceDir, err)
				continue
			}

			// Replace home directory path with ~ in backup message
			backupMessage := "Backup of " + filepath.Base(sourceDir) + " created successfully! Archive saved to " + filepath.Join(config.OutputDir, filepath.Base(sourceDir)+"_"+time.Now().Format("2006-01-02_15-04-05")+".tar.gz") + "\n"
			backupMessage = strings.Replace(backupMessage, os.Getenv("HOME"), "~", -1)
			backupMessages = append(backupMessages, backupMessage)
		}

		// Calculate next backup time
		var nextBackupTime time.Time
		if config.Interval > 0 {
			duration := time.Duration(config.Interval) * time.Hour
			nextBackupTime = time.Now().Add(duration)
		}

		if config.DiscordWebhookURL != "" {
			var messages []string

			if len(backupMessages) > 0 {
				combinedBackupMessage := strings.Join(backupMessages, "")
				messages = append(messages, combinedBackupMessage)
			}

			if !nextBackupTime.IsZero() {
				nextBackupMessage := "Next backup will run at **`" + nextBackupTime.Format("2006-01-02 15:04:05") + "`**"
				messages = append(messages, nextBackupMessage)
			}

			if len(messages) > 0 {
				if err := notifications.SendToDiscordWebhook(config.DiscordWebhookURL, messages); err != nil {
					fmt.Println("Error sending message to Discord:", err)
				}
			}
		}

		// Clean up old backups
		if err := cleaner.Cleaner(configFilePath); err != nil {
			fmt.Println("Error cleaning up old backups:", err)
		}

		// Sleep until the next backup time, if configured
		if config.Interval > 0 {
			if nextBackupTime.IsZero() {
				duration := time.Duration(config.Interval) * time.Hour
				nextBackupTime = time.Now().Add(duration)
			}
			color.Cyan("\nNext backup will run at %s\n", nextBackupTime.Format("2006-01-02 15:04:05"))
			time.Sleep(time.Until(nextBackupTime))
		} else {
			break
		}

	}
}
