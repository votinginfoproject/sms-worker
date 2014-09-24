package pollingLocation

import (
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/responses"
)

func BuildMessage(res *civicApi.Response, messages *responses.Responses) []string {
	if len(res.Error.Errors) == 0 && len(res.PollingLocations) > 0 {
		return success(res, messages)
	} else {
		return failure(res, messages)
	}
}

func success(res *civicApi.Response, messages *responses.Responses) []string {
	pl := res.PollingLocations[0]
	response := messages.PollingLocation.Text["en"]["prefix"] + "\n"

	if len(pl.Address.LocationName) > 0 {
		response = response + pl.Address.LocationName + "\n"
	}

	if len(pl.Address.Line1) > 0 {
		response = response + pl.Address.Line1 + "\n"
		response = response + pl.Address.City + ", "
		response = response + pl.Address.State + " "
		response = response + pl.Address.Zip
	}

	if len(pl.PollingHours) > 0 {
		response = response + "\n" + messages.PollingLocation.Text["en"]["hours"] + " " + pl.PollingHours
	}

	return []string{response}
}

func failure(res *civicApi.Response, messages *responses.Responses) []string {
	if len(res.Error.Errors) > 0 {
		if res.Error.Errors[0].Reason == "parseError" {
			return []string{messages.Errors.Text["en"]["addressParse"]}
		}
	}

	return []string{messages.Errors.Text["en"]["generalBackend"]}
}
