package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/date"
	"github.com/stretchr/testify/assert"
)

func TestEntryPlatformIteratorOrderByPriorityAsc(t *testing.T) {
	types := []EntryPlatformType{}
	for ep := range EntryPlatformIteratorOrderByPriorityAsc() {
		types = append(types, ep.Type)
	}

	expect := []EntryPlatformType{
		EntryPlatformTypeHatena,
		EntryPlatformTypeZenn,
	}
	if !reflect.DeepEqual(types, expect) {
		t.Errorf("expect %v, actual %v", expect, types)
	}
}

func TestIsGoalAchieved(t *testing.T) {
	type args struct {
		publishedAt time.Time
		now         time.Time
		goalType    GoalType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"not achieve when entry is published 8 days ago for recent week goal",
			args{
				time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
				GoalTypeRecentWeek,
			},
			false,
		},
		{
			"achieve when entry is published exactly 7 days ago for recent week goal",
			args{
				time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
				GoalTypeRecentWeek,
			},
			true,
		},
		{
			"achieve when entry is published 30 days ago for recent month goal",
			args{
				time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
				GoalTypeRecentMonth,
			},
			true,
		},
		{
			"not achieve when entry is published 31 days ago for recent month goal",
			args{
				time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
				GoalTypeRecentMonth,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGoalAchieved(tt.args.publishedAt, tt.args.now, tt.args.goalType); got != tt.want {
				t.Errorf("IsGoalAchieved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatest(t *testing.T) {
	entryHatena := Entry{
		Title:       "Rustを学ぶ",
		Body:        "Rustには所有権という概念があります",
		PublishedAt: time.Date(2025, 1, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST())),
		Platform: EntryPlatform{
			Type:     EntryPlatformTypeHatena,
			Priority: 1,
		},
	}
	entryZenn := Entry{
		Title:       "Haskellを学ぶ",
		Body:        "Haskellは関数型言語です",
		PublishedAt: time.Date(2025, 1, 2, 12, 0, 0, 0, lo.ToPtr(date.LocationJST())),
		Platform: EntryPlatform{
			Type:     EntryPlatformTypeZenn,
			Priority: 2,
		},
	}

	tests := []struct {
		name    string
		entries []*Entry
		want    mo.Option[*Entry]
	}{
		{
			"returns latest entry when timestampls differ",
			[]*Entry{
				&entryHatena,
				&entryZenn,
			},
			mo.Some(&entryZenn),
		},
		{
			"returns higher priority when timestamps same",
			[]*Entry{
				&entryHatena,
				&entryZenn,
			},
			mo.Some(&entryZenn),
		},
		{
			"returns none when entries slice is empty",
			[]*Entry{},
			mo.None[*Entry](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Latest(tt.entries)
			assert.Equal(t, tt.want, got)
		})
	}
}
