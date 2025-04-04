package middleware

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	RoleOwner  = 1
	RoleAdmin  = 2
	RoleStaff  = 3
	RoleViewer = 4
)

// RolePermissions определяет разрешения для каждого метода gRPC
var RolePermissions = map[string][]int{
	"/staff.StaffService/Create": {RoleOwner},
	"/staff.StaffService/Update": {RoleOwner, RoleAdmin},
	"/staff.StaffService/Delete": {RoleOwner},
	"/staff.StaffService/List":   {RoleOwner, RoleAdmin, RoleStaff, RoleViewer},
	"/staff.StaffService/Get":    {RoleOwner, RoleAdmin, RoleStaff, RoleViewer},
}

type AuthInterceptor struct {
	sessionManager SessionManager
}

type SessionManager interface {
	GetStaffRoleByToken(ctx context.Context, token string) (int, error)
}

func NewAuthInterceptor(sessionManager SessionManager) *AuthInterceptor {
	return &AuthInterceptor{
		sessionManager: sessionManager,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Пропускаем методы авторизации
		if info.FullMethod == "/staff.StaffService/Login" ||
			info.FullMethod == "/staff.StaffService/RefreshToken" {
			return handler(ctx, req)
		}

		// Получаем токен из метаданных
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		token := values[0]
		log.Printf("token: %s", token)
		// Получаем роль пользователя
		roleID, err := i.sessionManager.GetStaffRoleByToken(ctx, token)
		if err != nil {
			log.Printf("failed to get staff role by token: %v", err)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Проверяем права доступа
		allowedRoles, exists := RolePermissions[info.FullMethod]
		if !exists {
			// Если метод не указан в разрешениях, разрешаем только Owner
			allowedRoles = []int{RoleOwner}
		}

		hasPermission := false
		for _, role := range allowedRoles {
			if roleID == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("role %d does not have permission to access %s", roleID, info.FullMethod))
		}

		return handler(ctx, req)
	}
}
