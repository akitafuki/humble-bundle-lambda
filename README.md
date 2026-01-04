# Sales Bundle Lambda

An AWS Lambda function written in Go that scrapes sales data from a configured URL and uploads the raw JSON payload to an S3 bucket.

## Overview

This project uses [Colly](http://go-colly.org/) to visit a webpage, extracts JSON data from a specific script tag (`script#landingPage-json-data`), and saves it as a `.json` file in an S3 bucket.

## Configuration

The Lambda function requires the following environment variables:

- `RSS_FEED_URL`: The URL of the page to scrape (e.g., a Humble Bundle landing page).
- `S3_BUCKET_NAME`: The name of the S3 bucket where the JSON files will be stored.
- `AWS_REGION`: (Optional) The AWS region for the S3 bucket (defaults to `us-east-1`).

## Local Development

### Prerequisites

- Go 1.25+
- AWS CLI configured with appropriate permissions.

### Testing Locally

You can run an integration test that performs a real scrape and upload. Ensure your environment variables are set before running the test:

```bash
cd captureFeed

# Set environment variables
export RSS_FEED_URL="https://www.humblebundle.com/books"
export S3_BUCKET_NAME="your-s3-bucket-name"

# Run the integration test
go test -v -run TestHandleLambdaEvent_Integration
```

*Note: If environment variables are missing, the test will be skipped.*

## Deployment

To build the binary for AWS Lambda (Amazon Linux 2023):

```bash
cd captureFeed
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
zip deployment.zip bootstrap
```

## Deploying with AWS SAM

This project includes a `template.yaml` for AWS SAM (Serverless Application Model). This automates the creation of the Lambda function, the S3 bucket, and the necessary IAM permissions.

### Prerequisites

- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) installed.
- Docker (optional, but recommended for local testing with SAM).

### Steps

1. **Build the project**
   The SAM CLI will compile the Go code and prepare the artifacts.
   ```bash
   sam build
   ```

2. **Deploy to AWS**
   This command will guide you through the deployment process (setting region, parameter values, etc.).
   ```bash
   sam deploy --guided
   ```
   
   During the guided prompt:
   - **Stack Name**: e.g., `sales-bundle-stack`
   - **AWS Region**: e.g., `us-east-1`
   - **Parameter RssFeedUrl**: Accept the default or provide your own.
   - **Confirm changes before deploy**: `Y`
   - **Allow SAM CLI IAM role creation**: `Y`

Once deployed, the S3 bucket name will be displayed in the outputs.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
