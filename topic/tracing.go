package topic

import (
	"context"
	brain "github.com/anthontaylor/brain-debt"
	opentracing "github.com/opentracing/opentracing-go"
	"time"
)

type tracingService struct {
	Service
}

func NewTracingService(s Service) Service {
	return &tracingService{Service: s}
}

func (t *tracingService) Add(ctx context.Context, id *brain.UserID, name string) (_ *brain.TopicID, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add_topic")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "add",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Add(ctx, id, name)
}

func (t *tracingService) Get(ctx context.Context, id *brain.UserID) (_ []brain.Topic, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get_topics")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "find",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Get(ctx, id)
}

func (t *tracingService) Update(ctx context.Context, id *brain.UserID, topic *brain.Topic) (_ *brain.Topic, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "update_topic")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "update",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Update(ctx, id, topic)
}

func (t *tracingService) Delete(ctx context.Context, id *brain.UserID, topicID *brain.TopicID) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete_topic")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "delete",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Delete(ctx, id, topicID)
}
