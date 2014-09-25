package responseGenerator

import (
	"log"
	"strings"

	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/data"
	"github.com/votinginfoproject/sms-worker/response_generator/elo"
	"github.com/votinginfoproject/sms-worker/response_generator/polling_location"
	"github.com/votinginfoproject/sms-worker/response_generator/registration"
	"github.com/votinginfoproject/sms-worker/responses"
	"github.com/votinginfoproject/sms-worker/users"
)

type Generator struct {
	civic    civicApi.Querier
	content  *responses.Content
	triggers map[string]map[string]string
}

func New(civic civicApi.Querier) *Generator {
	rawContent, err := data.Asset("raw/data.yml")
	if err != nil {
		log.Panic("[ERROR] Failed to load responses : ", err)
	}

	content, triggers := responses.Load(rawContent)
	return &Generator{civic, content, triggers}
}

func (r *Generator) Generate(user *users.Users, number string, message string, routine int) []string {
	userData, firstContact, err := user.GetOrCreate(number)
	if err != nil {
		log.Printf("[ERROR] [%d] User store error : %s", routine, err)
		return []string{r.content.Errors.Text["en"]["generalBackend"]}
	}

	message = strings.TrimSpace(message)
	message = strings.ToLower(message)

	language := userData["language"]
	action := r.triggers[language][message]

	if len(action) == 0 {
		success, newLanguage := r.checkIfOtherLanguage(message)
		language = newLanguage
		if success == true {
			action = "ChangeLanguage"
		}
	}

	switch action {
	case "Elo":
		return r.elo(userData["address"], language, firstContact, routine)
	case "Registration":
		return r.registration(userData["address"], language, firstContact, routine)
	case "Help":
		if firstContact == true {
			return []string{r.content.Intro.Text[language]["all"]}
		} else {
			return []string{r.content.Help.Text[language]["menu"], r.content.Help.Text[language]["languages"]}
		}
	case "About":
		if firstContact == true {
			return []string{r.content.Intro.Text[language]["all"]}
		} else {
			return []string{r.content.About.Text[language]["all"]}
		}
	case "Intro":
		return []string{r.content.Intro.Text[language]["all"]}
	case "ChangeLanguage":
		return r.changeLanguage(user, number, language)
	case "PollingLocation":
		if len(userData["address"]) == 0 && firstContact == true {
			return []string{r.content.Intro.Text[language]["all"]}
		} else if len(userData["address"]) == 0 && firstContact == false {
			return []string{r.content.Errors.Text[language]["needAddress"] + "\n\n" + r.content.Help.Text[language]["languages"]}
		} else {
			return r.pollingLocation(userData, user, number, userData["address"], firstContact, routine)
		}
	default:
		return r.pollingLocation(userData, user, number, message, firstContact, routine)
	}
}

func (r *Generator) checkIfOtherLanguage(message string) (bool, string) {
	for language, _ := range r.triggers {
		if len(r.triggers[language][message]) > 0 {
			return true, language
		}
	}

	return false, ""
}

func (r *Generator) changeLanguage(user *users.Users, number string, language string) []string {
	err := user.ChangeLanguage(number, language)
	if err != nil {
		return []string{r.content.Errors.Text[language]["generalBackend"]}
	}

	return []string{r.content.Help.Text[language]["menu"], r.content.Help.Text[language]["languages"]}
}

func (r *Generator) elo(address string, language string, firstContact bool, routine int) []string {
	if address == "" {
		if firstContact == true {
			return []string{r.content.Intro.Text[language]["all"]}
		} else {
			return []string{r.content.Errors.Text[language]["needAddress"] + "\n\n" + r.content.Help.Text[language]["languages"]}
		}
	}

	res, err := r.civic.Query(address)
	if err != nil {
		log.Printf("[ERROR] [%d] Civic API failure : %s", routine, err)
		return []string{r.content.Errors.Text[language]["generalBackend"]}
	}

	return elo.BuildMessage(res, language, r.content)
}

func (r *Generator) registration(address string, language string, firstContact bool, routine int) []string {
	if address == "" {
		if firstContact == true {
			return []string{r.content.Intro.Text[language]["all"]}
		} else {
			return []string{r.content.Errors.Text[language]["needAddress"] + "\n\n" + r.content.Help.Text[language]["languages"]}
		}
	}

	res, err := r.civic.Query(address)
	if err != nil {
		log.Printf("[ERROR] [%d] Civic API failure : %s", routine, err)
		return []string{r.content.Errors.Text[language]["generalBackend"]}
	}

	return registration.BuildMessage(res, language, r.content)
}

func (r *Generator) pollingLocation(userData map[string]string, user *users.Users, number string, message string, firstContact bool, routine int) []string {
	newUser := false
	if len(userData["address"]) == 0 {
		newUser = true
	}

	res, err := r.civic.Query(message)
	if err != nil {
		log.Printf("[ERROR] [%d] Civic API failure : %s", routine, err)
		return []string{r.content.Errors.Text[userData["language"]]["generalBackend"]}
	}

	messages, success := pollingLocation.BuildMessage(res, userData["language"], newUser, firstContact, r.content)
	if success == true {
		user.SetAddress(number, message)
	}

	return messages
}
