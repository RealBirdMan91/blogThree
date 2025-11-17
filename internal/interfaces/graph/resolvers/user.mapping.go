package resolvers

import (
	"blogThree/internal/interfaces/graph/model"
	"blogThree/internal/user/domain"
)

func toUserModel(u *domain.User) *model.User {
	return &model.User{
		ID:        u.ID().String(),
		Email:     u.Email().String(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}

func toUserModelList(us []*domain.User) []*model.User {
	out := make([]*model.User, 0, len(us))
	for _, u := range us {
		out = append(out, toUserModel(u))
	}
	return out
}
