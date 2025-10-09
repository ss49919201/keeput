package config

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

func InitForLocal() error {
	if !isLocal() {
		return nil
	}
	return godotenv.Load()
}

var env = sync.OnceValue(func() string {
	return os.Getenv("ENV")
})

func isLocal() bool {
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

var s3BucketName = sync.OnceValue(func() string {
	return os.Getenv("S3_BUCKET_NAME")
})

func S3BucketName() string {
	return s3BucketName()
}
