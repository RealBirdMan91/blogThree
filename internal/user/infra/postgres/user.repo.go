package postgres

import (
	"blogThree/internal/user/app"
	"blogThree/internal/user/domain"
	"context"
	"database/sql"
	"errors"
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
