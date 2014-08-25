package response

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/votinginfoproject/sms-worker/civic_api"
)

var makeRequestSuccess = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../civic_api/test_data/google_civic_success.json")

	return data, nil
}

var makeRequestError = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("../civic_api/test_data/google_civic_error.json")

	return data, nil
}

var makeRequestFailure = func(endpoint string) ([]byte, error) {
	return nil, errors.New("something bad has happened")
}

func TestResponsePollingLocationSuccess(t *testing.T) {
	c := civicApi.New("", "", makeRequestSuccess)
	r := New(c)

	expected := "Your polling place is:\nSun Valley Neighborhood Center\n115 W 6th St\nSun Valley, NV 00000\nHours: 7am-7pm"
	assert.Equal(t, expected, r.Generate(""))
}

func TestResponsePollingLocationError(t *testing.T) {
	c := civicApi.New("", "", makeRequestError)
	r := New(c)

	expected := "the civic api returned an error"
	assert.Equal(t, expected, r.Generate(""))
}

func TestResponsePollingLocationFailure(t *testing.T) {
	c := civicApi.New("", "", makeRequestFailure)
	r := New(c)

	expected := "an error has occurred"
	assert.Equal(t, expected, r.Generate(""))
}
