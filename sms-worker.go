package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/env"
	"github.com/votinginfoproject/sms-worker/logger"
	"github.com/votinginfoproject/sms-worker/queue"
	"github.com/votinginfoproject/sms-worker/response"
	"github.com/votinginfoproject/sms-worker/sms"
	"github.com/votinginfoproject/sms-worker/util"
	"github.com/yvasiyarov/gorelic"
)

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

	queue := queue.New()
	queue.Connect()

	var wg sync.WaitGroup

	for i := 0; i < routines; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, routine int) {
			defer wg.Done()

			log.Print("[INFO] Started routine ", routine)

			for {
				message, number, rawMsg, getErr := queue.GetMessage(routine)
				if getErr != nil {
					log.Printf("[ERROR] [%d] %s", routine, getErr)
					continue
				}

				reply := res.Generate(message)

				log.Printf("[INFO] [%d] Sending '%s' To %s", routine, reply, number)

				sms.Send(reply, number)

				delErr := queue.DeleteMessage(rawMsg)
				if delErr != nil {
					log.Printf("[ERROR] [%d] %s", routine, delErr)
					continue
				}
			}
		}(&wg, i)
	}

	wg.Wait()
}
