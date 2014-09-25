package registration

import (
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/responses"
)

func BuildMessage(res *civicApi.Response, language string, content *responses.Content) []string {
	url := getRegistrationUrl(res)
	if len(url) == 0 {
		return []string{content.Errors.Text[language]["noRegistrationInfo"]}
	}

	return []string{content.Registration.Text[language]["prefix"] + " " + url}
}

func getRegistrationUrl(res *civicApi.Response) string {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	var url string
	url = res.State[0].ElectionAdministrationBody.ElectionRegistrationUrl

	return url
}
