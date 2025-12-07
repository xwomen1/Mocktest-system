package main

import (
	"fmt"
	"log"
	"os"

	"upm-simple/pkg/config"
)

func main() {
	fmt.Println("=== Environment Variable Override Example ===")

	os.Setenv("UPM_SERVER_PORT", "6000")
	os.Setenv("UPM_LOGGING_LEVEL", "warn")
	os.Setenv("UPM_NATS_URL", "nats://nats-test:4222")

	fmt.Println("Environment variables set:")
	fmt.Println("  UPM_SERVER_PORT=6000")
	fmt.Println("  UPM_LOGGING_LEVEL=warn")
	fmt.Println("  UPM_NATS_URL=nats://nats-test:4222")

	fmt.Println("\nLoading configuration...")
	loader := config.NewLoader()
	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("\nLoaded configuration (with env overrides):")
	fmt.Printf("  Server Port: %d (env overrode default)\n", cfg.Server.Port)
	fmt.Printf("  NATS URL: %s (env overrode default)\n", cfg.NATS.URL)
	fmt.Printf("  Log Level: %s (env overrode default)\n", cfg.Logging.Level)

	fmt.Println("\n? Environment variable override test completed!")
}
