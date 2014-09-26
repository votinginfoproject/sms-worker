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
	userDb   *users.Db
}

func New(civic civicApi.Querier, userDb *users.Db) *Generator {
	rawContent, err := data.Asset("raw/data.yml")
	if err != nil {
		log.Panic("[ERROR] Failed to load responses : ", err)
	}

	content, triggers := responses.Load(rawContent)
	return &Generator{civic, content, triggers, userDb}
}

func (r *Generator) Generate(number string, message string, routine int) []string {
	user, err := r.userDb.GetOrCreate(number)
	if err != nil {
		log.Printf("[ERROR] [%d] User store error : %s", routine, err)
		return []string{r.content.Errors.Text["en"]["generalBackend"]}
	}

	message = strings.TrimSpace(message)
	message = strings.ToLower(message)

	action := r.triggers[user.Language][message]

	if len(action) == 0 {
		success, newLanguage := r.checkIfOtherLanguage(message)
		if success == true {
			user.Language = newLanguage
			action = "ChangeLanguage"
		}
	}

	log.Printf("[INFO] [%d] Taking action '%s'", routine, action)

	lctm := r.lastContactTimeMessage(user)

	messages := r.performAction(action, user, message, routine)

	if len(lctm) > 0 {
		messages = append(messages, lctm)
	}

	return messages
}

func (r *Generator) lastContactTimeMessage(user *users.User) string {
	message := ""

	lcInt, _ := strconv.ParseInt(user.LastContactTime, 10, 64)
	lcTime := time.Unix(lcInt, 0)
	duration := time.Since(lcTime)

	if duration > (7*24*time.Hour) && len(user.Data["address"]) > 0 {
		message = r.content.LastContact.Text[user.Language]["prefix"] + "\n" + user.Data["address"]
	}

	return message
}

func (r *Generator) performAction(action string, user *users.User, message string, routine int) []string {
	var messages []string

	switch action {
	case "Elo":
		messages = r.elo(user.Data["address"], user.Language, user.FirstContact, routine)
	case "Registration":
		messages = r.registration(user.Data["address"], user.Language, user.FirstContact, routine)
	case "Help":
		if user.FirstContact == true {
			messages = []string{r.content.Intro.Text[user.Language]["all"]}
		} else {
			messages = []string{r.content.Help.Text[user.Language]["menu"], r.content.Help.Text[user.Language]["languages"]}
		}
	case "About":
		if user.FirstContact == true {
			messages = []string{r.content.Intro.Text[user.Language]["all"]}
		} else {
			messages = []string{r.content.About.Text[user.Language]["all"]}
		}
	case "Intro":
		messages = []string{r.content.Intro.Text[user.Language]["all"]}
	case "ChangeLanguage":
		messages = r.changeLanguage(user.Data["phone_number"], user.Language)
	case "PollingLocation":
		if len(user.Data["address"]) == 0 && user.FirstContact == true {
			messages = []string{r.content.Intro.Text[user.Language]["all"]}
		} else if len(user.Data["address"]) == 0 && user.FirstContact == false {
			messages = []string{r.content.Errors.Text[user.Language]["needAddress"] + "\n\n" + r.content.Help.Text[user.Language]["languages"]}
		} else {
			messages = r.pollingLocation(user, user.Data["address"], routine)
		}
	default:
		messages = r.pollingLocation(user, message, routine)
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
	err := r.userDb.ChangeLanguage(number, language)
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

func (r *Generator) pollingLocation(user *users.User, message string, routine int) []string {
	newUser := false
	if len(user.Data["address"]) == 0 {
		newUser = true
	}

	res, err := r.civic.Query(message)
	if err != nil {
		log.Printf("[ERROR] [%d] Civic API failure : %s", routine, err)
		return []string{r.content.Errors.Text[user.Data["language"]]["generalBackend"]}
	}

	messages, success := pollingLocation.BuildMessage(res, user.Data["language"], newUser, user.FirstContact, r.content)
	if success == true {
		r.userDb.SetAddress(user.Data["phone_number"], message)
	}

	return messages
}
