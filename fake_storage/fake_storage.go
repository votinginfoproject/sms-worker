package fakeStorage

import "errors"

type Local struct {
	db map[string]map[string]string
}

func New() *Local {
	return &Local{make(map[string]map[string]string)}
}

func (s *Local) GetItem(key string) (map[string]string, error) {
	attrs := s.db[key]

	if attrs == nil {
		return make(map[string]string), errors.New("Item not found")
	}

	return attrs, nil
}

func (s *Local) CreateItem(key string, attrs map[string]string) error {
	s.db[key] = attrs
	attrs["phone_number"] = key
	return nil
}

func (s *Local) UpdateItem(key string, attrs map[string]string) error {
	if s.db[key] == nil {
		s.db[key] = attrs
	} else {
		for attrKey, value := range attrs {
			s.db[key][attrKey] = value
		}
	}
	return nil
}
