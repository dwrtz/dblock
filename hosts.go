package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	markerStart = "# BEGIN dblock"
	markerEnd   = "# END dblock"
)

func EnableBlocking(config *Config) error {
	if err := CheckPermissions(config.HostsFile); err != nil {
		return err
	}

	contentBytes, err := ioutil.ReadFile(config.HostsFile)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	// Create backup
	backupDir := filepath.Join(GetHomeDir(), ".dblock", "backups")
	err = EnsureDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to ensure backup directory: %v", err)
	}

	backupPath := filepath.Join(backupDir, "hosts.bak")
	if err := ioutil.WriteFile(backupPath, contentBytes, 0644); err != nil {
		return err
	}

	// Change ownership of the backup file
	uid, gid, err := getOriginalUserIDs()
	if err == nil {
		os.Chown(backupPath, uid, gid)
	}

	newContent := RemoveExistingBlocks(content)
	blockEntries := GenerateBlockEntries(config)

	// Trim trailing spaces and newlines
	newContent = strings.TrimRight(newContent, " \n")

	// Decide whether to add a newline before the dblock section
	if len(newContent) > 0 {
		newContent += "\n\n"
	} else {
		newContent += "\n"
	}

	// Append the dblock section
	newContent += fmt.Sprintf("%s\n%s\n%s", markerStart, blockEntries, markerEnd)

	// Ensure the file ends with a single newline
	newContent += "\n"

	if err := ioutil.WriteFile(config.HostsFile, []byte(newContent), 0644); err != nil {
		return err
	}

	return nil
}

func DisableBlocking(config *Config) error {
	if err := CheckPermissions(config.HostsFile); err != nil {
		return err
	}

	contentBytes, err := ioutil.ReadFile(config.HostsFile)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	// Create backup
	backupDir := filepath.Join(GetHomeDir(), ".dblock", "backups")
	err = EnsureDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to ensure backup directory: %v", err)
	}

	backupPath := filepath.Join(backupDir, "hosts.bak")
	if err := ioutil.WriteFile(backupPath, contentBytes, 0644); err != nil {
		return err
	}

	// Change ownership of the backup file
	uid, gid, err := getOriginalUserIDs()
	if err == nil {
		os.Chown(backupPath, uid, gid)
	}

	newContent := RemoveExistingBlocks(content)

	// Trim trailing spaces and newlines
	newContent = strings.TrimRight(newContent, " \n")

	// Ensure the file ends with a single newline
	if len(newContent) > 0 {
		newContent += "\n"
	}

	if err := ioutil.WriteFile(config.HostsFile, []byte(newContent), 0644); err != nil {
		return err
	}

	return nil
}

func GetStatus(config *Config) (string, error) {
	content, err := ioutil.ReadFile(config.HostsFile)
	if err != nil {
		return "", err
	}

	if strings.Contains(string(content), markerStart) {
		return "Enabled", nil
	}
	return "Disabled", nil
}

func GenerateBlockEntries(config *Config) string {
	var entries []string
	for _, domain := range config.Domains {
		entries = append(entries, generateDomainEntries(domain)...)
	}
	for _, subdomain := range config.Subdomains {
		entries = append(entries, generateEntry(subdomain))
	}
	return strings.Join(entries, "\n")
}

func generateDomainEntries(domain string) []string {
	subdomains := []string{domain, "www." + domain}
	var entries []string
	for _, subdomain := range subdomains {
		entries = append(entries, generateEntry(subdomain))
	}
	return entries
}

func generateEntry(domain string) string {
	return fmt.Sprintf("127.0.0.1\t%s\n::1\t%s", domain, domain)
}

func RemoveExistingBlocks(content string) string {
	lines := strings.Split(content, "\n")
	var newLines []string
	skip := false
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.Contains(trimmedLine, markerStart) {
			skip = true
			// Remove any preceding blank lines
			for len(newLines) > 0 && strings.TrimSpace(newLines[len(newLines)-1]) == "" {
				newLines = newLines[:len(newLines)-1]
			}
			continue
		}
		if strings.Contains(trimmedLine, markerEnd) {
			skip = false
			continue
		}
		if !skip {
			newLines = append(newLines, line)
		}
	}
	// Remove any trailing blank lines
	for len(newLines) > 0 && strings.TrimSpace(newLines[len(newLines)-1]) == "" {
		newLines = newLines[:len(newLines)-1]
	}
	return strings.Join(newLines, "\n")
}

func CheckPermissions(filePath string) error {
	if unix.Access(filePath, unix.W_OK) != nil {
		return fmt.Errorf("insufficient permissions to modify %s", filePath)
	}
	return nil
}
