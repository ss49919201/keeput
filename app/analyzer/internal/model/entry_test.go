package model_test

import (
	"cmp"
	"reflect"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/date"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestEntryPlatformIteratorOrderByPriorityAsc(t *testing.T) {
	types := []model.EntryPlatformType{}
	for ep := range model.EntryPlatformIteratorOrderByPriorityAsc() {
		types = append(types, ep.Type)
	}

	expect := []model.EntryPlatformType{
		model.EntryPlatformTypeHatena,
		model.EntryPlatformTypeZenn,
	}
	if !reflect.DeepEqual(types, expect) {
		t.Errorf("expect %v, actual %v", expect, types)
	}
}

func TestIsGoalAchieved(t *testing.T) {
	type args struct {
		publishedAt time.Time
		now         time.Time
		goalType    model.GoalType
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
				model.GoalTypeRecentWeek,
			},
			false,
		},
		{
			"achieve when entry is published exactly 7 days ago for recent week goal",
			args{
				time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
				model.GoalTypeRecentWeek,
			},
			true,
		},
		{
			"achieve when entry is published 30 days ago for recent month goal",
			args{
				time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
				model.GoalTypeRecentMonth,
			},
			true,
		},
		{
			"not achieve when entry is published 31 days ago for recent month goal",
			args{
				time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
				model.GoalTypeRecentMonth,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := model.IsGoalAchieved(tt.args.publishedAt, tt.args.now, tt.args.goalType); got != tt.want {
				t.Errorf("IsGoalAchieved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatest(t *testing.T) {
	type option struct {
		publishedAt time.Time
	}
	entryHatena := func(opt option) *model.Entry {
		return &model.Entry{
			Title:       "Rustを学ぶ",
			Body:        "Rustには所有権という概念があります",
			PublishedAt: cmp.Or(opt.publishedAt, time.Date(2025, 1, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))),
			Platform: model.EntryPlatform{
				Type:     model.EntryPlatformTypeHatena,
				Priority: 1,
			},
		}
	}
	entryZenn := func(opt option) *model.Entry {
		return &model.Entry{
			Title:       "Haskellを学ぶ",
			Body:        "Haskellは関数型言語です",
			PublishedAt: cmp.Or(opt.publishedAt, time.Date(2025, 1, 2, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))),
			Platform: model.EntryPlatform{
				Type:     model.EntryPlatformTypeZenn,
				Priority: 2,
			},
		}
	}

	tests := []struct {
		name    string
		entries []*model.Entry
		want    mo.Option[*model.Entry]
	}{
		{
			"returns latest entry when timestamps differ",
			[]*model.Entry{
				entryHatena(option{publishedAt: time.Date(2025, 1, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))}),
				entryHatena(option{publishedAt: time.Date(2025, 2, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))}),
			},
			mo.Some(entryHatena(option{publishedAt: time.Date(2025, 2, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))})),
		},
		{
			"returns higher priority when timestamps same",
			[]*model.Entry{
				entryZenn(option{publishedAt: time.Date(2025, 1, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))}),
				entryHatena(option{publishedAt: time.Date(2025, 1, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))}),
			},
			mo.Some(entryHatena(option{publishedAt: time.Date(2025, 1, 1, 12, 0, 0, 0, lo.ToPtr(date.LocationJST()))})),
		},
		{
			"returns none when entries slice is empty",
			[]*model.Entry{},
			mo.None[*model.Entry](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.Latest(tt.entries)
			assert.Equal(t, tt.want, got)
		})
	}
}
