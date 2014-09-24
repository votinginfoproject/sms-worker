package poll

import (
	"log"
	"sync"

	"github.com/votinginfoproject/sms-worker/queue"
	"github.com/votinginfoproject/sms-worker/response_generator"
	"github.com/votinginfoproject/sms-worker/sms"
)

func Start(q queue.ExternalQueueService, rg *responseGenerator.Generator, sms sms.ExternalSmsServce, wg *sync.WaitGroup, routine int) {
	defer wg.Done()

	log.Print("[INFO] Started routine ", routine)

	for {
		message, number, rawMsg, getErr := q.GetMessage(routine)
		if getErr != nil {
			log.Printf("[ERROR] [%d] %s", routine, getErr)
			continue
		}

		reply := rg.Generate(message, routine)

		log.Printf("[INFO] [%d] Sending '%s' To %s", routine, reply, number)

		smsErr := sms.Send(reply, number)
		if smsErr != nil {
			log.Printf("[ERROR] [%d] %s", routine, smsErr)
			continue
		}

		delErr := q.DeleteMessage(rawMsg)
		if delErr != nil {
			log.Printf("[ERROR] [%d] %s", routine, delErr)
			continue
		}
	}
}
