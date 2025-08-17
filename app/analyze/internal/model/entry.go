package model

import "time"

type Entry struct {
	Title       string
	Body        string
	PublishedAt time.Time
}
