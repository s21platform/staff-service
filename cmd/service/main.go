package main

import (
	"log"
	"net"

	"github.com/s21platform/staff-service/internal/config"
	"github.com/s21platform/staff-service/internal/repository/postgres"
	v0 "github.com/s21platform/staff-service/internal/service/v0"
	staffv0 "github.com/s21platform/staff-service/pkg/staff/v0"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.NewConfig()

	dbRepo := postgres.New(cfg)

	service := v0.New(dbRepo)

	lis, err := net.Listen("tcp", ":"+cfg.Service.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	staffv0.RegisterStaffServiceServer(grpcServer, service)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
