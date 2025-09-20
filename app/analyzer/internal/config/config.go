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

var lockerURLCloudflareWorker = sync.OnceValue(func() string {
    return os.Getenv("LOCKER_URL_CLOUDFLARE_WORKER")
})

func LockerURLCloudflareWorker() string {
    return lockerURLCloudflareWorker()
}

var slackWebhookURL = sync.OnceValue(func() string {
    return os.Getenv("SLACK_WEBHOOK_URL")
})

func SlackWebhookURL() string {
    return slackWebhookURL()
}
