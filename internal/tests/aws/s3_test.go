package aws_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	awscore "github.com/vkhangstack/hexagonal-architecture/internal/adapters/objectStorage"
)

func setUpS3() *awscore.S3Adapter {
	var err error
	client, err := awscore.NewS3Adapter(context.Background(), awscore.S3Config{
		Endpoint:        "http://localhost:9001",
		SecretAccessKey: "",
		AccessKeyID:     "",
		Bucket:          "test-bucket",
		UsePathStyle:    true, // required for MinIO and other local S3-compatible endpoints
	})

	if err != nil {
		fmt.Printf("Failed to create S3 client: %v\n", err)
		panic(err)
	}
	return client
}

func TestS3PutObject(t *testing.T) {
	s3Client := setUpS3()

	ctx := context.Background()
	err := s3Client.Put(ctx, awscore.PutInput{
		Key:         "test.txt",
		Body:        strings.NewReader("Hello, S3!"),
		ContentType: "text/plain",
	})
	if err != nil {
		fmt.Printf("Failed to put object: %v\n", err)
		panic(err)
	}
}

func TestS3GetObject(t *testing.T) {
	s3Client := setUpS3()

	ctx := context.Background()
	body, err := s3Client.Get(ctx, "name.png")
	if err != nil {
		t.Fatalf("Failed to get object: %v", err)
	}
	defer body.Close()

	info, err := s3Client.GetInfo(ctx, "name.png")
	if err != nil {
		t.Fatalf("Failed to get object info: %v", err)
	}

	t.Logf("Object Key: %s, Size: %d, ContentType: %s\n", info.Key, info.Size, info.ContentType)
}
