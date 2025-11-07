package resolvers

import (
	authApp "blogThree/internal/auth/app"
	userApp "blogThree/internal/user/app"
)

//go:generate go run github.com/99designs/gqlgen generate
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserSvc userApp.UserService
	AuthSvc authApp.AuthService
}
