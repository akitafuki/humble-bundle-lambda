package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gocolly/colly/v2"
)

func HandleLambdaEvent(ctx context.Context) error {
	// Check for required environment variables
	bucketName := os.Getenv("SALESBUNDLES_BUCKET_NAME")
	if bucketName == "" {
		return errors.New("SALESBUNDLES_BUCKET_NAME environment variable is not set")
	}

	feedURL := os.Getenv("RSS_FEED_URL")
	if feedURL == "" {
		return errors.New("RSS_FEED_URL environment variable is not set")
	}

	// Load the Shared AWS Configuration.
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		fmt.Println("Error loading config:", err)
		return err
	}

	// Create an S3 client.
	svc := s3.NewFromConfig(cfg)

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong with the collector:", err)
	})

	c.OnHTML("script#landingPage-json-data", func(e *colly.HTMLElement) {
		jsonData := []byte(e.Text)

		// Validate that the text is actually JSON
		if !json.Valid(jsonData) {
			fmt.Println("Error: Extracted text is not valid JSON")
			return
		}

		// Create a unique key for the S3 object
		timestamp := time.Now().Format("2006-01-02T15-04-05")
		key := fmt.Sprintf("bundle-data-%s.json", timestamp)

		// Upload the JSON file to S3
		_, err := svc.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        bytes.NewReader(jsonData),
			ContentType: aws.String("application/json"),
		})

		if err != nil {
			fmt.Printf("Failed to upload data to S3 bucket '%s' with key '%s': %v\n", bucketName, key, err)
			return
		}

		fmt.Printf("Successfully uploaded data to s3://%s/%s\n", bucketName, key)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err = c.Visit(feedURL)
	if err != nil {
		fmt.Println("Error visiting URL:", err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}