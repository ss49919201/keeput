package model

import (
	"cmp"
	"slices"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/date"
)

type Entry struct {
	Title       string
	Body        string
	PublishedAt time.Time
	Platform    EntryPlatform
}

type GoalType int

const (
	GoalTypeRecentWeek  GoalType = iota + 1
	GoalTypeRecentMonth GoalType = iota + 1
)

// 現在日の00:00(JST)からn日遡った日時以降に公開されていれば目標達成とみなす
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

const (
	EntryPlatformTypeZero EntryPlatformType = iota
	EntryPlatformTypeZenn
	EntryPlatformTypeHatena
)

type EntryPlatform struct {
	Type     EntryPlatformType
	Priority int
}

func EntryPlatformHatena() EntryPlatform {
	return EntryPlatform{
		Type:     EntryPlatformTypeHatena,
		Priority: 1,
	}
}

func EntryPlatformZenn() EntryPlatform {
	return EntryPlatform{
		Type:     EntryPlatformTypeZenn,
		Priority: 2,
	}
}

func Latest(entries []*Entry) mo.Option[*Entry] {
	if len(entries) < 1 {
		return mo.None[*Entry]()
	}
	cloned := slices.Clone(entries)
	return mo.Some(
		slices.MinFunc(cloned, func(a, b *Entry) int {
			return cmp.Or(
				b.PublishedAt.Compare(a.PublishedAt),
				cmp.Compare(a.Platform.Priority, b.Platform.Priority),
			)
		}),
	)
}
