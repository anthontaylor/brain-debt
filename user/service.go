package user

import (
	"errors"

	brain "github.com/anthontaylor/brain-debt"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Add(f, l string) (brain.UserID, error)
	Find(brain.UserID) (*brain.User, error)
}

type service struct {
	users brain.UserRepository
}

func (s *service) Add(f, l string) (brain.UserID, error) {
	if f == "" || l == "" {
		return "", ErrInvalidArgument
	}
	return s.users.Add(&brain.User{FirstName: f, LastName: l})
}

func (s *service) Find(id brain.UserID) (*brain.User, error) {
	if id == "" {
		return nil, ErrInvalidArgument
	}
	return s.users.Find(id)
}

func NewService(users brain.UserRepository) Service {
	return &service{users: users}
}
