package responseGenerator

import (
	"log"

	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/data"
	"github.com/votinginfoproject/sms-worker/response_generator/polling_location"
	"github.com/votinginfoproject/sms-worker/responses"
)

type Generator struct {
	civic    civicApi.Querier
	messages *responses.Responses
	triggers map[string]map[string]string
}

func New(civic civicApi.Querier) *Generator {
	rawMessages, err := data.Asset("raw/data.yml")
	if err != nil {
		log.Panic("[ERROR] Failed to load responses : ", err)
	}

	messages, triggers := responses.Load(rawMessages)
	return &Generator{civic, messages, triggers}
}

func (r *Generator) Generate(message string) []string {
	res, err := r.civic.Query(message)
	if err != nil {
		return []string{r.messages.Errors.Text["en"]["generalBackend"]}
	}

	return pollingLocation.BuildMessage(res, r.messages)
}
