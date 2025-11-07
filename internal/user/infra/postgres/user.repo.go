package postgres

import (
	"blogThree/internal/user/app"
	"blogThree/internal/user/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PostgresUserRepo struct{ db *sql.DB }

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo { return &PostgresUserRepo{db: db} }

var _ app.UserRepository = (*PostgresUserRepo)(nil)

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.Email().String(),
		user.PasswordHash().String(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	return err
}

func (r *PostgresUserRepo) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	const q = `SELECT 1 FROM users WHERE lower(email)=lower($1) LIMIT 1`
	var one int
	err := r.db.QueryRowContext(ctx, q, email.String()).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
		LIMIT 1
	`

	var (
		id        uuid.UUID
		emailStr  string
		hashStr   string
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, q, email.String()).
		Scan(&id, &emailStr, &hashStr, &createdAt, &updatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	emailVO, err := domain.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	hashVO, err := domain.NewPasswordHash(hashStr)
	if err != nil {
		return nil, err
	}

	user, err := domain.RehydrateUser(id, emailVO, hashVO, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var (
		emailStr  string
		hashStr   string
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, q, id).
		Scan(&id, &emailStr, &hashStr, &createdAt, &updatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	emailVO, err := domain.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	hashVO, err := domain.NewPasswordHash(hashStr)
	if err != nil {
		return nil, err
	}

	user, err := domain.RehydrateUser(id, emailVO, hashVO, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}
