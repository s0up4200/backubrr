package cleaner

import (
	"fmt"
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

	// Walk through backup directory
	cutoffTime := time.Now().Add(-7 * 24 * time.Hour)
	deletedBackups := false
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only consider regular files with .tar.gz extension
		if !info.Mode().IsRegular() || filepath.Ext(path) != ".tar.gz" {
			return nil
		}

		// Check if file is older than 7 days
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

	if deletedBackups {
		fmt.Println("Old backups deleted successfully.")
	} else {
		fmt.Println("No old backups found. Cleanup not needed.")
	}

	return nil
}
