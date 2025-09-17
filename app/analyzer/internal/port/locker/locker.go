package locker

import (
	"context"

	"github.com/samber/mo"
)

type Acquire = func(ctx context.Context, lockID string) mo.Result[bool]
type Release = func(ctx context.Context, lockID string) error
