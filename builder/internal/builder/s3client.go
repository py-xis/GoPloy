package builder

import (
	"context"
	"log"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client() *s3.Client {
	// s3AccessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	// if s3AccessKeyID == "" {
	// 	log.Fatal("builder - NewS3Client : S3_ACCESS_KEY_ID environment variable not set")
	// }

	// s3SecretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	// if s3SecretAccessKey == "" {
	// 	log.Fatal("builder - NewS3Client: S3_SECRET_ACCESS_KEY environment variable not set")
	// }

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			// Fill
			credentials.NewStaticCredentialsProvider("ACCESS KEY", "SECRET", ""),
		),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	return s3.NewFromConfig(cfg)
}