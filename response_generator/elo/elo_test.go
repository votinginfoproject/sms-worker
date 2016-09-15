package elo_test

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

func TestEloFailureNewUserFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["en"]["all"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "elo", 0))
}

func TestEloFailureNewUserNotFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString})

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["needAddress"] + "\n\n" + content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "elo", 0))
}

func TestEloFailureExistingUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString, "address": "real"})

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccessEmpty)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["noElectionOfficial"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "elo", 0))
}

func TestEloSuccessExistingUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString, "address": "real"})

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{"Your local election official is:\nDan Burk\nPhone: (775) 328-3670\nEmail: dburk@washoecounty.us\nhttp://www.sos.ri.gov"}
	assert.Equal(t, expected, g.Generate("+15551235555", "elo", 0))
}
