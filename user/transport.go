package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"

	brain "github.com/anthontaylor/brain-debt"
)

func MakeHandler(us Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	addUserHandler := kithttp.NewServer(
		makeAddUserEndpoint(us),
		decodeAddUserRequest,
		encodeResponse,
		opts...,
	)

	getUserHandler := kithttp.NewServer(
		makeGetUserEndpoint(us),
		decodeGetUserRequest,
		encodeResponse,
		opts...,
	)

	updateUserHandler := kithttp.NewServer(
		makeUpdateUserEndpoint(us),
		decodeUpdateUserRequest,
		encodeResponse,
		opts...,
	)

	deleteUserHandler := kithttp.NewServer(
		makeDeleteUserEndpoint(us),
		decodeDeleteUserRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/user/", addUserHandler).Methods("POST")
	r.Handle("/user/{id}", getUserHandler).Methods("GET")
	r.Handle("/user/{id}", updateUserHandler).Methods("PUT")
	r.Handle("/user/{id}", deleteUserHandler).Methods("DELETE")

	return r
}

var errBadRoute = errors.New("bad route")

func decodeAddUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return addUserRequest{FirstName: body.FirstName, LastName: body.LastName}, nil
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	bID := brain.UserID(id)
	return getUserRequest{ID: &bID}, nil
}

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	var body struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	bID := brain.UserID(id)
	return updateUserRequest{ID: &bID, FirstName: body.FirstName, LastName: body.LastName}, nil
}

func decodeDeleteUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	bID := brain.UserID(id)
	return deleteUserRequest{ID: &bID}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case brain.ErrUnknownUser:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
