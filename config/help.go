package config

import (
	"fmt"
	"os"
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

func PrintHelp() {
	fmt.Printf(`
A command-line tool for backing up files and directories.

Usage:
backubrr [flags]

Flags:
--config string    path to config file (default "config.yaml")        Specifies the path to the configuration file.
-h, --help         show this message                                  Displays this help message.
version      show version information                           Displays version, commit, and date information.

Configuration options:
  source_dirs        A list of directories to back up.
  output_dir         The directory where backup files are saved.
  retention_days     The number of days to retain backup files.
  interval           Run every X hours.
  discord            Send notifications to Discord after a backup run.

`)
}

func PrintVersion() {
	fmt.Printf("Backubrr %s %s %s\n", version, commit[:7], date)
}
