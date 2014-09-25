package users

import (
	"log"
	"strconv"
	"time"

	"github.com/votinginfoproject/sms-worker/storage"
)

type Users struct {
	s storage.ExternalStorageService
}

func New(s storage.ExternalStorageService) *Users {
	return &Users{s}
}

func (u *Users) GetOrCreate(key string) (map[string]string, bool, error) {
	isFirstContact := false
	item, getErr := u.s.GetItem(key)
	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)

	// since updating creates an item if one doesn't exist, verify that last_contact
	// is set
	if len(item["last_contact"]) == 0 {
		item["last_contact"] = timeString
	}

	if getErr != nil {
		isFirstContact = true
		attrs := map[string]string{"language": "en", "last_contact": timeString}
		createErr := u.s.CreateItem(key, map[string]string{"language": "en"})

		if createErr != nil {
			log.Printf("[ERROR] unable to create user with number: '%s' : %s", key, createErr)
			return make(map[string]string), isFirstContact, createErr
		} else {
			return attrs, isFirstContact, nil
		}
	}

	timeErr := u.s.UpdateItem(key, map[string]string{"last_contact": timeString})
	if timeErr != nil {
		log.Printf("[ERROR] unable to update last_contact for user with number: '%s' : %s", key, timeErr)
		return make(map[string]string), isFirstContact, timeErr
	}

	return item, isFirstContact, nil
}

func (u *Users) ChangeLanguage(key, language string) error {
	err := u.s.UpdateItem(key, map[string]string{"language": language})
	if err != nil {
		log.Printf("[ERROR] unable to update language for user with number: '%s' : %s", key, err)
		return err
	}

	return nil
}

func (u *Users) SetAddress(key, address string) error {
	err := u.s.UpdateItem(key, map[string]string{"address": address})
	if err != nil {
		log.Printf("[ERROR] unable to set address for user with number: '%s' : %s", key, err)
		return err
	}

	return nil
}
