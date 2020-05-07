package topic

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

	addTopicHandler := kithttp.NewServer(
		makeAddTopicEndpoint(us),
		decodeAddTopicRequest,
		encodeResponse,
		opts...,
	)
	getTopicsHandler := kithttp.NewServer(
		makeGetTopicsEndpoint(us),
		decodeGetTopicsRequest,
		encodeResponse,
		opts...,
	)

	updateTopicHandler := kithttp.NewServer(
		makeUpdateTopicEndpoint(us),
		decodeUpdateTopicRequest,
		encodeResponse,
		opts...,
	)

	deleteTopicHandler := kithttp.NewServer(
		makeDeleteTopicEndpoint(us),
		decodeDeleteTopicRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/topics/", getTopicsHandler).Methods("GET")
	r.Handle("/topics/", addTopicHandler).Methods("POST")
	r.Handle("/topics/{id}", updateTopicHandler).Methods("PUT")
	r.Handle("/topics/{id}", deleteTopicHandler).Methods("DELETE")

	return r
}

func decodeGetTopicsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	userID := brain.UserID(body.UserID)
	return getTopicsRequest{ID: &userID}, nil
}

func decodeAddTopicRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	userID := brain.UserID(body.UserID)
	return addTopicRequest{ID: &userID, Name: body.Name}, nil
}

func decodeUpdateTopicRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	topicID, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	var body struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	userID := brain.UserID(body.UserID)
	tID := brain.TopicID(topicID)
	topic := &brain.Topic{ID: tID, Name: body.Name}
	return updateTopicRequest{UserID: &userID, Topic: topic}, nil
}

func decodeDeleteTopicRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	topicID, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	userID := brain.UserID(body.UserID)
	tID := brain.TopicID(topicID)
	return deleteTopicRequest{UserID: &userID, TopicID: &tID}, nil
}

var errBadRoute = errors.New("bad route")

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
	case brain.ErrTopicNotFound:
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
