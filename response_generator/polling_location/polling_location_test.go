package pollingLocation_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/response_generator"
)

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

func TestPollingLocationSuccess(t *testing.T) {
	c := civicApi.New("", "", makeRequestSuccess)
	g := responseGenerator.New(c)

	expected := []string{"Your polling place is:\nSun Valley Neighborhood Center\n115 W 6th St\nSun Valley, NV 00000\nHours: 7am-7pm"}
	assert.Equal(t, expected, g.Generate(""))
}

func TestPollingLocationError(t *testing.T) {
	c := civicApi.New("", "", makeRequestError)
	g := responseGenerator.New(c)

	expected := []string{"the civic api returned an error"}
	assert.Equal(t, expected, g.Generate(""))
}

func TestPollingLocationFailure(t *testing.T) {
	c := civicApi.New("", "", makeRequestFailure)
	g := responseGenerator.New(c)

	expected := []string{"an error has occurred"}
	assert.Equal(t, expected, g.Generate(""))
}
