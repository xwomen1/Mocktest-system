package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"

    pb "upm-simple/internal"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type server struct {
    pb.UnimplementedServiceRegistryServer
    mu       sync.RWMutex
    services map[string]*pb.Service
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    if req.Service == nil {
        return nil, status.Error(codes.InvalidArgument, "service is required")
    }

    id := fmt.Sprintf("%s-%s-%d", req.Service.Name, req.Service.Host, req.Service.Port)
    req.Service.Id = id

    s.mu.Lock()
    if s.services == nil {
        s.services = make(map[string]*pb.Service)
    }
    s.services[id] = req.Service
    s.mu.Unlock()

    fmt.Printf("Registered: %s\n", req.Service.Name)

    return &pb.RegisterResponse{
        Response: &pb.CommonResponse{
            Success: true,
            Message: "Registered",
        },
        ServiceId: id,
    }, nil
}

func (s *server) Discover(ctx context.Context, req *pb.DiscoverRequest) (*pb.DiscoverResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var found []*pb.Service
    for _, svc := range s.services {
        if req.ServiceName == "" || svc.Name == req.ServiceName {
            found = append(found, svc)
        }
    }

    return &pb.DiscoverResponse{
        Response: &pb.CommonResponse{
            Success: true,
            Message: fmt.Sprintf("Found %d services", len(found)),
        },
        Services: found,
    }, nil
}

func main() {
    fmt.Println("Starting Service Registry on :50051")

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }

    s := grpc.NewServer()
    pb.RegisterServiceRegistryServer(s, &server{})

    if err := s.Serve(lis); err != nil {
        log.Fatal(err)
    }
}
