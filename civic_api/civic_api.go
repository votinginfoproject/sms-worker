package civicApi

import (
	"encoding/json"

	"net/url"
)

type DropOffLocation struct {
	Address struct {
		LocationName string `json:"locationName"`
		Line1        string `json:"line1"`
		City         string `json:"city"`
		State        string `json:"state"`
		Zip          string `json:"zip"`
	} `json:"address"`
	Notes        string `json:"notes"`
	PollingHours string `json:"pollingHours"`
	Sources      []struct {
		Name     string `json:"name"`
		Official bool   `json:"official"`
	} `json:"sources"`
}

type Response struct {
	PollingLocations []struct {
		Address struct {
			LocationName string `json:"locationName"`
			Line1        string `json:"line1"`
			City         string `json:"city"`
			State        string `json:"state"`
			Zip          string `json:"zip"`
		} `json:"address"`
		Notes        string `json:"notes"`
		PollingHours string `json:"pollingHours"`
		Sources      []struct {
			Name     string `json:"name"`
			Official bool   `json:"official"`
		} `json:"sources"`
	} `json:"pollingLocations"`
	DropOffLocations []DropOffLocation `json:"dropOffLocations"`
	Error struct {
		Errors []struct {
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
			Message string `json:"message"`
		} `json:"errors"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	State []struct {
		LocalJurisdiction struct {
			ElectionAdministrationBody struct {
				ElectionOfficials []struct {
					Name              string `json:"name"`
					OfficePhoneNumber string `json:"officePhoneNumber"`
					EmailAddress      string `json:"emailAddress"`
				} `json:"electionOfficials"`
				ElectionRegistrationUrl string `json:"electionRegistrationUrl"`
			} `json:"electionAdministrationBody"`
		} `json:"local_jurisdiction"`
	} `json:"state"`
}

type Querier interface {
	Query(address string) (*Response, error)
}

type requestor func(endpoint string) ([]byte, error)

type CivicApi struct {
	endpoint     *url.URL
	key          string
	electionId   string
	officialOnly string
	makeRequest  requestor
}

func New(key string, electionId string, officialOnly string, makeRequest requestor) *CivicApi {
	endpoint, _ := url.Parse("https://www.googleapis.com/")
	endpoint.Path += "civicinfo/v2/voterinfo"

	return &CivicApi{endpoint, key, electionId, officialOnly, makeRequest}
}

func (c *CivicApi) Query(address string) (*Response, error) {
	parameters := url.Values{}
	parameters.Add("key", c.key)
	if len(c.electionId) > 0 {
		parameters.Add("electionId", c.electionId)
	}
	parameters.Add("address", address)
	parameters.Add("officialOnly", c.officialOnly)
	c.endpoint.RawQuery = parameters.Encode()

	body, err := c.makeRequest(c.endpoint.String())
	if err != nil {
		return nil, err
	}

	res := &Response{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
