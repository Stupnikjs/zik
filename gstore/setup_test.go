package gstore

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
)

var TestBucketName = "someidealtestbucketname"

func TestMain(m *testing.M) {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	code := m.Run()

	cleanup()
	// Exit with the received code
	os.Exit(code)
}

func setup() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatal(err)
	}

	bucket := client.Bucket(TestBucketName)

	bucketAttrs := &storage.BucketAttrs{
		Location: "EU",
	}

	projectID := "zikstudio"
	// Creates the new bucket.

	if err := bucket.Create(ctx, projectID, bucketAttrs); err != nil {
		fmt.Println("here")
		log.Fatalf("Failed to create bucket: %v", err)
	}

}

func cleanup() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatal(err)
	}

	bucket := client.Bucket(TestBucketName)

	err = bucket.Delete(ctx)
	fmt.Println("here", err)

}
