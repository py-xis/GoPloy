package builder

import (
	"context"
	"fmt"
	"os"
	"log"

	"github.com/redis/go-redis/v9"
)

func initRedisClient() *redis.Client {
	opt, err := redis.ParseURL("rediss://default:<TOKEN>@on-buffalo-26598.upstash.io:6379")
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	return redis.NewClient(opt)
}

var client = initRedisClient()

var projectID = os.Getenv("PROJECT_ID")

func PublishLog(message string) {
	ctx := context.Background()
	channel := fmt.Sprintf("logs:%s", projectID)
	client.Publish(ctx, channel, fmt.Sprintf(`{"log": %q}`, message))
}