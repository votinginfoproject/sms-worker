package responses

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadTestData() []byte {
	data, _ := ioutil.ReadFile("test_data/data.yml")

	return data
}

func TestParseYaml(t *testing.T) {
	d, _ := Load(loadTestData())
	assert.Equal(t, "testa", d.PollingLocation.Text["en"]["first"])
	assert.Equal(t, "atest", d.PollingLocation.Text["es"]["first"])
	assert.Equal(t, 2, len(d.PollingLocation.Triggers["en"]))
	assert.Equal(t, 2, len(d.PollingLocation.Triggers["es"]))
}

func TestBuildTriggerLookup(t *testing.T) {
	_, l := Load(loadTestData())
	assert.Equal(t, "PollingLocation", l["es"]["uno"])
	assert.Equal(t, "PollingLocation", l["en"]["one"])
	assert.Equal(t, "ElectionOfficial", l["es"]["tres"])
	assert.Equal(t, "ElectionOfficial", l["en"]["three"])
	assert.Equal(t, "", l["en"]["nope"])
}
