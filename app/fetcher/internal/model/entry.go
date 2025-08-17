package model

import "time"

type EntryForRead struct {
	Title       string
	Body        string
	PublishedAt time.Time
}
