package postgres

import (
	"blogThree/internal/auth/app"
	"blogThree/internal/auth/domain"
	"context"
	"database/sql"
)

type RefreshTokenRepo struct {
	db *sql.DB
}

var _ app.RefreshTokenRepository = (*RefreshTokenRepo)(nil)

func NewRefreshTokenRepo(db *sql.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, t domain.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (id, user_id, token_hash, revoked, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, q,
		t.ID, t.UserID, t.TokenHash, t.Revoked, t.ExpiresAt, t.CreatedAt,
	)
	return err
}

func (r *RefreshTokenRepo) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, revoked, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
		LIMIT 1
	`

	var t domain.RefreshToken

	err := r.db.QueryRowContext(ctx, q, hash).Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.Revoked,
		&t.ExpiresAt,
		&t.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *RefreshTokenRepo) Update(ctx context.Context, t *domain.RefreshToken) error {
	const q = `
		UPDATE refresh_tokens
		SET revoked = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, q, t.Revoked, t.ID)
	return err
}
