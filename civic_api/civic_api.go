package civicApi

import (
	"encoding/json"

	"net/url"
)

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
	Error struct {
		Errors []struct {
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
			Message string `json:"message"`
		} `json:"errors"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Querier interface {
	Query(address string) (*Response, error)
}

type CivicApi struct {
	endpoint   *url.URL
	key        string
	electionId string
}

func New(key string, electionId string) *CivicApi {
	endpoint, _ := url.Parse("https://www.googleapis.com/")
	endpoint.Path += "civicinfo/v2/voterinfo"

	return &CivicApi{endpoint, key, electionId}
}

type Requestor func(endpoint string) ([]byte, error)

func (c *CivicApi) Query(address string, makeRequest Requestor) (*Response, error) {
	parameters := url.Values{}
	parameters.Add("key", c.key)
	parameters.Add("electionId", c.electionId)
	parameters.Add("address", address)
	c.endpoint.RawQuery = parameters.Encode()

	body, err := makeRequest(c.endpoint.String())
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
