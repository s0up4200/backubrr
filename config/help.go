package config

import (
	"fmt"
)

func PrintHelp() {
	fmt.Printf(`
A command-line tool for backing up files and directories.

Usage:
backubrr [flags]

Flags:
--config string    path to config file (default "config.yaml")        Specifies the path to the configuration file. Optional.
--passphrase       encryption key passphrase                          Specifies the passphrase to use for GPG encryption. Optional.
-h, --help         show this message                                  Displays this help message.
version            show version information                           Displays version, commit, and date information.

Configuration options:
  source_dirs        A list of directories to back up.
  output_dir         The directory where backup files are saved.
  encryption_key     Set a passphrase for encryption (Can also be called directly as --passphrase)
  retention_days     The number of days to retain backup files.
  interval           Run every X hours.
  discord            Send notifications to Discord after a backup run.

`)
}
