package inmem

import (
	"errors"
	"sync"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/google/uuid"
)

type userRepository struct {
	mu    sync.RWMutex
	users map[brain.UserID]*brain.User
}

func (r *userRepository) Add(u *brain.User) (brain.UserID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	userID := brain.UserID(uuid.New().String())
	r.users[userID] = u
	return userID, nil
}

func (r *userRepository) Find(id brain.UserID) (*brain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.users[id]; ok {
		return val, nil
	}
	return nil, errors.New("User not found")
}

func NewUserRepository() brain.UserRepository {
	return &userRepository{
		users: make(map[brain.UserID]*brain.User),
	}
}
