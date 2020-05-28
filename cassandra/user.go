package cassandra

import (
	"context"
	"fmt"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type userRepository struct {
	conn *gocql.Session
}

func NewUserRepository(conn *gocql.Session) brain.UserRepository {
	return &userRepository{conn: conn}
}

func (r *userRepository) Add(ctx context.Context, u *brain.User) (*brain.UserID, error) {
	userID := uuid.New().String()
	userGuid, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("user_guid=%s %w", userID, brain.ErrInvalidUserID)
	}
	query := r.conn.Query(`INSERT INTO user (user_guid, first_name, last_name) VALUES (?, ?, ?)`).
		Consistency(gocql.LocalQuorum).
		Bind(userGuid, u.FirstName, u.LastName)
	if err := query.Exec(); err != nil {
		return nil, fmt.Errorf("user_guid=%s %w err=%v", userID, brain.ErrAddingUser, err)
	}
	id := brain.UserID(userID)
	return &id, nil
}

func (r *userRepository) Find(ctx context.Context, id *brain.UserID) (*brain.User, error) {
	userGuid, err := gocql.ParseUUID(string(*id))
	if err != nil {
		return nil, fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidUserID)
	}
	query := r.conn.Query(`SELECT first_name, last_name FROM user WHERE user_guid = ?`, userGuid).
		Consistency(gocql.LocalQuorum).
		Bind(userGuid)

	var user brain.User
	if err := query.Scan(&user.FirstName, &user.LastName); err != nil {
		return nil, fmt.Errorf("user_guid=%v %w err=%v", *id, brain.ErrUnknownUser, err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id *brain.UserID, u brain.User) (*brain.User, error) {
	userGuid, err := gocql.ParseUUID(string(*id))
	if err != nil {
		return nil, fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidUserID)
	}
	query := r.conn.Query(`UPDATE user SET first_name = ?, last_name = ? WHERE user_guid = ?`).
		Consistency(gocql.LocalQuorum).
		Bind(u.FirstName, u.LastName, userGuid)
	if err := query.Exec(); err != nil {
		return nil, fmt.Errorf("user_guid=%v %w err=%v", *id, brain.ErrUnknownUser, err)
	}
	return &u, nil
}

func (r *userRepository) Delete(ctx context.Context, id *brain.UserID) error {
	userGuid, err := gocql.ParseUUID(string(*id))
	if err != nil {
		return fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidUserID)
	}
	query := r.conn.Query("DELETE FROM user WHERE user_guid = ?").
		Consistency(gocql.LocalQuorum).
		Bind(userGuid)
	if err := query.Exec(); err != nil {
		return fmt.Errorf("user_guid=%v %w err=%v", *id, brain.ErrUnknownUser, err)
	}
	return nil
}
