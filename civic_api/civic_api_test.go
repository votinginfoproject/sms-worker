package civicApi

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var makeRequestSuccess = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("test_data/google_civic_success.json")

	return data, nil
}

var makeRequestError = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("test_data/google_civic_error.json")

	return data, nil
}

func TestQuerySuccess(t *testing.T) {
	c := New("", "", makeRequestSuccess)
	res, _ := c.Query("")
	assert.Equal(t, len(res.Error.Errors), 0)
	assert.Equal(t, res.PollingLocations[0].Address.Line1, "115 W 6th St")
}

func TestQueryError(t *testing.T) {
	c := New("", "", makeRequestError)
	res, _ := c.Query("")
	assert.Equal(t, len(res.Error.Errors), 1)
	assert.Equal(t, res.Error.Errors[0].Reason, "keyInvalid")
}
