package main

import (
	"log"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	awsRegion       = "us-east-1"
	queueNamePrefix = "gaia-sim-"
)

var (
	sessionSQS = sqs.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})))
)

func sendBatch(batchRequestEntries []*sqs.SendMessageBatchRequestEntry, queues *sqs.ListQueuesOutput) {
	if _, err = sessionSQS.SendMessageBatch(&sqs.SendMessageBatchInput{
		Entries:  batchRequestEntries,
		QueueUrl: queues.QueueUrls[0],
	}); err != nil {
		log.Printf(err.Error())
	}
}

func removeEmpties(batch []*sqs.SendMessageBatchRequestEntry) []*sqs.SendMessageBatchRequestEntry {
	var newBatch []*sqs.SendMessageBatchRequestEntry
	for _, msg := range batch {
		if msg != nil {
			newBatch = append(newBatch, msg)
		}
	}
	return newBatch
}

func sendSqsMsg(instanceIndex []int) {
	var (
		maxMessages         = 0
		batchRequestEntries = make([]*sqs.SendMessageBatchRequestEntry, 10)
		queues              *sqs.ListQueuesOutput
	)

	if queues, err = sessionSQS.ListQueues(&sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(queueNamePrefix),
	}); err != nil {
		log.Fatal(err.Error())
	}
	sort.Ints(instanceIndex)
	for index := range instanceIndex {
		batchRequestEntries[maxMessages] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(index)),
			MessageBody: aws.String("Instance " + strconv.Itoa(index)), // Required field, we don't care about the body right now
		}
		maxMessages++
		// SQS only accepts batches of max 10 messages
		if maxMessages == 10 {
			sendBatch(batchRequestEntries, queues)
			batchRequestEntries = make([]*sqs.SendMessageBatchRequestEntry, 10)
			maxMessages = 0
		}
		// We want the queue length to be one less than the number of instances that were created
		// An empty queue will prompt the last running instance to send the simulation finished message
		if index == len(instanceIndex)-2 {
			// Can't have nil elements in the list or the sqs send function will segfault
			batchRequestEntries = removeEmpties(batchRequestEntries)
			sendBatch(batchRequestEntries, queues)
			break
		}
	}
}
