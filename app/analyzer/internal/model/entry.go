package model

import (
	"iter"
	"time"

	"github.com/samber/lo"
	"github.com/ss49919201/keeput/app/analyzer/internal/date"
)

type Entry struct {
	Title       string
	Body        string
	PublishedAt time.Time
}

type GoalType int

const (
	GoalTypeRecentWeek  GoalType = iota + 1
	GoalTypeRecentMonth GoalType = iota + 1
)

func IsGoalAchieved(publishedAt, now time.Time, goalType GoalType) bool {
	beginningOfToday := date.BeginningOfDay(now)
	beginningOfBeforeXdays := date.AddDays(
		beginningOfToday,
		lo.If(goalType == GoalTypeRecentWeek, -7).
			ElseIf(goalType == GoalTypeRecentMonth, -30).
			Else(-7),
	)

	return !publishedAt.Before(beginningOfBeforeXdays)
}

type EntryPlatformType int

func (e EntryPlatformType) IsZero() bool {
	return e == EntryPlatformTypeZero
}

const (
	EntryPlatformTypeZero EntryPlatformType = iota
	EntryPlatformTypeZenn
	EntryPlatformTypeHatena
)

type EntryPlatform struct {
	entryPlatformType EntryPlatformType
	priority          int
}

func (e *EntryPlatform) Type() EntryPlatformType {
	return e.entryPlatformType
}

// DANGER: 最代入禁止
// 要素の順序は優先度の昇順
var entryPlatforms = []*EntryPlatform{
	{
		entryPlatformType: EntryPlatformTypeHatena,
		priority:          1,
	},
	{
		entryPlatformType: EntryPlatformTypeZenn,
		priority:          2,
	},
}

func EntryPlatformIteratorOrderByPriorityAsc() iter.Seq[*EntryPlatform] {
	return func(yield func(*EntryPlatform) bool) {
		for i := range len(entryPlatforms) {
			if !yield(entryPlatforms[i]) {
				break
			}
		}
	}
}
