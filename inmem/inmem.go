package inmem

import (
	"sync"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/google/uuid"
)

type userRepository struct {
	mu    sync.RWMutex
	users map[brain.UserID]*brain.User
}

func (r *userRepository) Add(u *brain.User) (*brain.UserID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	userID := brain.UserID(uuid.New().String())
	r.users[userID] = u
	return &userID, nil
}

func (r *userRepository) Find(id *brain.UserID) (*brain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.users[*id]; ok {
		return val, nil
	}
	return nil, brain.ErrUnknownUser
}

func (r *userRepository) Update(id *brain.UserID, u brain.User) (*brain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[*id]; ok {
		r.users[*id] = &u
		return &u, nil
	}
	return nil, brain.ErrUnknownUser
}

func (r *userRepository) Delete(id *brain.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[*id]; ok {
		delete(r.users, *id)
		return nil
	}
	return brain.ErrUnknownUser
}

func NewUserRepository() brain.UserRepository {
	return &userRepository{
		users: make(map[brain.UserID]*brain.User),
	}
}
