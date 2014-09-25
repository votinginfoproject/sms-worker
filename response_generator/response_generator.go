package responseGenerator

import (
	"log"
	"strconv"
	"strings"
	"time"

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
	user     *users.Users
}

func New(civic civicApi.Querier, user *users.Users) *Generator {
	rawContent, err := data.Asset("raw/data.yml")
	if err != nil {
		log.Panic("[ERROR] Failed to load responses : ", err)
	}

	content, triggers := responses.Load(rawContent)
	return &Generator{civic, content, triggers, user}
}

func (r *Generator) Generate(number string, message string, routine int) []string {
	userData, firstContact, lastContactTime, err := r.user.GetOrCreate(number)
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

	log.Printf("[INFO] [%d] Taking action '%s'", routine, action)
	return r.checkLastContactTime(r.performAction(action, userData, language, message, firstContact, routine), userData, language, lastContactTime)
}

func (r *Generator) checkLastContactTime(messages []string, userData map[string]string, language string, lastContactTime string) []string {
	lcInt, _ := strconv.ParseInt(lastContactTime, 10, 64)
	lcTime := time.Unix(lcInt, 0)
	duration := time.Since(lcTime)

	if duration > (7*24*time.Hour) && len(userData["address"]) > 0 {
		messages = append(messages, r.content.LastContact.Text[language]["prefix"]+"\n"+userData["address"])
	}

	return messages
}

func (r *Generator) performAction(action string, userData map[string]string, language string, message string, firstContact bool, routine int) []string {
	var messages []string

	switch action {
	case "Elo":
		messages = r.elo(userData["address"], language, firstContact, routine)
	case "Registration":
		messages = r.registration(userData["address"], language, firstContact, routine)
	case "Help":
		if firstContact == true {
			messages = []string{r.content.Intro.Text[language]["all"]}
		} else {
			messages = []string{r.content.Help.Text[language]["menu"], r.content.Help.Text[language]["languages"]}
		}
	case "About":
		if firstContact == true {
			messages = []string{r.content.Intro.Text[language]["all"]}
		} else {
			messages = []string{r.content.About.Text[language]["all"]}
		}
	case "Intro":
		messages = []string{r.content.Intro.Text[language]["all"]}
	case "ChangeLanguage":
		messages = r.changeLanguage(userData["phone_number"], language)
	case "PollingLocation":
		if len(userData["address"]) == 0 && firstContact == true {
			messages = []string{r.content.Intro.Text[language]["all"]}
		} else if len(userData["address"]) == 0 && firstContact == false {
			messages = []string{r.content.Errors.Text[language]["needAddress"] + "\n\n" + r.content.Help.Text[language]["languages"]}
		} else {
			messages = r.pollingLocation(userData, userData["address"], firstContact, routine)
		}
	default:
		messages = r.pollingLocation(userData, message, firstContact, routine)
	}

	return messages
}

func (r *Generator) checkIfOtherLanguage(message string) (bool, string) {
	for language, _ := range r.triggers {
		if len(r.triggers[language][message]) > 0 {
			return true, language
		}
	}

	return false, ""
}

func (r *Generator) changeLanguage(number string, language string) []string {
	err := r.user.ChangeLanguage(number, language)
	if err != nil {
		return []string{r.content.Errors.Text[language]["generalBackend"]}
	}

	return []string{r.content.Help.Text[language]["menu"], r.content.Help.Text[language]["languages"]}
}

func (r *Generator) elo(address string, language string, firstContact bool, routine int) []string {
	if len(address) == 0 {
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
	if len(address) == 0 {
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

func (r *Generator) pollingLocation(userData map[string]string, message string, firstContact bool, routine int) []string {
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
		r.user.SetAddress(userData["phone_number"], message)
	}

	return messages
}
