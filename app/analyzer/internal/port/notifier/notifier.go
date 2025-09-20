package notifier

import "context"

// NotifyAnalysisReport sends a notification based on goal achievement.
// true -> success message, false -> failure message.
type NotifyAnalysisReport = func(context.Context, bool) error
