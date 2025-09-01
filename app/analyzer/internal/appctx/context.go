package appctx

import (
	"context"
	"time"
)

type key int

const (
	keyNow key = iota + 1
)

func SetNow(ctx context.Context, now time.Time) context.Context {
	return context.WithValue(ctx, keyNow, now)
}

func GetNow(ctx context.Context) (time.Time, bool) {
	now, ok := ctx.Value(keyNow).(time.Time)
	return now, ok
}

func GetNowOr(ctx context.Context, fallback time.Time) time.Time {
	now, ok := GetNow(ctx)
	if !ok {
		now = fallback
	}
	return now
}
