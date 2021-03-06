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
	"github.com/votinginfoproject/sms-worker/poll"
	"github.com/votinginfoproject/sms-worker/queue"
	"github.com/votinginfoproject/sms-worker/response_generator"
	"github.com/votinginfoproject/sms-worker/sms"
	"github.com/votinginfoproject/sms-worker/storage"
	"github.com/votinginfoproject/sms-worker/users"
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

	st := storage.New()
	user := users.New(st)
	api := civicApi.New(os.Getenv("CIVIC_API_KEY"), os.Getenv("CIVIC_API_ELECTION_ID"), os.Getenv("CIVIC_API_OFFICIAL_ONLY"), util.MakeRequest)
	rg := responseGenerator.New(api, user)

	sms := sms.New(os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_TOKEN"), os.Getenv("TWILIO_NUMBER"))

	q := queue.New()
	q.Connect()

	var wg sync.WaitGroup

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go poll.Start(q, rg, sms, &wg, i)
	}

	wg.Wait()
}
