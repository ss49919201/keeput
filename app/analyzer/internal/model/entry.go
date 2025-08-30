package model

import (
	"iter"
	"time"
)

type Entry struct {
	Title       string
	Body        string
	PublishedAt time.Time
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

// DANGER: 最大入禁止
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
