package main

import (
	"fmt"
	"os"
	"strconv"
)

func PrintHelp() {
	fmt.Print(`Usage: dblock [command] [options]

Commands:
  enable      Enable blocking of specified domains
  disable     Disable blocking
  status      Get current blocking status
  list        List configured domains and subdomains
  help        Display this help message

Options:
  -t, --timeout    Timeout in minutes after which the operation is reversed
  -c, --config     Path to the configuration file

Examples:
  sudo dblock enable
  sudo dblock enable -t 60
  sudo dblock disable
  dblock status
  dblock list
`)
}

func ListDomains(config *Config) {
	fmt.Println("Configured domains:")
	for _, domain := range config.Domains {
		fmt.Println(" -", domain)
	}
	fmt.Println("Configured subdomains:")
	for _, subdomain := range config.Subdomains {
		fmt.Println(" -", subdomain)
	}
}

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}

	uid, gid, err := getOriginalUserIDs()
	if err == nil {
		os.Chown(dirName, uid, gid)
	}

	return nil
}

func getOriginalUserIDs() (int, int, error) {
	uidStr := os.Getenv("SUDO_UID")
	gidStr := os.Getenv("SUDO_GID")

	if uidStr == "" || gidStr == "" {
		return -1, -1, fmt.Errorf("SUDO_UID or SUDO_GID not set")
	}

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid SUDO_UID: %v", err)
	}

	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid SUDO_GID: %v", err)
	}

	return uid, gid, nil
}
