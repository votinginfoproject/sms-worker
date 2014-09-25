package pollingLocation_test

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

var makeRequestParseError = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../../civic_api/test_data/google_civic_parse_error.json")

	return data, nil
}

var makeRequestNotFoundError = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../../civic_api/test_data/google_civic_not_found_error.json")

	return data, nil
}

var makeRequestFailure = func(endpoint string) ([]byte, error) {
	return nil, errors.New("something bad has happened")
}

func TestPollingLocationSuccessNewUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestSuccess)
	g := responseGenerator.New(c)

	expected := []string{
		"Your polling place is:\nSun Valley Neighborhood Center\n115 W 6th St\nSun Valley, NV 00000\nHours: 7am-7pm",
		content.Help.Text["en"]["menu"] + "\n" + content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
}

func TestPollingLocationSuccessExistingUserCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": timeString})

	u := users.New(s)

	c := civicApi.New("", "", makeRequestSuccess)
	g := responseGenerator.New(c)

	expected := []string{
		"spanish-Your polling place is:\nSun Valley Neighborhood Center\n115 W 6th St\nSun Valley, NV 00000\nHours: 7am-7pm",
		content.Help.Text["es"]["menu"] + "\n" + content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "spoll", 0))
}

func TestPollingLocationSuccessExistingUserNewAddress(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": timeString})

	u := users.New(s)

	c := civicApi.New("", "", makeRequestSuccess)
	g := responseGenerator.New(c)

	expected := []string{
		"spanish-Your polling place is:\nSun Valley Neighborhood Center\n115 W 6th St\nSun Valley, NV 00000\nHours: 7am-7pm",
		content.Help.Text["es"]["menu"] + "\n" + content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
	updatedUser, _, _ := u.GetOrCreate("+15551235555")
	assert.Equal(t, "111 address street", updatedUser["address"])
}

func TestPollingLocationParseErrorNewUserFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestParseError)
	g := responseGenerator.New(c)

	expected := []string{content.Intro.Text["en"]["all"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
}

func TestPollingLocationParseErrorNewUserNotFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString})

	c := civicApi.New("", "", makeRequestParseError)
	g := responseGenerator.New(c)

	expected := []string{content.Errors.Text["en"]["addressParseNewUser"] + "\n\n" + content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
}

func TestPollingLocationParseErrorExistingUser(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "address": "123 valid street", "last_contact": timeString})

	u := users.New(s)

	c := civicApi.New("", "", makeRequestParseError)
	g := responseGenerator.New(c)

	expected := []string{content.Errors.Text["en"]["addressParseExistingUser"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
}

func TestPollingLocationNotFoundError(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestNotFoundError)
	g := responseGenerator.New(c)

	expected := []string{content.Errors.Text["en"]["noElectionInfo"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
}

func TestPollingLocationEmpty(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestSuccessEmpty)
	g := responseGenerator.New(c)

	expected := []string{content.Errors.Text["en"]["noElectionInfo"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
}

func TestPollingLocationFailure(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestFailure)
	g := responseGenerator.New(c)

	expected := []string{content.Errors.Text["en"]["generalBackend"]}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "111 address street", 0))
}
