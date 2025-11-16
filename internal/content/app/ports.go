package app

import (
	"blogThree/internal/content/domain"
	"context"

	"github.com/google/uuid"
)

// -------------------- COMMANDS --------------------

type PostCommandService interface {
	CreatePost(ctx context.Context, authorID uuid.UUID, rawTitle, rawBody string) (*domain.Post, error)
	// später: UpdatePost, DeletePost ...
}

type PostCommandRepository interface {
	Create(ctx context.Context, p *domain.Post) error
}

type UserReader interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

// -------------------- QUERIES --------------------

type PostListFilter struct {
	AuthorID *uuid.UUID
	Limit    int
	Offset   int
}

type PostQueryService interface {
	GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
	ListPosts(ctx context.Context, f PostListFilter) ([]*domain.Post, error)
}

type PostQueryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error)
	List(ctx context.Context, f PostListFilter) ([]*domain.Post, error)
}
