package locker

import "github.com/samber/mo"

type Acquire = func(lockID string) mo.Result[bool]
type Release = func(lockID string) error
