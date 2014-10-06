package testHelpers

import (
	"github.com/votinginfoproject/sms-worker/data"
	"github.com/votinginfoproject/sms-worker/responses"
)

func GetContent() *responses.Content {
	rawContent, _ := data.Asset("raw/data.yml")
	content, _ := responses.Load(rawContent)
	return content
}
