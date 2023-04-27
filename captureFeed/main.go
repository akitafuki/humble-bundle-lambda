package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
)

func FindGameBundle(category []string, bundle string) bool {
	for _, c := range category {
		if c == bundle {
			return true
		}
	}
	return false
}

func getEnvVariable(key string) string {
	return os.Getenv(key)
}

func HandleLambdaEvent() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext("https://blog.humblebundle.com/feed/", ctx)

	r := regexp.MustCompile(`/?\?p=(?P<postid>\d{4,5})`)

	for _, item := range feed.Items {
		c := item.Categories
		if FindGameBundle(c, "Game Bundle") {
			log.Println(item.Title)
			log.Println(item.Link)

			s := r.FindStringSubmatch(item.GUID)
			log.Println(s[1])
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
