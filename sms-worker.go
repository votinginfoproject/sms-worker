package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/env"
	"github.com/votinginfoproject/sms-worker/logger"
	"github.com/votinginfoproject/sms-worker/response"
	"github.com/votinginfoproject/sms-worker/sms"
	"github.com/votinginfoproject/sms-worker/util"
	"github.com/yvasiyarov/gorelic"
)

type Data struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func main() {
	env.Load()
	host, _ := os.Hostname()

	if os.Getenv("ENVIRONMENT") == "production" {
		agent := gorelic.NewAgent()
		agent.NewrelicName = "sms-worker" + "-" + host
		agent.NewrelicLicense = os.Getenv("NEWRELIC_TOKEN")
		agent.NewrelicPollInterval = 15
		agent.Run()
	}

	procs, err := strconv.Atoi(os.Getenv("PROCS"))
	if err != nil {
		log.Fatal("[ERROR] you must specify procs in the .env file")
	}
	runtime.GOMAXPROCS(procs)

	routines, err := strconv.Atoi(os.Getenv("ROUTINES"))
	if err != nil {
		log.Fatal("[ERROR] you must specify routines in the .env file")
	}

	log.SetOutput(logger.New())

	api := civicApi.New(os.Getenv("CIVIC_API_KEY"), os.Getenv("CIVIC_API_ELECTION_ID"), util.MakeRequest)
	res := response.New(api)

	sms := sms.New(os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_TOKEN"), os.Getenv("TWILIO_NUMBER"))

	accessKey := os.Getenv("ACCESS_KEY_ID")
	secretKey := os.Getenv("SECRET_ACCESS_KEY")

	auth := aws.Auth{AccessKey: accessKey, SecretKey: secretKey}
	sqs := sqs.New(auth, aws.USEast)

	queueName := os.Getenv("QUEUE_PREFIX") + "-" + strings.ToLower(os.Getenv("ENVIRONMENT"))

	var wg sync.WaitGroup

	for i := 0; i < routines; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, routine int) {
			defer wg.Done()
			queue, err := sqs.GetQueue(queueName)

			if err != nil {
				log.Fatal("[ERROR] Failed to get queue ", err)
			}

			log.Print("[INFO] Started routine ", routine)

			for {
				message, getErr := getMessage(queue)
				if getErr != nil {
					log.Printf("[ERROR] [%d] %s", routine, getErr)
					continue
				}

				data := &Data{}
				log.Printf("[INFO] [%d] Received %s", routine, string(message.Body))
				json.Unmarshal([]byte(message.Body), data)

				msg := res.Generate(data.Message)
				log.Printf("[INFO] [%d] Sending '%s' To %s", routine, msg, data.Number)
				sms.Send(msg, data.Number)

				_, delErr := queue.DeleteMessage(message)
				if delErr != nil {
					log.Printf("[ERROR] [%d] %s", routine, delErr)
					continue
				}
			}
		}(&wg, i)
	}

	wg.Wait()
}

func getMessage(queue *sqs.Queue) (*sqs.Message, error) {
	for {
		received, err := queue.ReceiveMessage(1)
		if err != nil {
			return nil, err
		}

		if len(received.Messages) == 0 {
			time.Sleep(3 * time.Second)
		} else {
			return &received.Messages[0], nil
		}
	}
}
