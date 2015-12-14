package dropOffLocation

import (
	"fmt"

	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/responses"
	"github.com/votinginfoproject/sms-worker/users"
)

func BuildMessage(res *civicApi.Response, user *users.User, content *responses.Content) ([]string, bool) {
	if len(res.DropOffLocations) > 0 {
		return success(res, user.Language, content), true
	} else if len(res.Error.Errors) == 0 && len(res.DropOffLocations) == 0 {
		return []string{content.Errors.Text[user.Language]["noElectionInfo"]}, true
	} else {
		return failure(res, user, content), false
	}
}

func success(res *civicApi.Response, language string, content *responses.Content) []string {
	dol := res.DropOffLocations[0]
	response := content.DropOffLocation.Text[language]["prefix"] + "\n"

	if len(dol.Address.LocationName) > 0 {
		response = response + dol.Address.LocationName + "\n"
	}

	if len(dol.Address.Line1) > 0 {
		response = fmt.Sprintf("%s%s\n%s, %s %s", response, dol.Address.Line1, dol.Address.City, dol.Address.State, dol.Address.Zip)
	}

	if len(dol.PollingHours) > 0 {
		response = response + "\n" +
			content.DropOffLocation.Text[language]["hours"] +
			" " + dol.PollingHours
	}

	return []string{response, content.Help.Text[language]["menu"], content.Help.Text[language]["languages"]}
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
