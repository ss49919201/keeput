package config

import (
	"os"
	"sync"
)

var feedURLZenn = sync.OnceValue(func() string {
	return os.Getenv("FEED_URL_ZENN")
})

func FeedURLZenn() string {
	return feedURLZenn()
}
