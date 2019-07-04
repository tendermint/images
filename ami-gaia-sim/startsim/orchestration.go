package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"strconv"
)

const (
	awsRegion       = "us-east-1"
	queueNamePrefix = "gaia-sim-"
)

func awsErrHandler(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		log.Fatal(awsErr.Error())
	}
	log.Fatal(err.Error())
}

func sendSqsMsg(seeds map[int]string) {
	maxMessages := 0
	batchRequestEntries := make([]*sqs.SendMessageBatchRequestEntry, len(seeds)-1)
	sqsSvc := sqs.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})))

	queues, err := sqsSvc.ListQueues(&sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(queueNamePrefix),
	})
	if err != nil {
		awsErrHandler(err)
	}
	for index := range seeds {
		maxMessages++
		batchRequestEntries[index] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(index)),
			MessageBody: aws.String("tick"), // Required field, we don't care about the body right now
		}
		if maxMessages == 10 || maxMessages >= len(seeds)-1 {
			_, err := sqsSvc.SendMessageBatch(&sqs.SendMessageBatchInput{
				Entries:  batchRequestEntries,
				QueueUrl: queues.QueueUrls[0],
			})
			if err != nil {
				awsErrHandler(err)
			}
			maxMessages = 0
			batchRequestEntries = batchRequestEntries[:0]
		}
		if index == len(batchRequestEntries) {
			fmt.Printf("%d", index)
			break
		}

	}
}
