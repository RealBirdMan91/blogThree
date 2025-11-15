package app

import (
	"blogThree/internal/content/domain"
	"context"

	"github.com/google/uuid"
)

type PostRepository interface {
	Create(ctx context.Context, p *domain.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Post, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Post, error)
	ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*domain.Post, error)
}

// Cross-BC Reader (User-BC)
type UserReader interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
