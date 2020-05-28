package brain

import (
	"context"
	"errors"
)

type User struct {
	FirstName string
	LastName  string
}

type UserID string

type UserRepository interface {
	Add(ctx context.Context, user *User) (*UserID, error)
	Find(ctx context.Context, id *UserID) (*User, error)
	Update(ctx context.Context, id *UserID, user User) (*User, error)
	Delete(ctx context.Context, id *UserID) error
}

var (
	ErrUnknownUser   = errors.New("unknown user")
	ErrInvalidUserID = errors.New("error invalid user id")
	ErrAddingUser    = errors.New("error adding user")
)
