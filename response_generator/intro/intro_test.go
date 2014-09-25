package intro_test

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

func TestIntroWithCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequest)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["en"]["all"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "vote", 0))
}
