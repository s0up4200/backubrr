package main

import (
    "log"
    "os"
    "os/user"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type Config struct {
    OutputDir  string   `yaml:"output_dir"`
    SourceDirs []string `yaml:"source_dirs"`
}

func LoadConfig(path string) (*Config, error) {
    // Load configuration from file
    var config Config
    configFile, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer configFile.Close()
    decoder := yaml.NewDecoder(configFile)
    if err = decoder.Decode(&config); err != nil {
        return nil, err
    }

    // Check that output_dir is specified
    if config.OutputDir == "" {
        currentUser, err := user.Current()
        if err != nil {
            return nil, err
        }
        config.OutputDir = filepath.Join(currentUser.HomeDir, "backups")
        log.Printf("output_dir not specified in configuration file, using default: %s\n", config.OutputDir)
    }

    return &config, nil
}
