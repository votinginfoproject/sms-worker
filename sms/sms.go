package sms

import (
	"fmt"
	"net/http"
	"net/url"
)

type ExternalSmsServce interface {
	Send(messages []string, to string) error
}

type Twilio struct {
	endpoint *url.URL
	from     string
}

func New(sid string, token string, from string) *Twilio {
	endpoint, _ := url.Parse(fmt.Sprintf("https://%s:%s@api.twilio.com", sid, token))
	endpoint.Path += fmt.Sprintf("/2010-04-01/Accounts/%s/Messages", sid)

	return &Twilio{endpoint, from}
}

func (t *Twilio) Send(messages []string, to string) error {
	for _, message := range messages {
		data := url.Values{}
		data.Set("From", t.from)
		data.Set("To", to)
		data.Set("Body", message)

		r, err := http.PostForm(t.endpoint.String(), data)
		if err != nil {
			return err
		}

		defer r.Body.Close()
	}

	return nil
}
