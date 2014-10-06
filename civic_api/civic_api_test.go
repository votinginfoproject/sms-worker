package civicApi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/votinginfoproject/sms-worker/civic_api/fixtures"
)

func TestQuerySuccess(t *testing.T) {
	c := New("", "", "", civicApiFixtures.MakeRequestSuccess)
	res, e := c.Query("")

	fmt.Println(res)
	fmt.Println(e)
	assert.Equal(t, 0, len(res.Error.Errors), 0)
	assert.Equal(t, "1 BENEVOLENT ST", res.PollingLocations[0].Address.Line1)
	assert.Equal(t, "http://www.sos.ri.gov/elections/voters/register/", res.State[0].LocalJurisdiction.ElectionAdministrationBody.ElectionRegistrationUrl)
	assert.Equal(t, "Dan Burk", res.State[0].LocalJurisdiction.ElectionAdministrationBody.ElectionOfficials[0].Name)
}

func TestQuerySuccessEmpty(t *testing.T) {
	c := New("", "", "", civicApiFixtures.MakeRequestSuccessEmpty)
	res, _ := c.Query("")

	assert.Equal(t, 0, len(res.Error.Errors))
	assert.Equal(t, 0, len(res.PollingLocations))
	assert.Equal(t, 0, len(res.State))
}

func TestQuerySuccessEmptyState(t *testing.T) {
	c := New("", "", "", civicApiFixtures.MakeRequestSuccessEmptyState)
	res, _ := c.Query("")

	assert.Equal(t, 0, len(res.Error.Errors))
	assert.Equal(t, 1, len(res.PollingLocations))
	assert.Equal(t, 1, len(res.State))
	assert.Equal(t, "", res.State[0].LocalJurisdiction.ElectionAdministrationBody.ElectionRegistrationUrl)
}

func TestQueryError(t *testing.T) {
	c := New("", "", "", civicApiFixtures.MakeRequestParseError)
	res, _ := c.Query("")

	assert.Equal(t, 1, len(res.Error.Errors))
	assert.Equal(t, "parseError", res.Error.Errors[0].Reason)
}
