package pollingLocation

import (
	"fmt"

	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/responses"
	"github.com/votinginfoproject/sms-worker/users"
)

func BuildMessage(res *civicApi.Response, user *users.User, content *responses.Content) ([]string, bool) {
	if len(res.PollingLocations) > 0 {
		return success(res, user.Language, content), true
	} else if len(res.Error.Errors) == 0 && len(res.PollingLocations) == 0 {
		return []string{content.Errors.Text[user.Language]["noElectionInfo"]}, true
	} else {
		return failure(res, user, content), false
	}
}

func success(res *civicApi.Response, language string, content *responses.Content) []string {
	pl := res.PollingLocations[0]
	response := content.PollingLocation.Text[language]["prefix"] + "\n"

	if len(pl.Address.LocationName) > 0 {
		response = response + pl.Address.LocationName + "\n"
	}

	if len(pl.Address.Line1) > 0 {
		response = fmt.Sprintf("%s%s\n%s, %s %s", response, pl.Address.Line1, pl.Address.City, pl.Address.State, pl.Address.Zip)
	}

	if len(pl.PollingHours) > 0 {
		response = response + "\n" +
			content.PollingLocation.Text[language]["hours"] +
			" " + pl.PollingHours
	}

	if len(res.PollingLocations) > 1 {
		return []string{response, content.PollingLocation.Text[language]["multi"], content.Help.Text[language]["menu"], content.Help.Text[language]["languages"]}
	} else {
		return []string{response, content.Help.Text[language]["menu"], content.Help.Text[language]["languages"]}
	}
}

func failure(res *civicApi.Response, user *users.User, content *responses.Content) []string {
	var reason string
	if len(res.Error.Errors) > 0 {
		reason = res.Error.Errors[0].Reason
	}

	switch reason {
	case "parseError":
		if user.IsNewUser() == true && user.FirstContact == true {
			return []string{content.Intro.Text[user.Language]["all"]}
		} else if user.IsNewUser() == true && user.FirstContact == false {
			return []string{
				content.Errors.Text[user.Language]["addressParseNewUser"] +
					"\n\n" + content.Help.Text[user.Language]["languages"]}
		} else {
			return []string{content.Errors.Text[user.Language]["addressParseExistingUser"]}
		}
	case "notFound":
		return []string{content.Errors.Text[user.Language]["noElectionInfo"]}
	default:
		return []string{content.Errors.Text[user.Language]["generalBackend"]}
	}
}
