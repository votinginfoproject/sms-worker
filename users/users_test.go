package users

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
	item, firstContact, _, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "en", item["language"])
	assert.Equal(t, true, len(item["last_contact"]) > 0)
	assert.Equal(t, true, firstContact)
}

func TestGetExistingUser(t *testing.T) {
	s := fakeStorage.New()
	u := New(s)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": "0"})
	item, firstContact, lastContactTime, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "es", item["language"])
	assert.Equal(t, false, firstContact)
	assert.Equal(t, "0", lastContactTime)
	assert.NotEqual(t, 0, item["last_contact"])
}

func TestChangeLanguage(t *testing.T) {
	s := fakeStorage.New()
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": "0"})
	u := New(s)
	u.ChangeLanguage("+15551235555", "xx")
	item, firstContact, _, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "xx", item["language"])
	assert.Equal(t, false, firstContact)
}

func TestSetAddress(t *testing.T) {
	s := fakeStorage.New()
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": "0"})
	u := New(s)
	u.SetAddress("+15551235555", "123 test street test city test 12345")
	item, firstContact, _, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "123 test street test city test 12345", item["address"])
	assert.Equal(t, false, firstContact)
}
