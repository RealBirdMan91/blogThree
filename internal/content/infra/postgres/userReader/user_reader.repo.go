package userReader

import (
	"blogThree/internal/content/app"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Reader struct{ db *sql.DB }

func New(db *sql.DB) *Reader { return &Reader{db: db} }

var _ app.UserReader = (*Reader)(nil)

func (r *Reader) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	const q = `SELECT 1 FROM users WHERE id = $1`
	var one int
	err := r.db.QueryRowContext(ctx, q, id).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
