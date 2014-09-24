package user

import (
	"log"
	"strconv"
	"time"

	"github.com/votinginfoproject/sms-worker/storage"
)

type User struct {
	s storage.ExternalStorageService
}

func New(s storage.ExternalStorageService) *User {
	return &User{s}
}

func (u *User) GetOrCreate(key string) (map[string]string, error) {
	item, getErr := u.s.GetItem(key)
	time := time.Now().Unix()
	timeString := strconv.FormatInt(time, 10)

	// since updating creates an item if one doesn't exist, verify that last_contact
	// is set
	if len(item["last_contact"]) == 0 {
		item["last_contact"] = timeString
	}

	if getErr != nil {
		attrs := map[string]string{"language": "en", "last_contact": timeString}
		createErr := u.s.CreateItem(key, map[string]string{"language": "en"})

		if createErr != nil {
			log.Printf("[ERROR] unable to create user with number: '%s' : %s", key, createErr)
			return make(map[string]string), createErr
		} else {
			return attrs, nil
		}
	}

	timeErr := u.s.UpdateItem(key, map[string]string{"last_contact": timeString})
	if timeErr != nil {
		log.Printf("[ERROR] unable to update last_contact for user with number: '%s' : %s", key, timeErr)
		return make(map[string]string), timeErr
	}

	return item, nil
}

func (u *User) ChangeLanguage(key, language string) error {
	err := u.s.UpdateItem(key, map[string]string{"language": language})
	if err != nil {
		log.Printf("[ERROR] unable to update language for user with number: '%s' : %s", key, err)
		return err
	}

	return nil
}

func (u *User) SetAddress(key, address string) error {
	err := u.s.UpdateItem(key, map[string]string{"address": address})
	if err != nil {
		log.Printf("[ERROR] unable to set address for user with number: '%s' : %s", key, err)
		return err
	}

	return nil
}
