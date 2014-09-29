package changeLanguage

import (
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/data"
	"github.com/votinginfoproject/sms-worker/fake_storage"
	"github.com/votinginfoproject/sms-worker/response_generator"
	"github.com/votinginfoproject/sms-worker/responses"
	"github.com/votinginfoproject/sms-worker/users"
)

func setup() {
	log.SetOutput(ioutil.Discard)
}

func getContent() *responses.Content {
	rawContent, _ := data.Asset("raw/data.yml")
	content, _ := responses.Load(rawContent)
	return content
}

var content = getContent()

var makeRequest = func(endpoint string) ([]byte, error) {
	return []byte{}, nil
}

func TestChangeLanguageWithLanguageCommandNotFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString})

	c := civicApi.New("", "", makeRequest)
	g := responseGenerator.New(c, u)

	expected := []string{content.Help.Text["es"]["menu"], content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "español", 0))
}

func TestChangeLanguageWithOtherCommandNotFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString})

	c := civicApi.New("", "", makeRequest)
	g := responseGenerator.New(c, u)

	expected := []string{content.Help.Text["es"]["menu"], content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "encuesta", 0))
}

func TestChangeLanguageWithLanguageCommandFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequest)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["es"]["all"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "español", 0))
}

func TestChangeLanguageWithOtherCommandFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequest)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["es"]["all"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "encuesta", 0))
}
