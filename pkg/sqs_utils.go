package receive

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Message struct {
	MessageId              string            `json:"messageId"`
	Body                   string            `json:"body"`
	MD5OfBody              string            `json:"md5OfBody"`
	MD5OfMessageAttributes string            `json:"md5OfMessageAttributes"`
	MessageAttributes      map[string]string `json:"messageAttributes"`
	Attributes             map[string]string `json:"attributes"`
	ReceiptHandle          string            `json:"receiptHandle"`
}

func ReceiveMessages(queueUrl string, maxMessages int, region string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: int32(maxMessages),
		MessageAttributeNames: []string{
			"All",
		},
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
	}

	result, err := client.ReceiveMessage(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to receive messages, %v", err)
	}

	if len(result.Messages) == 0 {
		log.Println("No messages received")
		return
	}

	var messages []Message
	for _, msg := range result.Messages {
		attributes := make(map[string]string)
		for key, value := range msg.Attributes {
			attributes[key] = value
		}

		messageAttributes := make(map[string]string)
		for key, value := range msg.MessageAttributes {
			if value.StringValue != nil {
				messageAttributes[key] = *value.StringValue
			}
		}

		messages = append(messages, Message{
			MessageId:              *msg.MessageId,
			Body:                   *msg.Body,
			MD5OfBody:              *msg.MD5OfBody,
			MD5OfMessageAttributes: *msg.MD5OfMessageAttributes,
			MessageAttributes:      messageAttributes,
			Attributes:             attributes,
			ReceiptHandle:          *msg.ReceiptHandle,
		})
	}

	jsonOutput, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal messages to JSON, %v", err)
	}

	fmt.Println(string(jsonOutput))
}
