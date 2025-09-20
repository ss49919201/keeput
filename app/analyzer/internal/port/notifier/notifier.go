package notifier

import (
	"github.com/samber/mo"
)

type NotificationRequest struct {
	IsGoalAchieved bool `json:"is_goal_achieved"`
}

type Notify = func(*NotificationRequest) mo.Result[struct{}]