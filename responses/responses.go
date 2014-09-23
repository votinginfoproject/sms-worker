package responses

import (
	"log"

	"gopkg.in/yaml.v1"
)

type Response struct {
	Text     map[string]map[string]string "text"
	Triggers map[string][]string          "triggers"
}

type Responses struct {
	PollingLocation  Response "pollingLocation"
	ElectionOfficial Response "electionOfficial"
}

func Load(raw []byte) *Responses {
	r := &Responses{}

	err := yaml.Unmarshal(raw, r)
	if err != nil {
		log.Panic("[ERROR] Failed to parse responses : ", err)
	}

	return r
}
