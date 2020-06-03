package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
)

type MyResponse struct {
	Message string `json:"Answer:"`
}

func FindGameBundle(category []string, bundle string) bool {
	for _, c := range category {
		if c == bundle {
			return true
		}
	}
	return false
}

func HandleLambdaEvent() (MyResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext("https://blog.humblebundle.com/feed/", ctx)

	items := feed.Items
	for _, item := range items {
		c := item.Categories
		if FindGameBundle(c, "Game Bundle") {
			log.Println(item.Title)
			log.Println(item.Link)
		}
	}

	return MyResponse{Message: "foo"}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
