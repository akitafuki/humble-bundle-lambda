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

type Bundle struct {
	URL       string `dynamodbav:"url"`
	Title     string `dynamodbav:"title"`
	CrawledAt string `dynamodbav:"crawledat"`
	EndDate   string `dynamodbav:"enddate"`
}

func HandleLambdaEvent() error {
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
			var newBundle Bundle

			newBundle.URL = record.(map[string]interface{})["product_url"].(string)
			newBundle.Title = record.(map[string]interface{})["tile_name"].(string)
			newBundle.CrawledAt = time.Now().Format("2006-01-02T15:04:05")
			newBundle.EndDate = record.(map[string]interface{})["end_date|datetime"].(string)

			av, err := dynamodbattribute.MarshalMap(newBundle)
			if err != nil {
				fmt.Println(fmt.Println(err))
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
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err = c.Visit(os.Getenv("RSS_FEED_URL"))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
