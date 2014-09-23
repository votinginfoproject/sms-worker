package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/votinginfoproject/sms-worker/fake_storage"
)

func TestNew(t *testing.T) {
	u := New(fakeStorage.New())
	assert.NotNil(t, u)
}

func TestGetNewUser(t *testing.T) {
	s := fakeStorage.New()
	u := New(s)
	item, err := u.GetOrCreate("+12318384215")
	assert.Nil(t, err)
	assert.Equal(t, "en", item["language"])
	assert.Equal(t, true, len(item["last_contact"]) > 0)
}

func TestGetExistingUser(t *testing.T) {
	s := fakeStorage.New()
	u := New(s)
	s.CreateItem("+12318384215", map[string]string{"language": "es", "last_contact": "0"})
	item, err := u.GetOrCreate("+12318384215")
	assert.Nil(t, err)
	assert.Equal(t, "es", item["language"])
	assert.NotEqual(t, 0, item["last_contact"])
}

func TestChangeLanguage(t *testing.T) {
	s := fakeStorage.New()
	s.CreateItem("+12318384215", map[string]string{"language": "es", "last_contact": "0"})
	u := New(s)
	u.ChangeLanguage("+12318384215", "xx")
	item, err := u.GetOrCreate("+12318384215")
	assert.Nil(t, err)
	assert.Equal(t, "xx", item["language"])
}
