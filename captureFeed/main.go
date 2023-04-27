package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"time"
)

func FindGameBundle(category []string, bundle string) bool {
	for _, c := range category {
		if c == bundle {
			return true
		}
	}
	return false
}

type item struct {
	URL       string
	Title     string
	CrawledAt time.Time
}

func HandleLambdaEvent() error {
	posts := []item{}

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML("info-section", func(e *colly.HTMLElement) {
		temp := item{}
		temp.URL = e.Attr("href")
		temp.Title = e.ChildText("span")
		temp.CrawledAt = time.Now()
		posts = append(posts, temp)
		log.Println(temp)
	})

	c.Visit(os.Getenv("RSS_FEED_URL"))

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
