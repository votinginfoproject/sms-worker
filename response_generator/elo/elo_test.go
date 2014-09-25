package elo_test

import (
	"errors"
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

var makeRequestSuccess = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../../civic_api/test_data/google_civic_success.json")

	return data, nil
}

var makeRequestSuccessEmpty = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../../civic_api/test_data/google_civic_success_empty.json")

	return data, nil
}

var makeRequestFailure = func(endpoint string) ([]byte, error) {
	return nil, errors.New("something bad has happened")
}

func TestEloFailureNewUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestSuccess)
	g := responseGenerator.New(c)

	expected := []string{content.Errors.Text["en"]["needAddress"] + "\n\n" + content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "elo", 0))
}

func TestEloFailureExistingUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString, "address": "real"})

	c := civicApi.New("", "", makeRequestSuccessEmpty)
	g := responseGenerator.New(c)

	expected := []string{content.Errors.Text["en"]["noElectionOfficial"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "elo", 0))
}

func TestEloSuccessExistingUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString, "address": "real"})

	c := civicApi.New("", "", makeRequestSuccess)
	g := responseGenerator.New(c)

	expected := []string{"Your local election official is:\nDan Burk\nPhone: (775) 328-3670\nEmail: dburk@washoecounty.us"}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "elo", 0))
}
