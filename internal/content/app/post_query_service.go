package app

import (
	"blogThree/internal/content/domain"
	"context"

	"github.com/google/uuid"
)

type queryService struct {
	repo PostQueryRepository
}

func NewQueryService(repo PostQueryRepository) PostQueryService {
	return &queryService{repo: repo}
}

func (s *queryService) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *queryService) ListPosts(ctx context.Context, f PostListFilter) ([]*domain.Post, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return s.repo.List(ctx, f)
}
