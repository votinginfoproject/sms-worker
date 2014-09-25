package pollingLocation

import (
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/responses"
)

func BuildMessage(res *civicApi.Response, language string, newUser bool, firstContact bool, content *responses.Content) ([]string, bool) {
	if len(res.PollingLocations) > 0 {
		return success(res, language, content), true
	} else if len(res.Error.Errors) == 0 && len(res.PollingLocations) == 0 {
		return []string{content.Errors.Text[language]["noElectionInfo"]}, true
	} else {
		return failure(res, language, newUser, firstContact, content), false
	}
}

func success(res *civicApi.Response, language string, content *responses.Content) []string {
	pl := res.PollingLocations[0]
	response := content.PollingLocation.Text[language]["prefix"] + "\n"

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
		response = response + "\n" + content.PollingLocation.Text["en"]["hours"] + " " + pl.PollingHours
	}

	return []string{response + "\n\n" + content.Help.Text[language]["menu"] + " " + content.Help.Text[language]["languages"]}
}

func failure(res *civicApi.Response, language string, newUser bool, firstContact bool, content *responses.Content) []string {
	var reason string
	if len(res.Error.Errors) > 0 {
		reason = res.Error.Errors[0].Reason
	}

	switch reason {
	case "parseError":
		if newUser == true && firstContact == true {
			return []string{content.Intro.Text[language]["all"]}
		} else if newUser == true && firstContact == false {
			return []string{content.Errors.Text[language]["addressParseNewUser"] + "\n\n" + content.Help.Text[language]["languages"]}
		} else {
			return []string{content.Errors.Text[language]["addressParseExistingUser"]}
		}
	case "notFound":
		return []string{content.Errors.Text[language]["noElectionInfo"]}
	default:
		return []string{content.Errors.Text[language]["generalBackend"]}
	}
}
