package model

import (
	"time"

	"github.com/samber/mo"
)

type AnalysisReport struct {
	IsGoalAchieved bool `json:"is_goal_achieved"`

	LatestEntry mo.Option[*Entry] `json:"latest_entry"`
}

func Analyze(latestEntry mo.Option[*Entry], now time.Time, goalType GoalType) *AnalysisReport {
	if latestEntry.IsAbsent() {
		return &AnalysisReport{
			IsGoalAchieved: false,
			LatestEntry:    mo.None[*Entry](),
		}
	}

	return &AnalysisReport{
		IsGoalAchieved: IsGoalAchieved(latestEntry.MustGet().PublishedAt, now, goalType),
		LatestEntry:    latestEntry,
	}
}
