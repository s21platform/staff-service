package v0

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s21platform/staff-service/internal/model"
	staffv0 "github.com/s21platform/staff-service/pkg/staff/v0"
)

var (
	// ErrNotFound возвращается, когда запрашиваемый объект не найден
	ErrNotFound = errors.New("not found")
	// ErrInvalidCredentials возвращается при неверных учетных данных
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrSessionExpired возвращается, когда срок действия сессии истек
	ErrSessionExpired = errors.New("session expired")
	// ErrInvalidToken возвращается при недействительном токене
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidInput возвращается, когда входные данные некорректны
	ErrInvalidInput = errors.New("invalid input")
)

// Константы для настройки сервиса
const (
	// DefaultAccessTokenTTL срок действия access токена по умолчанию (1 час)
	DefaultAccessTokenTTL = 1 * time.Hour
	// DefaultRefreshTokenTTL срок действия refresh токена по умолчанию (7 дней)
	DefaultRefreshTokenTTL = 7 * 24 * time.Hour
	// DefaultBcryptCost стоимость хеширования bcrypt
	DefaultBcryptCost = 10
)

// StaffService реализует gRPC API для управления персоналом
type StaffService struct {
	staffv0.UnimplementedStaffServiceServer
	repo DbRepo

	// Настройки сервиса
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	bcryptCost      int
}

