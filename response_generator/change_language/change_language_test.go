package changeLanguage

import (
	"io/ioutil"
	"log"
	"testing"

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

func TestChangeLanguageWithLanguageCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequest)
	g := responseGenerator.New(c)

	expected := []string{content.Help.Text["es"]["menu"], content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "español", 0))
}

func TestChangeLanguageWithOtherCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequest)
	g := responseGenerator.New(c)

	expected := []string{content.Help.Text["es"]["menu"], content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "spoll", 0))
}
