package main

import (
	"fmt"
	"log"

	"upm-simple/pkg/config"
)

func main() {
	fmt.Println("=== Configuration Basic Example ===")

	fmt.Println("\n1. Loading default configuration...")
	loader := config.NewLoader()
	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	printConfig(cfg)

	fmt.Println("\n2. Getting individual values...")
	fmt.Printf("  Server Port: %d\n", loader.GetInt("server.port"))
	fmt.Printf("  NATS URL: %s\n", loader.GetString("nats.url"))
	fmt.Printf("  Log Level: %s\n", loader.GetString("logging.level"))

	fmt.Println("\n Configuration test completed!")
}

func printConfig(cfg *config.Config) {
	fmt.Println("Configuration loaded:")
	fmt.Printf("  Environment: %s\n", cfg.Environment)
	fmt.Printf("  Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("  NATS: %s\n", cfg.NATS.URL)
	fmt.Printf("  Logging: Level=%s\n", cfg.Logging.Level)
}
