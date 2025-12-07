package main

import (
    "context"
    "fmt"
    "log"
    "time"

    pb "upm-simple/internal"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    fmt.Println("Testing Service Registry...")

    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
        grpc.WithTimeout(5*time.Second))
    if err != nil {
        log.Fatal("Cannot connect:", err)
    }
    defer conn.Close()

    client := pb.NewServiceRegistryClient(conn)

    // Register service
    fmt.Println("Registering service...")
    resp, err := client.Register(context.Background(), &pb.RegisterRequest{
        Service: &pb.Service{
            Name: "mock-engine",
            Host: "localhost",
            Port: 8080,
        },
    })
    if err != nil {
        log.Fatal("Register failed:", err)
    }
    fmt.Printf("Registered with ID: %s\n", resp.ServiceId)

    // Discover service
    fmt.Println("Discovering services...")
    discResp, err := client.Discover(context.Background(), &pb.DiscoverRequest{
        ServiceName: "mock-engine",
    })
    if err != nil {
        log.Fatal("Discover failed:", err)
    }
    fmt.Printf("Found %d services\n", len(discResp.Services))

    fmt.Println("Test PASSED!")
}
