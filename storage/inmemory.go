package storage

import (
	"fmt"
  "context"

	"github.com/amirkhaki/crossword/user"
)

type UserNotFoundError error
type GroupNotFoundError error

type inmemory struct {
	users []user.User
  groups []user.Group
}

func (im *inmemory) AddUser(ctx context.Context, u user.User) error {
	_, err := im.GetUser(ctx, u, func(u1, u2 user.User) bool {
		if u1.Username == u2.Username {
			return true
		}
		return false
	})

	err, ok := err.(UserNotFoundError)

	if err == nil {
		return fmt.Errorf("Adduser: user with given username (%s) exists!", u.Username)
	} else if !ok {
		return fmt.Errorf("Adduser: error while checking uniqueness: %w", err)
	}

	im.users = append(im.users, u)
	return nil
}

func (im *inmemory) GetUser(ctx context.Context, u user.User, equal func(user.User, user.User) bool) (user.User, error) {
	for _, v := range im.users {
		if equal(v, u) {
			return v, nil
		}
	}
	return u, UserNotFoundError(fmt.Errorf("GetUser: user not found"))
}

func (im *inmemory) AddGroup(ctx context.Context, u user.Group) error {
	_, err := im.GetGroup(ctx, u, func(u1, u2 user.Group) bool {
		if u1.Name == u2.Name {
			return true
		}
		return false
	})

	err, ok := err.(GroupNotFoundError)

	if err == nil {
		return fmt.Errorf("Adduser: user with given username (%s) exists!", u.Name)
	} else if !ok {
		return fmt.Errorf("Adduser: error while checking uniqueness: %w", err)
	}

	im.groups = append(im.groups, u)
	return nil
}

func (im *inmemory) GetGroup(ctx context.Context, u user.Group, equal func(user.Group, user.Group) bool) (user.Group, error) {
	for _, v := range im.groups {
		if equal(v, u) {
			return v, nil
		}
	}
	return u, GroupNotFoundError(fmt.Errorf("GetGroup: user not found"))
}


func NewInmemory() Storage {
  i := inmemory{}
  i.users = make([]user.User, 0)
  return &i
}
