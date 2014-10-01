package users

import (
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/votinginfoproject/sms-worker/storage"
)

type Db struct {
	s storage.ExternalStorageService
}

type User struct {
	Data            map[string]string
	FirstContact    bool
	Language        string
	LastContactTime string
}

func New(s storage.ExternalStorageService) *Db {
	return &Db{s}
}

func (u *Db) GetOrCreate(key string) (*User, error) {
	user, err := u.get(key)
	if err != nil {
		switch reflect.TypeOf(err).String() {
		case "*users.updateUserError":
			// got user, but failed to update last_contact time
			return &User{}, err
		default:
			// user not found, create new user
			user, err = u.create(key)
			if err != nil {
				return &User{}, err
			}
		}
	}

	return user, nil
}

func (u *Db) get(key string) (*User, error) {
	item, err := u.s.GetItem(key)
	if err != nil {
		return &User{}, &getUserError{err.Error()}
	}

	time := time.Now().Unix()
	newLastContactTime := strconv.FormatInt(time, 10)
	oldLastContactTime := item["last_contact"]

	err = u.s.UpdateItem(key, map[string]string{"last_contact": newLastContactTime})
	if err != nil {
		log.Printf("[ERROR] unable to update last_contact for user with number: '%s' : %s", key, err)
		return &User{}, &updateUserError{err.Error()}
	}

	return &User{item, false, item["language"], oldLastContactTime}, nil
}

func (u *Db) create(key string) (*User, error) {
	time := time.Now().Unix()
	lastContactTime := strconv.FormatInt(time, 10)

	attrs := map[string]string{"phone_number": key, "language": "en", "last_contact": lastContactTime}

	err := u.s.CreateItem(key, attrs)
	if err != nil {
		log.Printf("[ERROR] unable to create uler with number: '%s' : %s", key, err)
		return &User{}, &createUserError{err.Error()}
	}

	return &User{attrs, true, attrs["language"], lastContactTime}, nil
}

func (u *Db) ChangeLanguage(key, language string) error {
	err := u.s.UpdateItem(key, map[string]string{"language": language})
	if err != nil {
		log.Printf("[ERROR] unable to update language for user with number: '%s' : %s", key, err)
		return err
	}

	return nil
}

func (u *Db) SetAddress(key, address string) error {
	err := u.s.UpdateItem(key, map[string]string{"address": address})
	if err != nil {
		log.Printf("[ERROR] unable to set address for user with number: '%s' : %s", key, err)
		return err
	}

	return nil
}

func (u *User) IsNewUser() bool {
	newUser := false
	if len(u.Data["address"]) == 0 {
		newUser = true
	}

	return newUser
}
