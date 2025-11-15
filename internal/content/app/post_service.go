package app

import (
	"blogThree/internal/content/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo       PostRepository
	userReader UserReader
}

type PostService interface {
	CreatePost(ctx context.Context, authorID uuid.UUID, rawTitle, rawBody string) (*domain.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
	ListPosts(ctx context.Context, limit, offset int) ([]*domain.Post, error)
	ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*domain.Post, error)
}

func NewService(repo PostRepository, ur UserReader) PostService {
	return &Service{repo: repo, userReader: ur}
}

func (s *Service) CreatePost(ctx context.Context, authorID uuid.UUID, rawTitle, rawBody string) (*domain.Post, error) {
	// 1) Author existiert?
	ok, err := s.userReader.Exists(ctx, authorID)
	if err != nil {
		// Upstream Problem (User-BC/DB/Netz) -> Technical
		return nil, NewAuthorCheckFailed(err)
	}
	if !ok {
		return nil, ErrAuthorNotFound
	}

	// 2) VOs bauen (Domain-Fehler -> Validation in App mappen)
	title, err := domain.NewTitle(rawTitle)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyTitle) {
			return nil, NewInvalidTitleError("empty")
		}
		if errors.Is(err, domain.ErrTitleTooLong) {
			return nil, NewInvalidTitleError("too long")
		}
		return nil, NewInvalidTitleError("") // fallback
	}

	body, err := domain.NewBody(rawBody)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyBody) {
			return nil, NewInvalidBodyError()
		}
		return nil, NewInvalidBodyError()
	}

	// 3) Aggregate bauen
	post, err := domain.NewPost(authorID, title, body)
	if err != nil {
		return nil, NewUnknownPostError(err)
	}

	// 4) Persistieren
	if err := s.repo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err // Repo liefert ErrPostNotFound oder Technical
	}
	return post, nil
}

func (s *Service) ListPosts(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	posts, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *Service) ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*domain.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListByAuthor(ctx, authorID, limit, offset)

}
