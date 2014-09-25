package elo

import (
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/responses"
)

func BuildMessage(res *civicApi.Response, language string, messages *responses.Content) []string {
	name, email, phone := getElo(res)
	if len(name) == 0 {
		return []string{messages.Errors.Text[language]["noElectionOfficial"]}
	}

	message := messages.Elo.Text[language]["prefix"] + "\n" + name
	if len(email) > 0 {
		message = message + "\n" + messages.Elo.Text[language]["phone"] + " " + phone
	}

	if len(phone) > 0 {
		message = message + "\n" + messages.Elo.Text[language]["email"] + " " + email
	}

	return []string{message}
}

func getElo(res *civicApi.Response) (string, string, string) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	var name string
	var email string
	var phone string

	elo := res.State[0].LocalJurisdiction.ElectionAdministrationBody.ElectionOfficials[0]

	name = elo.Name
	email = elo.EmailAddress
	phone = elo.OfficePhoneNumber

	return name, email, phone
}
