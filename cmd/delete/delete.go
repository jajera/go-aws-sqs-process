package delete

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type DeleteMessageResult struct {
	MessageId string `json:"messageId"`
}

func DeleteMessage(queueUrl string, region string, messageReceiptHandle string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)

	input := &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: &messageReceiptHandle,
	}

	_, err = client.DeleteMessage(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to delete message, %v", err)
		return nil, err
	}

	log.Printf("Message has been deleted")
	return nil, nil
}
