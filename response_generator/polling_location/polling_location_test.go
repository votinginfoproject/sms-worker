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
	"github.com/votinginfoproject/sms-worker/fake_storage"
	"github.com/votinginfoproject/sms-worker/response_generator"
	"github.com/votinginfoproject/sms-worker/users"
)

func setup() {
	log.SetOutput(ioutil.Discard)
}

var makeRequestSuccess = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../../civic_api/test_data/google_civic_success.json")

	return data, nil
}

var makeRequestError = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../../civic_api/test_data/google_civic_parse_error.json")

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

	expected := []string{"Your polling place is:\nSun Valley Neighborhood Center\n115 W 6th St\nSun Valley, NV 00000\nHours: 7am-7pm"}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "", 0))
}

func TestPollingLocationSuccessExistingUser(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": timeString})

	u := users.New(s)

	c := civicApi.New("", "", makeRequestSuccess)
	g := responseGenerator.New(c)

	expected := []string{"spanish-Your polling place is:\nSun Valley Neighborhood Center\n115 W 6th St\nSun Valley, NV 00000\nHours: 7am-7pm"}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "", 0))
}

func TestPollingLocationError(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestError)
	g := responseGenerator.New(c)

	expected := []string{"That isn’t a recognized command. Text HELP to see all options."}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "", 0))
}

func TestPollingLocationFailure(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", makeRequestFailure)
	g := responseGenerator.New(c)

	expected := []string{"Sorry, we were unable to find your election day polling location."}
	assert.Equal(t, expected, g.Generate(u, "+15551235555", "", 0))
}
