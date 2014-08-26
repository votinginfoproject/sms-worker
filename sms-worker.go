package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/env"
	"github.com/votinginfoproject/sms-worker/response"
	"github.com/votinginfoproject/sms-worker/sms"
	"github.com/votinginfoproject/sms-worker/util"
)

type Data struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func main() {
	env.Load()

	api := civicApi.New(os.Getenv("CIVIC_API_KEY"), os.Getenv("CIVIC_API_ELECTION_ID"), util.MakeRequest)
	res := response.New(api)

	sms := sms.New(os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_TOKEN"), os.Getenv("TWILIO_NUMBER"))

	accessKey := os.Getenv("ACCESS_KEY_ID")
	secretKey := os.Getenv("SECRET_ACCESS_KEY")

	auth := aws.Auth{AccessKey: accessKey, SecretKey: secretKey}
	sqs := sqs.New(auth, aws.USEast)

	queueName := os.Getenv("QUEUE_PREFIX") + "-" + strings.ToLower(os.Getenv("ENVIRONMENT"))

	queue, err := sqs.GetQueue(queueName)
	if err != nil {
		log.Panic("Failed to get queue: ", err)
	}

	for {
		message, getErr := getMessage(queue)
		if getErr != nil {
			fmt.Println(getErr)
			continue
		}

		data := &Data{}
		json.Unmarshal([]byte(message.Body), data)

		msg := res.Generate(data.Message)
		fmt.Println(msg)
		sms.Send(msg, data.Number)

		_, delErr := queue.DeleteMessage(message)
		if delErr != nil {
			fmt.Println(getErr)
			continue
		}
	}
}

func getMessage(queue *sqs.Queue) (*sqs.Message, error) {
	for {
		received, err := queue.ReceiveMessage(1)
		if err != nil {
			return nil, err
		}

		if len(received.Messages) == 0 {
			fmt.Println("Queue is empty, waiting...")
			time.Sleep(5 * time.Second)
		} else {
			return &received.Messages[0], nil
		}
	}
}
