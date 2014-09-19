package queue

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"
)

type ExternalQueueService interface {
	Connect()
	GetMessage(routine int) (string, string, rawMessage, error)
	DeleteMessage(message rawMessage) error
}

type rawMessage interface {
}

type SQS struct {
	q *sqs.Queue
}

type Data struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func New() *SQS {
	return &SQS{nil}
}

func (s *SQS) Connect() {
	accessKey := os.Getenv("ACCESS_KEY_ID")
	secretKey := os.Getenv("SECRET_ACCESS_KEY")

	auth := aws.Auth{AccessKey: accessKey, SecretKey: secretKey}
	sqs := sqs.New(auth, aws.USEast)

	queueName := os.Getenv("QUEUE_PREFIX") + "-" + strings.ToLower(os.Getenv("ENVIRONMENT"))

	queue, err := sqs.GetQueue(queueName)
	if err != nil {
		log.Panic(err)
	}

	s.q = queue
}

func (s *SQS) GetMessage(routine int) (string, string, rawMessage, error) {
	for {
		received, err := s.q.ReceiveMessage(1)

		if err != nil {
			return "", "", nil, err
		}

		if len(received.Messages) == 0 {
			time.Sleep(3 * time.Second)
		} else {
			rawMsg := &received.Messages[0]

			data := &Data{}
			log.Printf("[INFO] [%d] Received %s", routine, rawMsg.Body)
			json.Unmarshal([]byte(rawMsg.Body), data)

			return data.Message, data.Number, rawMsg, nil
		}
	}
}

func (s *SQS) DeleteMessage(message rawMessage) error {
	_, delErr := s.q.DeleteMessage(message.(*sqs.Message))
	if delErr != nil {
		return delErr
	}

	return nil
}
