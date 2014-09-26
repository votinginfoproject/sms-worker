package users

import (
	"log"
	"strconv"
	"time"

	"github.com/votinginfoproject/sms-worker/storage"
)

type Db struct {
	s storage.ExternalStorageService
}

func New(s storage.ExternalStorageService) *Db {
	return &Db{s}
}

func (u *Db) GetOrCreate(key string) (map[string]string, bool, string, error) {
	isFirstContact := false
	item, getErr := u.s.GetItem(key)
	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)

	// since updating creates an item if one doesn't exist, verify that last_contact
	// is set
	if len(item["last_contact"]) == 0 {
		item["last_contact"] = timeString
	}

	lastContactTime := item["last_contact"]

	if getErr != nil {
		isFirstContact = true
		attrs := map[string]string{"phone_number": key, "language": "en", "last_contact": timeString}
		createErr := u.s.CreateItem(key, attrs)

		if createErr != nil {
			log.Printf("[ERROR] unable to create user with number: '%s' : %s", key, createErr)
			return make(map[string]string), isFirstContact, lastContactTime, createErr
		} else {
			return attrs, isFirstContact, lastContactTime, nil
		}
	}

	timeErr := u.s.UpdateItem(key, map[string]string{"last_contact": timeString})
	if timeErr != nil {
		log.Printf("[ERROR] unable to update last_contact for user with number: '%s' : %s", key, timeErr)
		return make(map[string]string), isFirstContact, lastContactTime, timeErr
	}

	return item, isFirstContact, lastContactTime, nil
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
