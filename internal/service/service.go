package service

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
	staff "github.com/s21platform/staff-service/pkg/staff"
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
	staff.UnimplementedStaffServiceServer
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
func (s *StaffService) Get(ctx context.Context, req *staff.GetIn) (*staff.GetOut, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff id")
	}

	staffModel, err := s.repo.StaffGetByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staffModel == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	return &staff.GetOut{
		Staff: convertStaffToProto(staffModel),
	}, nil
}

// CreateStaff создает нового сотрудника
func (s *StaffService) Create(ctx context.Context, req *staff.CreateIn) (*staff.CreateOut, error) {
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

	staffModel := &model.Staff{
		ID:           uuid.New(),
		Login:        req.Login,
		PasswordHash: string(hashedPassword),
		RoleID:       int(req.RoleId),
		Permissions:  permissions,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.StaffCreate(ctx, staffModel); err != nil {
		log.Printf("error: %v", err)
		return nil, status.Error(codes.Internal, "failed to create staff")
	}

	return &staff.CreateOut{
		Staff: convertStaffToProto(staffModel),
	}, nil
}

// UpdateStaff обновляет информацию о сотруднике
func (s *StaffService) Update(ctx context.Context, req *staff.UpdateIn) (*staff.UpdateOut, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff id")
	}

	staffModel, err := s.repo.StaffGetByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staffModel == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	if req.Login != nil {
		staffModel.Login = *req.Login
	}
	if req.RoleId != nil {
		staffModel.RoleID = int(*req.RoleId)
	}
	if req.Permissions != nil {
		staffModel.Permissions.Access = req.Permissions.Access
	}
	staffModel.UpdatedAt = time.Now()

	if err := s.repo.StaffUpdate(ctx, staffModel); err != nil {
		return nil, status.Error(codes.Internal, "failed to update staff")
	}

	return &staff.UpdateOut{
		Staff: convertStaffToProto(staffModel),
	}, nil
}

// DeleteStaff удаляет сотрудника
func (s *StaffService) Delete(ctx context.Context, req *staff.DeleteIn) (*staff.DeleteOut, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff id")
	}

	if err := s.repo.StaffDelete(ctx, id); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete staff")
	}

	return &staff.DeleteOut{}, nil
}

// ListStaff получает список сотрудников с фильтрацией и пагинацией
func (s *StaffService) List(ctx context.Context, req *staff.ListIn) (*staff.ListOut, error) {
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

	protoStaffList := make([]*staff.Staff, len(staffList))
	for i, staff := range staffList {
		protoStaffList[i] = convertStaffToProto(staff)
	}

	// Безопасное вычисление количества страниц
	pageCount := 0
	if pageSize > 0 {
		pageCount = (total + pageSize - 1) / pageSize
	}

	return &staff.ListOut{
		Staff:      protoStaffList,
		TotalCount: int32(total),
		PageCount:  int32(pageCount),
	}, nil
}

// ===== Реализация методов авторизации =====

// Login авторизация сотрудника по логину и паролю
func (s *StaffService) Login(ctx context.Context, req *staff.LoginIn) (*staff.LoginOut, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password are required")
	}

	staffModel, err := s.repo.StaffGetByLogin(ctx, req.Login)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staffModel == nil {
		log.Printf("staff not found")
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staffModel.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("invalid password")
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	session, err := s.createSession(ctx, staffModel.ID)
	if err != nil {
		log.Printf("failed to create session: %v", err)
		return nil, err
	}

	return &staff.LoginOut{
		AccessToken:  session.Token,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt.Unix(),
		Staff:        convertStaffToProto(staffModel),
	}, nil
}

// RefreshToken обновление токена сессии
func (s *StaffService) RefreshToken(ctx context.Context, req *staff.RefreshTokenIn) (*staff.RefreshTokenOut, error) {
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

	staffModel, err := s.repo.StaffGetByID(ctx, session.StaffID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staffModel == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	newSession, err := s.createSession(ctx, staffModel.ID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SessionDelete(ctx, session.Token); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete old session")
	}

	return &staff.RefreshTokenOut{
		AccessToken:  newSession.Token,
		RefreshToken: newSession.RefreshToken,
		ExpiresAt:    newSession.ExpiresAt.Unix(),
	}, nil
}

// Logout выход из системы и завершение сессии
func (s *StaffService) Logout(ctx context.Context, req *staff.LogoutIn) (*staff.LogoutOut, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	if err := s.repo.SessionDelete(ctx, req.AccessToken); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete session")
	}

	return &staff.LogoutOut{
		Success: true,
	}, nil
}

// CheckAuth проверка текущего статуса авторизации
func (s *StaffService) CheckAuth(ctx context.Context, req *staff.CheckAuthIn) (*staff.CheckAuthOut, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	session, err := s.repo.SessionGetByToken(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get session")
	}
	if session == nil {
		return &staff.CheckAuthOut{
			Authorized: false,
		}, nil
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.repo.SessionDelete(ctx, session.Token); err != nil {
			return nil, status.Error(codes.Internal, "failed to delete expired session")
		}
		return &staff.CheckAuthOut{
			Authorized: false,
		}, nil
	}

	staffModel, err := s.repo.StaffGetByID(ctx, session.StaffID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staffModel == nil {
		return &staff.CheckAuthOut{
			Authorized: false,
		}, nil
	}

	return &staff.CheckAuthOut{
		Authorized: true,
		Staff:      convertStaffToProto(staffModel),
	}, nil
}

// ChangePassword изменение пароля авторизованного пользователя
func (s *StaffService) ChangePassword(ctx context.Context, req *staff.ChangePasswordIn) (*staff.ChangePasswordOut, error) {
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

	staffModel, err := s.repo.StaffGetByID(ctx, session.StaffID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get staff")
	}
	if staffModel == nil {
		return nil, status.Error(codes.NotFound, "staff not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staffModel.PasswordHash), []byte(req.OldPassword)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), s.bcryptCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	staffModel.PasswordHash = string(hashedPassword)
	staffModel.UpdatedAt = time.Now()

	if err := s.repo.StaffUpdate(ctx, staffModel); err != nil {
		return nil, status.Error(codes.Internal, "failed to update staff")
	}

	if err := s.repo.SessionDeleteAllForStaff(ctx, staffModel.ID); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete sessions")
	}

	return &staff.ChangePasswordOut{
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
func convertStaffToProto(staffModel *model.Staff) *staff.Staff {
	return &staff.Staff{
		Id:       staffModel.ID.String(),
		Login:    staffModel.Login,
		RoleId:   int32(staffModel.RoleID),
		RoleName: staffModel.RoleName,
		Permissions: &staff.Permissions{
			Access: staffModel.Permissions.Access,
		},
		CreatedAt: staffModel.CreatedAt.Unix(),
		UpdatedAt: staffModel.UpdatedAt.Unix(),
	}
}

// generateToken генерирует случайный токен
func generateToken() string {
	return uuid.New().String()
}
