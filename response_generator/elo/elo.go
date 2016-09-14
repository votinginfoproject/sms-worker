package elo

import (
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/responses"
)

func BuildMessage(res *civicApi.Response, language string, content *responses.Content) []string {
	name, email, phone, url := getElo(res)
	if len(name) == 0 {
		return []string{content.Errors.Text[language]["noElectionOfficial"]}
	}

	message := content.Elo.Text[language]["prefix"] + "\n" + name
	if len(phone) > 0 {
		message = message + "\n" + content.Elo.Text[language]["phone"] + " " + phone
	}

	if len(email) > 0 {
		message = message + "\n" + content.Elo.Text[language]["email"] + " " + email
	}

	if len(url) > 0 {
		message = message + "\n" + url
	}

	return []string{message}
}

func getElo(res *civicApi.Response) (string, string, string, string) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	var name string
	var email string
	var phone string
	var url string

	eab := res.State[0].LocalJurisdiction.ElectionAdministrationBody
	elo := eab.ElectionOfficials[0]

	name = elo.Name
	email = elo.EmailAddress
	phone = elo.OfficePhoneNumber
	url = eab.ElectionInfoUrl

	return name, email, phone, url
}
