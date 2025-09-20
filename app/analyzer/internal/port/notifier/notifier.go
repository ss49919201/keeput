package notifier

import "context"

// NotifyByIsGoalAchieved sends a notification based on goal achievement.
// true -> success message, false -> failure message.
type NotifyByIsGoalAchieved = func(context.Context, bool) error

