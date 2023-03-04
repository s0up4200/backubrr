package config

import (
	"errors"
	"io/ioutil"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

// Config stores the configuration for the backup program.
type Config struct {
	SourceDirs    []string `yaml:"source_dirs"`
	OutputDir     string   `yaml:"output_dir"`
	RetentionDays int      `yaml:"retention_days"`
	Interval      int      `yaml:"interval"`
}

// LoadConfig loads the backup configuration from a YAML file.
func LoadConfig(filePath string) (*Config, error) {
	// Read config file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse YAML data
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if config.Interval < 0 {
		return nil, errors.New(color.HiRedString("error: Interval must be a positive number"))
	}

	return config, nil
}
