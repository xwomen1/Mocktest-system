package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"upm-simple/pkg/config"
)

func main() {
	fmt.Println("=== Configuration Watch Example ===")

	// Get absolute path to config file
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get working directory:", err)
	}

	// Go up from examples/config-watch to project root
	projectRoot := filepath.Join(baseDir)
	configPath := filepath.Join(projectRoot, "configs", "dev", "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file not found: %s\n", configPath)
		fmt.Println("Please run from project root: .\\scripts\\create-config.ps1")
		waitForExit()
		return
	}

	fmt.Printf("Config file: %s\n", configPath)

	// Load with specific file
	loader := config.NewLoader()
	cfg, err := loader.LoadFromFile(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("\nInitial configuration:\n")
	fmt.Printf("  Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("  NATS: %s\n", cfg.NATS.URL)
	fmt.Printf("  Log level: %s\n", cfg.Logging.Level)

	// Watch for changes
	fmt.Println("\n=== Starting config watcher ===")
	fmt.Println("To test:")
	fmt.Println("1. Open configs/dev/config.yaml in another editor")
	fmt.Println("2. Change server port to 6000")
	fmt.Println("3. Save the file")
	fmt.Println("4. See update here")
	fmt.Println("\nPress Ctrl+C to exit")

	updates := loader.Watch()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	updateCount := 0

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				fmt.Println("\nUpdate channel closed")
				return
			}

			if update.Error != nil {
				fmt.Printf("\nError: %v\n", update.Error)
				continue
			}

			if update.Config != nil {
				updateCount++
				fmt.Printf("\n UPDATE #%d at %s\n", updateCount, time.Now().Format("15:04:05"))
				fmt.Printf("   Server port: %d\n", update.Config.Server.Port)
				fmt.Printf("   Log level: %s\n", update.Config.Logging.Level)
				fmt.Printf("   NATS URL: %s\n", update.Config.NATS.URL)

				// Update current config
				cfg = update.Config
			}

		case sig := <-sigs:
			fmt.Printf("\nReceived %v, exiting...\n", sig)
			return

		case <-ticker.C:
			fmt.Print(".")
		}
	}
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C to exit")
	<-sigs
}
