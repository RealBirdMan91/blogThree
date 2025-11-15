package postgres

import (
	app "blogThree/internal/content/app"
	"blogThree/internal/content/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PostgresPostRepo struct{ db *sql.DB }

func NewPostgresPostRepo(db *sql.DB) *PostgresPostRepo { return &PostgresPostRepo{db: db} }

var _ app.PostRepository = (*PostgresPostRepo)(nil)

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

func (r *PostgresPostRepo) List(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	const q = `
        SELECT id, title, body, author_id, created_at, updated_at
        FROM posts
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
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

func (r *PostgresPostRepo) ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*domain.Post, error) {
	const q = `
		SELECT id, title, body, author_id, created_at, updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, q, authorID, limit, offset)
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
			aID       uuid.UUID
			createdAt time.Time
			updatedAt time.Time
		)
		if err := rows.Scan(&id, &titleStr, &bodyStr, &aID, &createdAt, &updatedAt); err != nil {
			return nil, app.NewPostsListFailed(err)
		}

		p, err := r.hydratePost(id, titleStr, bodyStr, aID, createdAt, updatedAt)
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
