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

var feedURLHatena = sync.OnceValue(func() string {
	return os.Getenv("FEED_URL_HATENA")
})

func FeedURLHatena() string {
	return feedURLHatena()
}

var logLevel = sync.OnceValue(func() string {
	return os.Getenv("LOG_LEVEL")
})

func LogLevel() string {
	return logLevel()
}
