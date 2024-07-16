package send

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SendMessageResult struct {
	MD5OfMessageAttributes       *string `json:"MD5OfMessageAttributes"`
	MD5OfMessageBody             *string `json:"MD5OfMessageBody"`
	MD5OfMessageSystemAttributes *string `json:"MD5OfMessageSystemAttributes"`
	MessageId                    *string `json:"MessageId"`
	SequenceNumber               *string `json:"SequenceNumber"`
}

func SendMessage(queueUrl string, region string, messageBody string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)

	input := &sqs.SendMessageInput{
		MessageBody: &messageBody,
		QueueUrl:    &queueUrl,
	}

	result, err := client.SendMessage(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to send message, %v", err)
		return nil, err
	}

	sendMessageResult := SendMessageResult{
		MD5OfMessageAttributes:       result.MD5OfMessageAttributes,
		MD5OfMessageBody:             result.MD5OfMessageBody,
		MD5OfMessageSystemAttributes: result.MD5OfMessageSystemAttributes,
		MessageId:                    result.MessageId,
		SequenceNumber:               result.SequenceNumber,
	}

	jsonOutput, err := json.Marshal(sendMessageResult)
	if err != nil {
		log.Fatalf("failed to marshal messages to JSON, %v", err)
		return nil, err
	}

	log.Printf("Send message result: %s", jsonOutput)
	return jsonOutput, nil
}
