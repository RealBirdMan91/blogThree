package app

import (
	"blogThree/internal/content/domain"
	"context"

	"github.com/google/uuid"
)

type commandService struct {
	repo       PostCommandRepository
	userReader UserReader
}

func NewCommandService(repo PostCommandRepository, ur UserReader) PostCommandService {
	return &commandService{
		repo:       repo,
		userReader: ur,
	}
}

func (s *commandService) CreatePost(ctx context.Context, authorID uuid.UUID, rawTitle, rawBody string) (*domain.Post, error) {
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
		return nil, NewInvalidTitleError(err.Error())
	}

	body, err := domain.NewBody(rawBody)
	if err != nil {
		return nil, NewInvalidBodyError(err.Error())
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