// NewStaffService создает новый экземпляр сервиса
func New(
	repo DbRepo,
	opts ...ServiceOption,
) *StaffService {
	s := &StaffService{
		repo:            repo,
		accessTokenTTL:  DefaultAccessTokenTTL,
		refreshTokenTTL: DefaultRefreshTokenTTL,
		bcryptCost:      DefaultBcryptCost,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// ServiceOption функциональная опция для настройки сервиса
type ServiceOption func(*StaffService)

// WithAccessTokenTTL устанавливает срок действия access токена
func WithAccessTokenTTL(ttl time.Duration) ServiceOption {
	return func(s *StaffService) {
		s.accessTokenTTL = ttl
	}
}

// WithRefreshTokenTTL устанавливает срок действия refresh токена
func WithRefreshTokenTTL(ttl time.Duration) ServiceOption {
	return func(s *StaffService) {
		s.refreshTokenTTL = ttl
	}
}

// WithBcryptCost устанавливает стоимость хеширования bcrypt
func WithBcryptCost(cost int) ServiceOption {
	return func(s *StaffService) {
		s.bcryptCost = cost
	}
}

// ===== Реализация методов управления персоналом =====

// GetStaff получает информацию о сотруднике по ID
func (s *StaffService) GetStaff(ctx context.Context, req *staffv0.GetStaffRequest) (*staffv0.GetStaffResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff id")
	}

	staff, err := s.repo.StaffGetByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staff == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	return &staffv0.GetStaffResponse{
		Staff: convertStaffToProto(staff),
	}, nil
}

// CreateStaff создает нового сотрудника
func (s *StaffService) CreateStaff(ctx context.Context, req *staffv0.CreateStaffRequest) (*staffv0.CreateStaffResponse, error) {
	if req.Login == "" || req.Password == "" || req.RoleId == 0 {
		log.Printf("invalid input: %v", req)
		return nil, status.Error(codes.InvalidArgument, "login, password and role_id are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	permissions := model.Permissions{}
	if req.Permissions != nil {
		permissions.Access = req.Permissions.Access
	}

	staff := &model.Staff{
		ID:           uuid.New(),
		Login:        req.Login,
		PasswordHash: string(hashedPassword),
		RoleID:       int(req.RoleId),
		Permissions:  permissions,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.StaffCreate(ctx, staff); err != nil {
		log.Printf("error: %v", err)
		return nil, status.Error(codes.Internal, "failed to create staff")
	}

	return &staffv0.CreateStaffResponse{
		Staff: convertStaffToProto(staff),
	}, nil
}

// UpdateStaff обновляет информацию о сотруднике
func (s *StaffService) UpdateStaff(ctx context.Context, req *staffv0.UpdateStaffRequest) (*staffv0.UpdateStaffResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff id")
	}

	staff, err := s.repo.StaffGetByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staff == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	if req.Login != nil {
		staff.Login = *req.Login
	}
	if req.RoleId != nil {
		staff.RoleID = int(*req.RoleId)
	}
	if req.Permissions != nil {
		staff.Permissions.Access = req.Permissions.Access
	}
	staff.UpdatedAt = time.Now()

	if err := s.repo.StaffUpdate(ctx, staff); err != nil {
		return nil, status.Error(codes.Internal, "failed to update staff")
	}

	return &staffv0.UpdateStaffResponse{
		Staff: convertStaffToProto(staff),
	}, nil
}

// DeleteStaff удаляет сотрудника
func (s *StaffService) DeleteStaff(ctx context.Context, req *staffv0.DeleteStaffRequest) (*staffv0.DeleteStaffResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff id")
	}

	if err := s.repo.StaffDelete(ctx, id); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete staff")
	}

	return &staffv0.DeleteStaffResponse{}, nil
}

// ListStaff получает список сотрудников с фильтрацией и пагинацией
func (s *StaffService) ListStaff(ctx context.Context, req *staffv0.ListStaffRequest) (*staffv0.ListStaffResponse, error) {
	// Установка значений по умолчанию для пагинации
	page := int(req.Page)
	if page < 1 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize < 1 {
		pageSize = 10 // значение по умолчанию
	}

	filter := &model.StaffFilter{
		Page:     page,
		PageSize: pageSize,
	}
	if req.SearchTerm != nil {
		filter.SearchTerm = *req.SearchTerm
	}
	if req.RoleId != nil {
		filter.RoleID = int(*req.RoleId)
	}

	staffList, total, err := s.repo.StaffList(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list staff")
	}

	protoStaffList := make([]*staffv0.Staff, len(staffList))
	for i, staff := range staffList {
		protoStaffList[i] = convertStaffToProto(staff)
	}

	// Безопасное вычисление количества страниц
	pageCount := 0
	if pageSize > 0 {
		pageCount = (total + pageSize - 1) / pageSize
	}

	return &staffv0.ListStaffResponse{
		Staff:      protoStaffList,
		TotalCount: int32(total),
		PageCount:  int32(pageCount),
	}, nil
}

// ===== Реализация методов авторизации =====

// Login авторизация сотрудника по логину и паролю
func (s *StaffService) Login(ctx context.Context, req *staffv0.LoginRequest) (*staffv0.LoginResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password are required")
	}

	staff, err := s.repo.StaffGetByLogin(ctx, req.Login)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staff == nil {
		log.Printf("staff not found")
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("invalid password")
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	session, err := s.createSession(ctx, staff.ID)
	if err != nil {
		log.Printf("failed to create session: %v", err)
		return nil, err
	}

	return &staffv0.LoginResponse{
		AccessToken:  session.Token,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt.Unix(),
		Staff:        convertStaffToProto(staff),
	}, nil
}

// RefreshToken обновление токена сессии
func (s *StaffService) RefreshToken(ctx context.Context, req *staffv0.RefreshTokenRequest) (*staffv0.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	session, err := s.repo.SessionGetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get session")
	}
	if session == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.repo.SessionDelete(ctx, session.Token); err != nil {
			return nil, status.Error(codes.Internal, "failed to delete expired session")
		}
		return nil, status.Error(codes.Unauthenticated, "refresh token expired")
	}

	staff, err := s.repo.StaffGetByID(ctx, session.StaffID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staff == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	newSession, err := s.createSession(ctx, staff.ID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SessionDelete(ctx, session.Token); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete old session")
	}

	return &staffv0.RefreshTokenResponse{
		AccessToken:  newSession.Token,
		RefreshToken: newSession.RefreshToken,
		ExpiresAt:    newSession.ExpiresAt.Unix(),
	}, nil
}

