package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HostsFile  string   `yaml:"hosts_file"`
	Domains    []string `yaml:"domains"`
	Subdomains []string `yaml:"subdomains"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Configuration file does not exist; create a default one
			err = CreateDefaultConfig(path)
			if err != nil {
				return nil, fmt.Errorf("failed to create default configuration file: %v", err)
			}
			// Try reading the file again after creating it
			data, err = ioutil.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read configuration file after creating it: %v", err)
			}
		} else {
			return nil, err
		}
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set default hosts file if not specified
	if config.HostsFile == "" {
		config.HostsFile = "/etc/hosts"
	}

	return &config, nil
}

func CreateDefaultConfig(path string) error {
	defaultConfig := Config{
		Domains: []string{
			"x.com",
			"twitter.com",
			"youtube.com",
			"reddit.com",
		},
		Subdomains: []string{
			"blog.example.com",
			"mail.example.org",
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal default configuration: %v", err)
	}

	// Ensure the configuration directory exists
	err = EnsureConfigDir()
	if err != nil {
		return fmt.Errorf("failed to ensure configuration directory: %v", err)
	}

	// Write the default configuration file
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write default configuration file: %v", err)
	}

	// Change ownership of the configuration file
	uid, gid, err := getOriginalUserIDs()
	if err == nil {
		os.Chown(path, uid, gid)
	}

	fmt.Println("Default configuration file created at", path)
	return nil
}

func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		logError(err)
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return home
}

func EnsureConfigDir() error {
	configDir := filepath.Join(GetHomeDir(), ".dblock")
	return EnsureDir(configDir)
}
