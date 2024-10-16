package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	// Initialize logging
	SetupLogging()

	// Parse command-line arguments
	if len(os.Args) < 2 {
		PrintHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	// Default configuration path
	configPath := filepath.Join(GetHomeDir(), ".dblock", "default.yaml")

	// Command-line flags
	fs := flag.NewFlagSet(command, flag.ExitOnError)
	timeout := fs.Int("t", 0, "Timeout in minutes after which the operation is reversed")
	configFile := fs.String("c", configPath, "Path to the configuration file")

	// Parse flags
	if err := fs.Parse(os.Args[2:]); err != nil {
		logError(err)
		os.Exit(1)
	}

	// Load configuration
	config, err := LoadConfig(*configFile)
	if err != nil {
		logError(err)
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	// Check for root privileges when modifying the hosts file
	if command == "enable" || command == "disable" {
		if !isRoot() {
			fmt.Println("Error: Insufficient permissions. Please run the command with 'sudo'.")
			os.Exit(1)
		}
	}

	switch command {
	case "enable":
		if err := EnableBlocking(config); err != nil {
			logError(err)
			fmt.Println("Error enabling blocking:", err)
			os.Exit(1)
		}
		fmt.Println("Blocking enabled.")
		if *timeout > 0 {
			fmt.Printf("Blocking will be disabled in %d minutes.\n", *timeout)
			// Create a channel to signal when the goroutine has finished
			done := make(chan struct{})
			// Start the goroutine
			go func() {
				time.Sleep(time.Duration(*timeout) * time.Minute)
				if err := DisableBlocking(config); err != nil {
					logError(err)
					fmt.Println("Error disabling blocking after timeout:", err)
				} else {
					fmt.Println("Blocking disabled after timeout.")
				}
				// Signal that the goroutine has completed
				close(done)
			}()
			// Wait for either the goroutine to finish or an interrupt signal
			waitForCompletion(done)
		}
	case "disable":
		if err := DisableBlocking(config); err != nil {
			logError(err)
			fmt.Println("Error disabling blocking:", err)
			os.Exit(1)
		}
		fmt.Println("Blocking disabled.")
		if *timeout > 0 {
			fmt.Printf("Blocking will be re-enabled in %d minutes.\n", *timeout)
			// Create a channel to signal when the goroutine has finished
			done := make(chan struct{})
			// Start the goroutine
			go func() {
				time.Sleep(time.Duration(*timeout) * time.Minute)
				if err := EnableBlocking(config); err != nil {
					logError(err)
					fmt.Println("Error enabling blocking after timeout:", err)
				} else {
					fmt.Println("Blocking re-enabled after timeout.")
				}
				// Signal that the goroutine has completed
				close(done)
			}()
			// Wait for either the goroutine to finish or an interrupt signal
			waitForCompletion(done)
		}
	case "status":
		status, err := GetStatus(config)
		if err != nil {
			logError(err)
			fmt.Println("Error getting status:", err)
			os.Exit(1)
		}
		fmt.Println("Blocking status:", status)
	case "list":
		ListDomains(config)
	case "help":
		PrintHelp()
	default:
		fmt.Println("Unknown command:", command)
		PrintHelp()
		os.Exit(1)
	}

}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		logError(err)
		fmt.Println("Error determining current user:", err)
		os.Exit(1)
	}
	return currentUser.Uid == "0"
}

func waitForCompletion(done chan struct{}) {
	// Set up signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-done:
		// Goroutine completed
		os.Exit(0)
	case sig := <-sigs:
		fmt.Printf("\nReceived signal: %s. Exiting.\n", sig)
		os.Exit(0)
	}
}
