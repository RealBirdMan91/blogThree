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
		AuthorID:  p.AuthorID().String(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}

func toPostModelList(ps []*contentdom.Post) []*model.Post {
	out := make([]*model.Post, 0, len(ps))
	for _, p := range ps {
		out = append(out, toPostModel(p))
	}
	return out
}
