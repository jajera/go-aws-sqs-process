package main

import (
	"context"
	"flag"
	"log"

	"go-aws-sqs-process/cmd/delete"
	"go-aws-sqs-process/cmd/receive"
	"go-aws-sqs-process/cmd/send"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	var defaultMaxMessages = 3
	task := flag.String("task", "receive", "The task to perform: send, receive, delete")
	queueUrl := flag.String("queueUrl", "", "The URL of the SQS queue")
	region := flag.String("region", "", "The AWS region")
	maxMessages := flag.Int("maxMessages", defaultMaxMessages, "The maximum number of messages to retrieve from the SQS queue (default is 3)")
	messageBody := flag.String("messageBody", "", "The body of the message to send")
	receiptHandle := flag.String("receiptHandle", "", "The receipt handle of the message to delete")
	flag.Parse()

	if !(*task == "delete" || *task == "receive" || *task == "send") {
		log.Fatalf("Valid tasks to perform are: delete, receive, send")
	}

	if *queueUrl == "" {
		log.Fatalf("queueUrl is required")
	}

	if *region != "" {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*region))
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		*region = cfg.Region
	}

	if *task == "delete" {
		if !(*maxMessages != 0 || *maxMessages != defaultMaxMessages) {
			log.Fatalf("maxMessages is not required for delete task")
		}
	} else if *task == "receive" {
		if *maxMessages <= 0 {
			log.Fatalf("maxMessages must be greater than 0 for receive task")
		}
	} else if *task == "send" {
		if !(*maxMessages != 0 || *maxMessages != defaultMaxMessages) {
			log.Fatalf("maxMessages is not required for send task")
		}
	}

	if *messageBody == "" && *task == "send" {
		log.Fatalf("messageBody is required for send task and must not be empty")
	}

	if *messageBody != "" {
		if *task == "delete" {
			log.Fatalf("messageBody is NOT required for delete task")
		} else if *task == "receive" {
			log.Fatalf("messageBody is NOT required for receive task")
		}
	}

	if *receiptHandle == "" && *task == "delete" {
		log.Fatalf("receiptHandle is required for delete task")
	}

	switch *task {
	case "send":
		if *messageBody == "" {
			log.Fatalf("messageBody is required for send task")
		}
		send.SendMessage(*queueUrl, *region, *messageBody)
	case "receive":
		receive.ReceiveMessages(*queueUrl, *region, *maxMessages)
	case "delete":
		delete.DeleteMessage(*queueUrl, *region, *receiptHandle)
	default:
		log.Fatalf("Unknown task: %s", *task)
	}
}
