package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/s21platform/staff-service/internal/config"
	"github.com/s21platform/staff-service/internal/model"
)

// Ошибки репозитория
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

// Repo реализует интерфейс DbRepo для работы с PostgreSQL
type Repo struct {
	db *sqlx.DB
}

// New создает новый экземпляр репозитория
func New(cfg *config.Config) *Repo {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &Repo{
		db: db,
	}
}

// ===== Методы для работы со Staff =====

// StaffGetByID получает информацию о сотруднике по ID
func (r *Repo) StaffGetByID(ctx context.Context, id uuid.UUID) (*model.Staff, error) {
	query, args, err := sq.
		Select("s.id", "s.login", "s.password_hash", "s.role_id", "r.name as role_name",
			"s.permissions", "s.created_at", "s.updated_at").
		From("staff s").
		LeftJoin("roles r ON s.role_id = r.id").
		Where(sq.Eq{"s.id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	staff := &model.Staff{}
	err = r.db.GetContext(ctx, staff, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	return staff, nil
}

// StaffGetByLogin получает информацию о сотруднике по логину
func (r *Repo) StaffGetByLogin(ctx context.Context, login string) (*model.Staff, error) {
	query, args, err := sq.
		Select("s.id", "s.login", "s.password_hash", "s.role_id", "r.name as role_name",
			"s.permissions", "s.created_at", "s.updated_at").
		From("staff s").
		LeftJoin("roles r ON s.role_id = r.id").
		Where(sq.Eq{"s.login": login}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	staff := &model.Staff{}
	err = r.db.GetContext(ctx, staff, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	return staff, nil
}

// StaffCreate создает нового сотрудника
func (r *Repo) StaffCreate(ctx context.Context, staff *model.Staff) error {
	query, args, err := sq.
		Insert("staff").
		Columns("id", "login", "password_hash", "role_id", "created_at", "updated_at").
		Values(staff.ID, staff.Login, staff.PasswordHash, staff.RoleID, staff.CreatedAt, staff.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	log.Printf("query: %s, args: %v", query, args)
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create staff: %w", err)
	}

	return nil
}

// StaffUpdate обновляет информацию о сотруднике
func (r *Repo) StaffUpdate(ctx context.Context, staff *model.Staff) error {
	query, args, err := sq.
		Update("staff").
		Set("login", staff.Login).
		Set("password_hash", staff.PasswordHash).
		Set("role_id", staff.RoleID).
		Set("updated_at", staff.UpdatedAt).
		Where(sq.Eq{"id": staff.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// StaffDelete удаляет сотрудника
func (r *Repo) StaffDelete(ctx context.Context, id uuid.UUID) error {
	query, args, err := sq.
		Delete("staff").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// StaffList получает список сотрудников с фильтрацией и пагинацией
func (r *Repo) StaffList(ctx context.Context, filter *model.StaffFilter) ([]*model.Staff, int, error) {
	// Построение базового запроса
	baseQuery := sq.
		Select("s.id", "s.login", "s.password_hash", "s.role_id", "r.name as role_name",
			"s.permissions", "s.created_at", "s.updated_at").
		From("staff s").
		LeftJoin("roles r ON s.role_id = r.id").
		PlaceholderFormat(sq.Dollar)

	// Добавление условий фильтрации
	if filter.SearchTerm != "" {
		baseQuery = baseQuery.Where(sq.Like{"s.login": "%" + filter.SearchTerm + "%"})
	}
	if filter.RoleID != 0 {
		baseQuery = baseQuery.Where(sq.Eq{"s.role_id": filter.RoleID})
	}

	// Получение общего количества записей
	countQuery := baseQuery.Column("COUNT(*) OVER()")
	query, args, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query: %w", err)
	}

	type staffWithCount struct {
		model.Staff
		Count int `db:"count"`
	}

	var staffList []staffWithCount
	err = r.db.SelectContext(ctx, &staffList, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get staff list: %w", err)
	}

	total := 0
	if len(staffList) > 0 {
		total = staffList[0].Count
	}

	// Добавление пагинации
	query, args, err = baseQuery.
		Offset(uint64((filter.Page - 1) * filter.PageSize)).
		Limit(uint64(filter.PageSize)).
		ToSql()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to build paginated query: %w", err)
	}

	var result []*model.Staff
	err = r.db.SelectContext(ctx, &result, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get staff list: %w", err)
	}

	return result, total, nil
}

// ===== Методы для работы с Session =====

// SessionCreate создает новую сессию
func (r *Repo) SessionCreate(ctx context.Context, session *model.Session) error {
	query, args, err := sq.
		Insert("sessions").
		Columns("id", "staff_id", "token", "refresh_token", "expires_at",
			"created_at", "last_activity_at").
		Values(session.ID, session.StaffID, session.Token, session.RefreshToken,
			session.ExpiresAt, session.CreatedAt, session.LastActivityAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("failed to create session (in repo): %v", err)
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// SessionGetByToken получает сессию по токену
func (r *Repo) SessionGetByToken(ctx context.Context, token string) (*model.Session, error) {
	query, args, err := sq.
		Select("id", "staff_id", "token", "refresh_token", "expires_at",
			"created_at", "last_activity_at").
		From("sessions").
		Where(sq.Eq{"token": token}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	session := &model.Session{}
	err = r.db.GetContext(ctx, session, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// SessionGetByRefreshToken получает сессию по refresh токену
func (r *Repo) SessionGetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	query, args, err := sq.
		Select("id", "staff_id", "token", "refresh_token", "expires_at",
			"created_at", "last_activity_at").
		From("sessions").
		Where(sq.Eq{"refresh_token": refreshToken}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	session := &model.Session{}
	err = r.db.GetContext(ctx, session, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// SessionDelete удаляет сессию по токену
func (r *Repo) SessionDelete(ctx context.Context, token string) error {
	query, args, err := sq.
		Delete("sessions").
		Where(sq.Eq{"token": token}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// SessionDeleteAllForStaff удаляет все сессии сотрудника
func (r *Repo) SessionDeleteAllForStaff(ctx context.Context, staffID uuid.UUID) error {
	query, args, err := sq.
		Delete("sessions").
		Where(sq.Eq{"staff_id": staffID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete sessions: %w", err)
	}

	return nil
}

// SessionUpdateTokens обновляет токены сессии
func (r *Repo) SessionUpdateTokens(ctx context.Context, session *model.Session) error {
	query, args, err := sq.
		Update("sessions").
		Set("token", session.Token).
		Set("refresh_token", session.RefreshToken).
		Set("expires_at", session.ExpiresAt).
		Set("last_activity_at", session.LastActivityAt).
		Where(sq.Eq{"id": session.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// ===== Методы для работы с Role =====

// RoleGetByID получает роль по ID
func (r *Repo) RoleGetByID(ctx context.Context, id int) (*model.Role, error) {
	query, args, err := sq.
		Select("id", "name").
		From("roles").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	role := &model.Role{}
	err = r.db.GetContext(ctx, role, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// RoleList получает список всех ролей
func (r *Repo) RoleList(ctx context.Context) ([]*model.Role, error) {
	query, args, err := sq.
		Select("id", "name").
		From("roles").
		OrderBy("id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var roles []*model.Role
	err = r.db.SelectContext(ctx, &roles, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return roles, nil
}

// GetStaffRoleByToken получает роль сотрудника по токену сессии
func (r *Repo) GetStaffRoleByToken(ctx context.Context, token string) (int, error) {
	query := `
		SELECT s.role_id
		FROM staff s
		JOIN sessions sess ON sess.staff_id = s.id
		WHERE sess.token = $1 AND sess.expires_at > NOW()
	`

	var roleID int
	err := r.db.GetContext(ctx, &roleID, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("failed to get staff role: %w", err)
	}

	return roleID, nil
}
