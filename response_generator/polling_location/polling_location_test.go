package pollingLocation_test

import (
	"fmt"
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

func TestPollingLocationSuccessNewUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{
		"Your polling place is:\nFIRST UNITARIAN CHURCH OF PROVIDENCE - 2ND FLR AUDITORIUM - B\n1 BENEVOLENT ST\nPROVIDENCE, RI 02906\nHours: 7am - 7pm",
		content.Help.Text["en"]["menu"],
		content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}

func TestPollingLocationWithDropOffSuccessNewUser(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccessWithDropOff)
	g := responseGenerator.New(c, u)

	expected := []string{
		"Your polling place is:\nFIRST UNITARIAN CHURCH OF PROVIDENCE - 2ND FLR AUDITORIUM - B\n1 BENEVOLENT ST\nPROVIDENCE, RI 02906\nHours: 7am - 7pm",
		"Your nearest drop box location is:\nPROVIDENCE LIBRARY - GUTENBERG BRANCH\n14 40TH ST\nPROVIDENCE, RI 02906\nHours: 7am - 7pm",
		content.Help.Text["en"]["menu"],
		content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}

func TestPollingLocationSuccessNewUserFirstcontactCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()

	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccessEmpty)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["en"]["all"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "poll", 0))
}

func TestPollingLocationSuccessNewUserCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString})

	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccessEmpty)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["needAddress"] + "\n\n" + content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "poll", 0))
}

func TestPollingLocationSuccessExistingUserCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": timeString, "address": "exists"})

	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{
		fmt.Sprintf("%s\nFIRST UNITARIAN CHURCH OF PROVIDENCE - 2ND FLR AUDITORIUM - B\n1 BENEVOLENT ST\nPROVIDENCE, RI 02906\n%s 7am - 7pm", content.PollingLocation.Text["es"]["prefix"], content.PollingLocation.Text["es"]["hours"]),
		content.Help.Text["es"]["menu"],
		content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "spoll", 0))
}

func TestPollingLocationSuccessWithDropOffExistingUserCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": timeString, "address": "exists"})

	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccessWithDropOff)
	g := responseGenerator.New(c, u)

	expected := []string{
		fmt.Sprintf("%s\nFIRST UNITARIAN CHURCH OF PROVIDENCE - 2ND FLR AUDITORIUM - B\n1 BENEVOLENT ST\nPROVIDENCE, RI 02906\n%s 7am - 7pm", content.PollingLocation.Text["es"]["prefix"], content.PollingLocation.Text["es"]["hours"]),
		fmt.Sprintf("%s\nPROVIDENCE LIBRARY - GUTENBERG BRANCH\n14 40TH ST\nPROVIDENCE, RI 02906\n%s 7am - 7pm", content.DropOffLocation.Text["es"]["prefix"], content.DropOffLocation.Text["es"]["hours"]),
		content.Help.Text["es"]["menu"],
		content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "spoll", 0))
}

func TestPollingLocationSuccessMultiExistingUserCommand(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": timeString, "address": "exists"})

	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccessMulti)
	g := responseGenerator.New(c, u)

	expected := []string{
		fmt.Sprintf("%s\nFIRST UNITARIAN CHURCH OF PROVIDENCE - 2ND FLR AUDITORIUM - B\n1 BENEVOLENT ST\nPROVIDENCE, RI 02906\n%s 7am - 7pm", content.PollingLocation.Text["es"]["prefix"], content.PollingLocation.Text["es"]["hours"]),
		content.PollingLocation.Text["es"]["multi"],
		content.Help.Text["es"]["menu"],
		content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "spoll", 0))
}

func TestPollingLocationSuccessExistingUserNewAddress(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": timeString})

	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccess)
	g := responseGenerator.New(c, u)

	expected := []string{
		fmt.Sprintf("%s\nFIRST UNITARIAN CHURCH OF PROVIDENCE - 2ND FLR AUDITORIUM - B\n1 BENEVOLENT ST\nPROVIDENCE, RI 02906\n%s 7am - 7pm", content.PollingLocation.Text["es"]["prefix"], content.PollingLocation.Text["es"]["hours"]),
		content.Help.Text["es"]["menu"],
		content.Help.Text["es"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
	updatedUser, _ := u.GetOrCreate("+15551235555")
	assert.Equal(t, "111 address street", updatedUser.Data["address"])
}

func TestPollingLocationParseErrorNewUserFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestParseError)
	g := responseGenerator.New(c, u)

	expected := []string{content.Intro.Text["en"]["all"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}

func TestPollingLocationParseErrorNewUserNotFirstContact(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "last_contact": timeString})

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestParseError)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["addressParseNewUser"] + "\n\n" + content.Help.Text["en"]["languages"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}

func TestPollingLocationParseErrorExistingUser(t *testing.T) {
	setup()
	s := fakeStorage.New()

	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)
	s.CreateItem("+15551235555", map[string]string{"language": "en", "address": "123 valid street", "last_contact": timeString})

	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestParseError)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["addressParseExistingUser"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}

func TestPollingLocationNotFoundError(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestNotFoundError)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["noElectionInfo"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}

func TestPollingLocationEmpty(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestSuccessEmpty)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["noElectionInfo"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}

func TestPollingLocationFailure(t *testing.T) {
	setup()
	s := fakeStorage.New()
	u := users.New(s)

	c := civicApi.New("", "", "", civicApiFixtures.MakeRequestFailure)
	g := responseGenerator.New(c, u)

	expected := []string{content.Errors.Text["en"]["generalBackend"]}
	assert.Equal(t, expected, g.Generate("+15551235555", "111 address street", 0))
}
