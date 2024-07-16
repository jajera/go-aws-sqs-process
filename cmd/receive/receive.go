package receive

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Message struct {
	MessageId              string            `json:"messageId"`
	Body                   string            `json:"body"`
	MD5OfBody              string            `json:"md5OfBody"`
	MD5OfMessageAttributes string            `json:"md5OfMessageAttributes,omitempty"`
	MessageAttributes      map[string]string `json:"messageAttributes"`
	Attributes             map[string]string `json:"attributes"`
	ReceiptHandle          string            `json:"receiptHandle"`
}

func ReceiveMessages(queueUrl string, region string, maxMessages int) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
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
		return nil, err
	}

	var messages []Message
	if len(result.Messages) == 0 {
		log.Println("No messages received")
		return json.Marshal(messages)
	}

	for _, msg := range result.Messages {
		attributes := make(map[string]string)
		if msg.Attributes != nil {
			for key, value := range msg.Attributes {
				attributes[key] = value
			}
		}

		messageAttributes := make(map[string]string)
		if msg.MessageAttributes != nil {
			for key, value := range msg.MessageAttributes {
				if value.StringValue != nil {
					messageAttributes[key] = *value.StringValue
				}
			}
		}

		var md5OfMessageAttributes string
		if msg.MD5OfMessageAttributes != nil {
			md5OfMessageAttributes = *msg.MD5OfMessageAttributes
		}

		messages = append(messages, Message{
			MessageId:              *msg.MessageId,
			Body:                   *msg.Body,
			MD5OfBody:              *msg.MD5OfBody,
			MD5OfMessageAttributes: md5OfMessageAttributes,
			MessageAttributes:      messageAttributes,
			Attributes:             attributes,
			ReceiptHandle:          *msg.ReceiptHandle,
		})
	}

	jsonOutput, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal messages to JSON, %v", err)
		return nil, err
	}

	log.Printf("Receive message(s) result: %s", jsonOutput)
	return jsonOutput, nil
}
