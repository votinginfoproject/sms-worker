package lastContact_test

import (
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/votinginfoproject/sms-worker/civic_api"
	"github.com/votinginfoproject/sms-worker/civic_api/fixtures"
	"github.com/votinginfoproject/sms-worker/fake_storage"
	"github.com/votinginfoproject/sms-worker/response_generator"
	"github.com/votinginfoproject/sms-worker/test_helpers"
	"github.com/votinginfoproject/sms-worker/users"
)

func setup() {
	log.SetOutput(ioutil.Discard)
}

var content = testHelpers.GetContent()

func TestSendAddress(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Add(-1 * 8 * 24 * time.Hour).Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString, "address": "test"})

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["en"]["all"], content.LastContact.Text["en"]["prefix"] + "\ntest"}
	assert.Equal(t, expected, g.Generate("+15551235555", "vote", 0))
}

func TestDontSendAddress(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString, "address": "test"})

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["en"]["all"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "vote", 0))
}
