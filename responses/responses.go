package responses

import (
	"log"

	"github.com/fatih/structs"
	"gopkg.in/yaml.v1"
)

type Response struct {
	Text     map[string]map[string]string "text"
	Triggers map[string][]string          "triggers"
}

type Responses struct {
	PollingLocation  Response "pollingLocation"
	ElectionOfficial Response "electionOfficial"
	Errors           Response "errors"
}

func Load(raw []byte) (*Responses, map[string]map[string]string) {
	r := &Responses{}

	err := yaml.Unmarshal(raw, r)
	if err != nil {
		log.Panic("[ERROR] Failed to parse responses : ", err)
	}

	return r, buildTriggerLookup(r)
}

func buildTriggerLookup(r *Responses) map[string]map[string]string {
	lookup := make(map[string]map[string]string)

	for _, field := range structs.New(r).Fields() {
		name := field.Name()

		triggers := field.Value().(Response).Triggers

		for language, words := range triggers {
			if lookup[language] == nil {
				lookup[language] = make(map[string]string)
			}
			for _, word := range words {
				lookup[language][word] = name
			}
		}
	}

	return lookup
}