// Logout выход из системы и завершение сессии
func (s *StaffService) Logout(ctx context.Context, req *staffv0.LogoutRequest) (*staffv0.LogoutResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	if err := s.repo.SessionDelete(ctx, req.AccessToken); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete session")
	}

	return &staffv0.LogoutResponse{
		Success: true,
	}, nil
}

// CheckAuth проверка текущего статуса авторизации
func (s *StaffService) CheckAuth(ctx context.Context, req *staffv0.CheckAuthRequest) (*staffv0.CheckAuthResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	session, err := s.repo.SessionGetByToken(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get session")
	}
	if session == nil {
		return &staffv0.CheckAuthResponse{
			Authorized: false,
		}, nil
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.repo.SessionDelete(ctx, session.Token); err != nil {
			return nil, status.Error(codes.Internal, "failed to delete expired session")
		}
		return &staffv0.CheckAuthResponse{
			Authorized: false,
		}, nil
	}

	staff, err := s.repo.StaffGetByID(ctx, session.StaffID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staff == nil {
		return &staffv0.CheckAuthResponse{
			Authorized: false,
		}, nil
	}

	return &staffv0.CheckAuthResponse{
		Authorized: true,
		Staff:      convertStaffToProto(staff),
	}, nil
}

// ChangePassword изменение пароля авторизованного пользователя
func (s *StaffService) ChangePassword(ctx context.Context, req *staffv0.ChangePasswordRequest) (*staffv0.ChangePasswordResponse, error) {
	if req.OldPassword == "" || req.NewPassword == "" || req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "old_password, new_password and access_token are required")
	}

	session, err := s.repo.SessionGetByToken(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get session")
	}
	if session == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.repo.SessionDelete(ctx, session.Token); err != nil {
			return nil, status.Error(codes.Internal, "failed to delete expired session")
		}
		return nil, status.Error(codes.Unauthenticated, "access token expired")
	}

	staff, err := s.repo.StaffGetByID(ctx, session.StaffID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staff == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.OldPassword)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), s.bcryptCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	staff.PasswordHash = string(hashedPassword)
	staff.UpdatedAt = time.Now()

	if err := s.repo.StaffUpdate(ctx, staff); err != nil {
		return nil, status.Error(codes.Internal, "failed to update staff")
	}

	if err := s.repo.SessionDeleteAllForStaff(ctx, staff.ID); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete sessions")
	}

	return &staffv0.ChangePasswordResponse{
		Success: true,
	}, nil
}

// ===== Вспомогательные методы =====

// createSession создает новую сессию для сотрудника
func (s *StaffService) createSession(ctx context.Context, staffID uuid.UUID) (*model.Session, error) {
	session := &model.Session{
		ID:             uuid.New(),
		StaffID:        staffID,
		Token:          generateToken(),
		RefreshToken:   generateToken(),
		ExpiresAt:      time.Now().Add(s.accessTokenTTL),
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}

	if err := s.repo.SessionCreate(ctx, session); err != nil {
		return nil, status.Error(codes.Internal, "failed to create session")
	}

	return session, nil
}

// convertStaffToProto преобразует модель Staff в proto-сообщение
func convertStaffToProto(staff *model.Staff) *staffv0.Staff {
	return &staffv0.Staff{
		Id:       staff.ID.String(),
		Login:    staff.Login,
		RoleId:   int32(staff.RoleID),
		RoleName: staff.RoleName,
		Permissions: &staffv0.Permissions{
			Access: staff.Permissions.Access,
		},
		CreatedAt: staff.CreatedAt.Unix(),
		UpdatedAt: staff.UpdatedAt.Unix(),
	}
}

func wrappedInt32(value int32) *int32 {
	return &value
}

// generateToken генерирует случайный токен
func generateToken() string {
	return uuid.New().String()
}
