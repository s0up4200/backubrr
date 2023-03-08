package cleaner

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func Cleaner(configFilePath string) error {
	// Load configuration from file
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Get output directory from configuration file
	outputDir := viper.GetString("output_dir")
	if outputDir == "" {
		return fmt.Errorf("output_dir not set in configuration file")
	}

	// Check if output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("output_dir does not exist: %s", outputDir)
	}

	// Get retention days from configuration file
	retentionDays := viper.GetInt("retention_days")
	if retentionDays == 0 {
		retentionDays = 7 // default retention period is 7 days
	}

	// Calculate cutoff time based on retention period
	cutoffTime := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)
	deletedBackups := false

	// Walk through output directory and subdirectories
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only consider regular files with .tar.gz extension
		if !info.Mode().IsRegular() || filepath.Ext(path) != ".tar.gz" {
			return nil
		}

		// Check if file is older than retention period
		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(path); err != nil {
				return err
			}
			deletedBackups = true
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Walk through output directory and subdirectories again to delete empty directories
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only consider directories that are empty
		if !info.IsDir() || !isEmptyDir(path) {
			return nil
		}

		if err := os.Remove(path); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.SetFlags(0)

	if deletedBackups {
		log.Printf("Old backups deleted successfully.")
	} else {
		log.Printf("No old backups found. Cleanup not needed.")
	}

	return nil
}

// Returns true if the directory is empty (contains no files or subdirectories)
func isEmptyDir(path string) bool {
	dir, err := os.Open(path)
	if err != nil {
		return false
	}
	defer dir.Close()

	_, err = dir.Readdir(1)
	if err == nil {
		// Directory is not empty
		return false
	}
	if err == io.EOF {
		// Directory is empty
		return true
	}
	// Error occurred while reading directory
	return false
}
