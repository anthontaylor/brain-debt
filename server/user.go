package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/anthontaylor/brain-debt/user"
)

type userHandler struct {
	s      user.Service
	logger kitlog.Logger
}

func (u *userHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", u.addUser)
	r.Get("/{userID}", u.getUser)

	return r
}

func (u *userHandler) addUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		u.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}

	id, err := u.s.Add(request.FirstName, request.LastName)
	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = struct {
		ID brain.UserID `json:"user_id"`
	}{
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		u.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}
}

func (u *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	userID := chi.URLParam(r, "userID")

	user, err := u.s.Find(brain.UserID(userID))
	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = struct {
		User *brain.User `json:"user"`
	}{User: user}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		u.logger.Log("error", err)
		encodeError(ctx, err, w)
		return
	}
}
