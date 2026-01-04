# captureFeed

A Go-based AWS Lambda function that scrapes a specific webpage, extracts JSON data embedded in the page, and archives it to an AWS S3 bucket.

## Overview

This function is designed to:
1.  Visit a URL specified by the `RSS_FEED_URL` environment variable.
2.  Use [Colly](http://go-colly.org/) to scrape the page.
3.  Extract JSON content from a `<script id="landingPage-json-data">` tag.
4.  Validate the JSON data.
5.  Upload the data as a timestamped JSON file to an S3 bucket specified by `SALESBUNDLES_BUCKET_NAME`.

## Status

**Current Version**: Uses AWS SDK for Go v2.

## Prerequisites

*   Go 1.24 or later
*   AWS Credentials configured (for local testing or deployment context)

## Configuration

The Lambda function requires the following environment variables:

| Variable | Description |
| :--- | :--- |
| `RSS_FEED_URL` | The URL of the webpage to scrape. |
| `SALESBUNDLES_BUCKET_NAME` | The name of the S3 bucket where the extracted JSON will be stored. |

## Development

### Dependencies

This project uses Go modules. To install dependencies:

```bash
go mod tidy
```

### Building for AWS Lambda

The project includes a `Makefile` to simplify building the binary for the AWS Lambda execution environment (Linux/ARM64) and packaging it into a zip file.

To build the project:

```bash
make build-captureFeed
```

This command will:
1.  Compile the code for `GOOS=linux` and `GOARCH=arm64`.
2.  Output an executable named `bootstrap`.
3.  Create a `captureFeed.zip` file containing the bootstrap binary, ready for upload to AWS Lambda.

## Local Testing

You can run the code locally, provided you set the required environment variables and have valid AWS credentials configured in your environment or `~/.aws/credentials`.

```bash
export RSS_FEED_URL="https://example.com/feed"
export SALESBUNDLES_BUCKET_NAME="my-target-bucket"
go run main.go
```
