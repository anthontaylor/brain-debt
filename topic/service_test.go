package topic

import (
	"testing"

	brain "github.com/anthontaylor/brain-debt"
	mock_topic "github.com/anthontaylor/brain-debt/topic/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	successTests := []struct {
		user_id brain.UserID
		name    string
		err     error
	}{
		{brain.UserID("123"), "operating systems", nil},
		{brain.UserID("123"), "compilers", nil},
	}
	mockService := mock_topic.NewMockService(ctrl)
	for _, test := range successTests {
		topic_id := brain.TopicID(uuid.New().String())
		mockService.EXPECT().Add(gomock.Any(), gomock.Any()).Return(&topic_id, nil)
		topics := NewService(mockService)
		got, err := topics.Add(&test.user_id, test.name)
		assert.Same(t, got, &topic_id)
		assert.Equal(t, test.err, err)
	}

	failureTests := []struct {
		user_id brain.UserID
		name    string
		err     error
	}{
		{brain.UserID("123"), "", ErrInvalidArgument},
		{brain.UserID(""), "123", ErrInvalidArgument},
	}

	for _, test := range failureTests {
		mockService.EXPECT().Add(gomock.Any(), gomock.Any()).MaxTimes(0)
		topics := NewService(mockService)
		got, err := topics.Add(&test.user_id, test.name)
		assert.Equal(t, test.err, err)
		assert.Nil(t, got)
	}

	nilTests := []struct {
		user_id *brain.UserID
		name    string
		err     error
	}{
		{nil, "", ErrInvalidArgument},
		{nil, "123", ErrInvalidArgument},
	}

	for _, test := range nilTests {
		mockService.EXPECT().Add(gomock.Any(), gomock.Any()).MaxTimes(0)
		topics := NewService(mockService)
		got, err := topics.Add(test.user_id, test.name)
		assert.Equal(t, test.err, err)
		assert.Nil(t, got)
	}
}

func TestGetTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	successTests := []struct {
		user_id brain.UserID
		topics  []brain.Topic
		err     error
	}{
		{brain.UserID("456"), []brain.Topic{}, nil},
	}

	mockService := mock_topic.NewMockService(ctrl)
	for _, test := range successTests {
		mockService.EXPECT().Get(gomock.Any()).Return(test.topics, nil)
		topics := NewService(mockService)
		got, err := topics.Get(&test.user_id)
		assert.Equal(t, got, test.topics)
		assert.Equal(t, test.err, err)
	}

	failureTests := []struct {
		user_id brain.UserID
		topics  []brain.Topic
		err     error
	}{
		{brain.UserID(""), nil, ErrInvalidArgument},
	}

	for _, test := range failureTests {
		mockService.EXPECT().Get(gomock.Any()).MaxTimes(0)
		topics := NewService(mockService)
		_, err := topics.Get(&test.user_id)
		assert.Equal(t, test.err, err)
	}
}
