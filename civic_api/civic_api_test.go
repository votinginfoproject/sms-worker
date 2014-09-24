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

var makeRequestSuccessEmpty = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("test_data/google_civic_success_empty.json")

	return data, nil
}

var makeRequestSuccessEmptyState = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("test_data/google_civic_success_empty_state.json")

	return data, nil
}

var makeRequestError = func(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile("test_data/google_civic_parse_error.json")

	return data, nil
}

func TestQuerySuccess(t *testing.T) {
	c := New("", "", makeRequestSuccess)
	res, _ := c.Query("")
	assert.Equal(t, 0, len(res.Error.Errors), 0)
	assert.Equal(t, "115 W 6th St", res.PollingLocations[0].Address.Line1)
	assert.Equal(t, "http://nvsos.gov/index.aspx?page=703", res.State[0].ElectionAdministrationBody.ElectionRegistrationUrl)
	assert.Equal(t, "Dan Burk", res.State[0].LocalJurisdiction.ElectionAdministrationBody.ElectionOfficials[0].Name)
}

func TestQuerySuccessEmpty(t *testing.T) {
	c := New("", "", makeRequestSuccessEmpty)
	res, _ := c.Query("")
	assert.Equal(t, 0, len(res.Error.Errors))
	assert.Equal(t, 0, len(res.PollingLocations))
	assert.Equal(t, 0, len(res.State))
}

func TestQuerySuccessEmptyState(t *testing.T) {
	c := New("", "", makeRequestSuccessEmptyState)
	res, _ := c.Query("")
	assert.Equal(t, 0, len(res.Error.Errors))
	assert.Equal(t, 1, len(res.PollingLocations))
	assert.Equal(t, 1, len(res.State))
	assert.Equal(t, "", res.State[0].ElectionAdministrationBody.ElectionRegistrationUrl)
}

func TestQueryError(t *testing.T) {
	c := New("", "", makeRequestError)
	res, _ := c.Query("")
	assert.Equal(t, 1, len(res.Error.Errors))
	assert.Equal(t, "parseError", res.Error.Errors[0].Reason)
}
