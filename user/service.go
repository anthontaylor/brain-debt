package user

import (
	"errors"

	brain "github.com/anthontaylor/brain-debt"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Add(brain.User) (*brain.UserID, error)
	Find(id *brain.UserID) (*brain.User, error)
	Update(id *brain.UserID, user brain.User) (*brain.User, error)
	Delete(id *brain.UserID) error
}

type service struct {
	users brain.UserRepository
}

func (s *service) Add(user brain.User) (*brain.UserID, error) {
	if user.FirstName == "" || user.LastName == "" {
		return nil, ErrInvalidArgument
	}
	return s.users.Add(&user)
}

func (s *service) Find(id *brain.UserID) (*brain.User, error) {
	if id == nil {
		return nil, ErrInvalidArgument
	}
	return s.users.Find(id)
}

func (s *service) Update(id *brain.UserID, user brain.User) (*brain.User, error) {
	if id == nil || user.FirstName == "" || user.LastName == "" {
		return nil, ErrInvalidArgument
	}
	return s.users.Update(id, user)
}

func (s *service) Delete(id *brain.UserID) error {
	if id == nil {
		return ErrInvalidArgument
	}
	return s.users.Delete(id)
}

func NewService(users brain.UserRepository) Service {
	return &service{users: users}
}
