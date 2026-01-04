build-captureFeed:
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
	zip captureFeed.zip bootstrap