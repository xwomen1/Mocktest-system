package main

import (
    "fmt"
    "time"

    "github.com/nats-io/nats.go"
)

func main() {
    fmt.Println("Testing NATS...")

    nc, err := nats.Connect("nats://localhost:4222")
    if err != nil {
        fmt.Println("NATS error:", err)
        return
    }
    defer nc.Close()

    fmt.Println("Connected to NATS")

    // Subscribe
    sub, _ := nc.Subscribe("test.upmp", func(m *nats.Msg) {
        fmt.Printf("Received: %s\n", string(m.Data))
    })
    defer sub.Unsubscribe()

    // Publish
    nc.Publish("test.upmp", []byte("Hello from UPM"))
    nc.Flush()

    time.Sleep(100 * time.Millisecond)
    fmt.Println("NATS test completed")
}
