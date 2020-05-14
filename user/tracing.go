package user

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

func (t *tracingService) Add(ctx context.Context, user brain.User) (id *brain.UserID, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add_user")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "add",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Add(ctx, user)
}

func (t *tracingService) Find(ctx context.Context, id *brain.UserID) (user *brain.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get_user")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "find",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Find(ctx, id)
}

func (t *tracingService) Update(ctx context.Context, id *brain.UserID, user brain.User) (_ *brain.User, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "update_user")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "update",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Update(ctx, id, user)
}

func (t *tracingService) Delete(ctx context.Context, id *brain.UserID) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete_user")
	defer span.Finish()
	defer func(begin time.Time) {
		span.LogKV(
			"method", "delete",
			"took", time.Since(begin).String(),
		)
	}(time.Now())
	return t.Service.Delete(ctx, id)
}
