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
	var backupMessages []string
	var passphrase string
	flag.StringVar(&configFilePath, "config", "config.yaml", "path to config file")
	flag.StringVar(&passphrase, "passphrase", "", "encryption key passphrase")
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

	// Check if encryption key is set in config
	if config.EncryptionKey != "" && passphrase != "" {
		errorMessage := "Encryption key is already set in config. Please remove the --passphrase argument or unset the encryption key in the config file."
		color.HiRed(errorMessage)
		os.Exit(1)
	}

	// Create destination directory if it doesn't exist
	err = os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Create backup for each source directory
		for _, sourceDir := range config.SourceDirs {
			err := backup.CreateBackup(config, sourceDir, passphrase)
			if err != nil {
				log.Printf("Error creating backup of %s: %s\n", sourceDir, err)
				continue
			}

			archiveName := fmt.Sprintf("%s_%s.tar.gz", filepath.Base(sourceDir), time.Now().Format("2006-01-02_15-04-05"))
			backupDir := filepath.Join(config.OutputDir, filepath.Base(sourceDir))
			var backupMessage string
			if config.EncryptionKey != "" {
				encryptedArchiveName := archiveName + ".enc"
				backupMessage = fmt.Sprintf("Backup of **`%s`** created successfully! Encrypted archive saved to **`%s`**\n", sourceDir, filepath.Join(backupDir, encryptedArchiveName))
			} else {
				backupMessage = fmt.Sprintf("Backup of **`%s`** created successfully! Archive saved to **`%s`**\n", sourceDir, filepath.Join(backupDir, archiveName))
			}
			backupMessage = strings.Replace(backupMessage, os.Getenv("HOME"), "~", -1)
			backupMessages = append(backupMessages, backupMessage)
		}

		// Combine backup messages into a single message
		backupMessage := strings.Join(backupMessages, "")

		// Calculate next backup time
		var nextBackupTime time.Time
		if config.Interval > 0 {
			duration := time.Duration(config.Interval) * time.Hour
			nextBackupTime = time.Now().Add(duration)
		}

		// Create next backup message
		var nextBackupMessage string
		if !nextBackupTime.IsZero() {
			nextBackupMessage = "\nNext backup will run at **`" + nextBackupTime.Format("2006-01-02 15:04:05") + "`**\n"
		}

		// Combine backup message and next scheduled backup message
		combinedMessage := backupMessage + nextBackupMessage

		// Send combined message to Discord
		if config.DiscordWebhookURL != "" {
			if err := notifications.SendToDiscordWebhook(config.DiscordWebhookURL, []string{combinedMessage}); err != nil {
				fmt.Println("Error sending message to Discord:", err)
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
