package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gocolly/colly/v2"
	"os"
	"time"
)

type bundle struct {
	URL       string
	Title     string
	CrawledAt time.Time
}

func HandleLambdaEvent() error {
	bundles := []bundle{}

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML("script#landingPage-json-data", func(e *colly.HTMLElement) {
		var js map[string]interface{}

		err := json.Unmarshal([]byte(e.Text), &js)

		if err != nil {
			panic(errors.New("can't parse the data in script#landingPage-json-data"))
		}

		var products = js["data"].(map[string]interface{})["games"].(map[string]interface{})["mosaic"].([]interface{})[0].(map[string]interface{})["products"]

		for _, record := range products.([]interface{}) {
			var newBundle bundle

			newBundle.URL = record.(map[string]interface{})["product_url"].(string)
			newBundle.Title = record.(map[string]interface{})["tile_name"].(string)
			newBundle.CrawledAt = time.Now()

			bundles = append(bundles, newBundle)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(os.Getenv("RSS_FEED_URL"))

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create a session.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create a DynamoDB client.
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(bundles)
	if err != nil {
		fmt.Println(fmt.Println(err))
		return err
	}

	// Create a PutItemInput object.
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String("humble-data"),
		Item:      av,
	}

	// Call the PutItem method of the DynamoDB client.
	_, err = svc.PutItem(putItemInput)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
