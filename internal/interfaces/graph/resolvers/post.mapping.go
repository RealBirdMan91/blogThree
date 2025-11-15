package resolvers

import (
	contentdom "blogThree/internal/content/domain"
	"blogThree/internal/interfaces/graph/model"
)

func toPostModel(p *contentdom.Post) *model.Post {
	return &model.Post{
		ID:        p.ID().String(),
		Title:     p.Title().String(),
		Body:      p.Body().String(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}
