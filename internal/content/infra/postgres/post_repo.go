package postgres

import (
	app "blogThree/internal/content/app"
	"blogThree/internal/content/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PostgresPostRepo struct{ db *sql.DB }

func NewPostgresPostRepo(db *sql.DB) *PostgresPostRepo { return &PostgresPostRepo{db: db} }

// Implementiert beide Ports:
var _ app.PostCommandRepository = (*PostgresPostRepo)(nil)
var _ app.PostQueryRepository = (*PostgresPostRepo)(nil)

// -------------------- COMMANDS --------------------
func (r *PostgresPostRepo) Create(ctx context.Context, p *domain.Post) error {
	const q = `
		INSERT INTO posts (id, title, body, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, q,
		p.ID(),
		p.Title().String(),
		p.Body().String(),
		p.AuthorID(),
		p.CreatedAt(),
		p.UpdatedAt(),
	)
	if err != nil {
		return app.NewPostPersistFailed(err)
	}
	return nil
}

// -------------------- QUERIES --------------------

func (r *PostgresPostRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	const q = `
		SELECT id, title, body, author_id, created_at, updated_at
		FROM posts WHERE id = $1 LIMIT 1
	`
	var (
		postID    uuid.UUID
		titleStr  string
		bodyStr   string
		authorID  uuid.UUID
		createdAt time.Time
		updatedAt time.Time
	)
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&postID, &titleStr, &bodyStr, &authorID, &createdAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrPostNotFound
	}
	if err != nil {
		return nil, app.NewPostSelectFailed(err)
	}

	return r.hydratePost(postID, titleStr, bodyStr, authorID, createdAt, updatedAt)
}

func (r *PostgresPostRepo) List(ctx context.Context, f app.PostListFilter) ([]*domain.Post, error) {
	base := `
        SELECT id, title, body, author_id, created_at, updated_at
        FROM posts
    `
	where := ""
	args := []any{}
	i := 1

	if f.AuthorID != nil {
		where = fmt.Sprintf(" WHERE author_id = $%d", i)
		args = append(args, *f.AuthorID)
		i++
	}

	limitOffset := fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, f.Limit, f.Offset)

	q := base + where + limitOffset

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, app.NewPostsListFailed(err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var (
			id        uuid.UUID
			titleStr  string
			bodyStr   string
			authorID  uuid.UUID
			createdAt time.Time
			updatedAt time.Time
		)
		if err := rows.Scan(&id, &titleStr, &bodyStr, &authorID, &createdAt, &updatedAt); err != nil {
			return nil, app.NewPostsListFailed(err)
		}

		p, err := r.hydratePost(id, titleStr, bodyStr, authorID, createdAt, updatedAt)
		if err != nil {
			return nil, app.NewPostsListFailed(err)
		}

		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, app.NewPostsListFailed(err)
	}
	return posts, nil
}

func (r *PostgresPostRepo) hydratePost(
	id uuid.UUID,
	titleStr, bodyStr string,
	authorID uuid.UUID,
	createdAt, updatedAt time.Time,
) (*domain.Post, error) {
	title, err := domain.NewTitle(titleStr)
	if err != nil {
		return nil, app.NewPostSelectFailed(err)
	}
	body, err := domain.NewBody(bodyStr)
	if err != nil {
		return nil, app.NewPostSelectFailed(err)
	}
	post, err := domain.RehydratePost(id, title, body, authorID, createdAt, updatedAt)
	if err != nil {
		return nil, app.NewPostSelectFailed(err)
	}
	return post, nil
}
