package responseGenerator

import (
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/response_generator/polling_location"
)

type Generator struct {
	civic civicApi.Querier
}

func New(civic civicApi.Querier) *Generator {
	return &Generator{civic}
}

func (r *Generator) Generate(message string) []string {
	res, err := r.civic.Query(message)
	if err != nil {
		return []string{"an error has occurred"}
	}

	return pollingLocation.BuildMessage(res)
}
