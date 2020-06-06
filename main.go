package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
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

func HandleLambdaEvent() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext("https://blog.humblebundle.com/feed/", ctx)

	r := regexp.MustCompile(`/?\?p=(?P<postid>\d{4,5})`)

	mySession := session.Must(session.NewSession())

	svc := rdsdataservice.New(mySession, aws.NewConfig().WithRegion("us-west-2"))

	for _, item := range feed.Items {
		c := item.Categories
		if FindGameBundle(c, "Game Bundle") {
			log.Println(item.Title)
			log.Println(item.Link)

			s := r.FindStringSubmatch(item.GUID)
			log.Println(s[1])

			SQLStatement := fmt.Sprintf(`SELECT %s FROM HumbleBundlePosts;`, s[1])

			req, resp := svc.ExecuteStatementRequest(&rdsdataservice.ExecuteStatementInput{
				Database:    aws.String("akitafuki"),
				ResourceArn: aws.String("arn:aws:rds:us-west-2:934040383244:cluster:humblebundle"),
				SecretArn:   aws.String("arn:aws:secretsmanager:us-west-2:934040383244:secret:humblebundle-OQp5wy"),
				Sql:         aws.String(SQLStatement),
			})

			err1 := req.Send()
			if err1 == nil { // resp is now filled
				fmt.Println("Response:", resp)
			} else {
				fmt.Println("error:", err1)
			}

			if len(resp.Records) == 0 {
				log.Println("post not found")
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
