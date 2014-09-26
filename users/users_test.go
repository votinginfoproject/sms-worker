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
	user, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "en", user.Data["language"])
	assert.Equal(t, true, len(user.LastContactTime) > 0)
	assert.Equal(t, true, user.FirstContact)
}

func TestGetExistingUser(t *testing.T) {
	s := fakeStorage.New()
	u := New(s)
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": "0"})
	user, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "es", user.Data["language"])
	assert.Equal(t, false, user.FirstContact)
	assert.Equal(t, "0", user.LastContactTime)
	assert.NotEqual(t, 0, user.Data["last_contact"])
}

func TestChangeLanguage(t *testing.T) {
	s := fakeStorage.New()
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": "0"})
	u := New(s)
	u.ChangeLanguage("+15551235555", "xx")
	user, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "xx", user.Data["language"])
	assert.Equal(t, false, user.FirstContact)
	assert.Equal(t, true, user.IsNewUser())
}

func TestSetAddress(t *testing.T) {
	s := fakeStorage.New()
	s.CreateItem("+15551235555", map[string]string{"language": "es", "last_contact": "0"})
	u := New(s)
	u.SetAddress("+15551235555", "123 test street test city test 12345")
	user, err := u.GetOrCreate("+15551235555")
	assert.Nil(t, err)
	assert.Equal(t, "123 test street test city test 12345", user.Data["address"])
	assert.Equal(t, false, user.FirstContact)
	assert.Equal(t, false, user.IsNewUser())
}
