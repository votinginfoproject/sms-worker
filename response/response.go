package response

import (
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/response/polling_location"
)

type Response struct {
	civic civicApi.Querier
}

func New(civic civicApi.Querier) *Response {
	return &Response{civic}
}

func (r *Response) Generate(message string) string {
	res, err := r.civic.Query(message)
	if err != nil {
		return "an error has occurred"
	}

	return pollingLocation.BuildMessage(res)
}
