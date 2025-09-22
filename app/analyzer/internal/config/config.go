package config

import (
	"os"
	"strings"
	"sync"
)

var env = sync.OnceValue(func() string {
	return os.Getenv("ENV")
})

func IsLocal() bool {
	return strings.ToLower(env()) == "local"
}

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

var lockerAPIKeyCloudflareWorker = sync.OnceValue(func() string {
	return os.Getenv("LOCKER_API_KEY_CLOUDFLARE_WORKER")
})

func LockerAPIKeyCloudflareWorker() string {
	return lockerAPIKeyCloudflareWorker()
}

var discordWebhookURL = sync.OnceValue(func() string {
	return os.Getenv("DISCORD_WEBHOOK_URL")
})

func DiscordWebhookURL() string {
	return discordWebhookURL()
}
