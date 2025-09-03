package model

import (
	"reflect"
	"testing"
	"time"
)

func TestEntryPlatformIteratorOrderByPriorityAsc(t *testing.T) {
	types := []EntryPlatformType{}
	for ep := range EntryPlatformIteratorOrderByPriorityAsc() {
		types = append(types, ep.Type())
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
