package brain

import (
	"errors"
)

type User struct {
	FirstName string
	LastName  string
}

type UserID string

type UserRepository interface {
	Add(user *User) (UserID, error)
	Find(id UserID) (*User, error)
}

var ErrUnknownUser = errors.New("unknown user")
