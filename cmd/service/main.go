package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/s21platform/staff-service/internal/config"
	"github.com/s21platform/staff-service/internal/middleware"
	"github.com/s21platform/staff-service/internal/repository/postgres"
	"github.com/s21platform/staff-service/internal/service"
	staff "github.com/s21platform/staff-service/pkg/staff"
)

func main() {
	cfg := config.NewConfig()

	dbRepo := postgres.New(cfg)

	srv := service.New(dbRepo)

	lis, err := net.Listen("tcp", ":"+cfg.Service.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем интерсептор для проверки ролей
	authInterceptor := middleware.NewAuthInterceptor(dbRepo)

	// Создаем gRPC сервер с интерсептором
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	staff.RegisterStaffServiceServer(grpcServer, srv)

	log.Printf("Starting staff service on port %s", cfg.Service.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
