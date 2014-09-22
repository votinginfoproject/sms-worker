package pollingLocation

import "github.com/votinginfoproject/sms-worker/civic_api"

func BuildMessage(res *civicApi.Response) []string {
	if len(res.Error.Errors) == 0 && len(res.PollingLocations) > 0 {
		return success(res)
	} else {
		return failure(res)
	}
}

func success(res *civicApi.Response) []string {
	pl := res.PollingLocations[0]
	response := "Your polling place is:\n"

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
		response = response + "\nHours: " + pl.PollingHours
	}

	return []string{response}
}

func failure(res *civicApi.Response) []string {
	return []string{"the civic api returned an error"}
}
