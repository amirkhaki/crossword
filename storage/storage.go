package storage

import (
	"context"

	"github.com/amirkhaki/crossword/config"
	"github.com/amirkhaki/crossword/user"
)

type User interface {
	AddUser(context.Context, user.User) error
	GetUser(context.Context, user.User, func(user.User, user.User) bool) (user.User, error)
	// TODO UpdateUser(context.Context, user.User) error
	// TODO DeleteUser(context.Context, user.User) error
}

type Group interface {
	AddGroup(context.Context, user.Group) error
	GetGroup(context.Context, user.Group, func(user.Group, user.Group) bool) (user.Group, error)
	// TODO UpdateGroup(context.Context, user.Group) error
	// TODO DeleteGroup(context.Context, user.Group) error
}

type Storage interface {
	User
	Group
}

func NewStorage(cfg config.Config) (Storage, error) {
	return nil, nil
}
