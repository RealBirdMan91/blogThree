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
