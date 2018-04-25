package civicApiFixtures

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var root string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	root = filepath.Dir(filename)
}

func MakeRequestSuccess(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_success.json"))

	return data, nil
}

func MakeRequestSuccessMulti(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_success_multi.json"))

	return data, nil
}

func MakeRequestSuccessWithDropOff(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_success_with_drop_off.json"))

	return data, nil
}

func MakeRequestSuccessNoPollingLocationsWithDropOff(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_success_no_polling_location_with_drop_off.json"))

	return data, nil
}

func MakeRequestSuccessEmpty(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_success_empty.json"))

	return data, nil
}

func MakeRequestSuccessEmptyState(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_success_empty_state.json"))

	return data, nil
}

func MakeRequestParseError(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_parse_error.json"))

	return data, nil
}

func MakeRequestNotFoundError(endpoint string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filepath.Join(root, "google_civic_not_found_error.json"))

	return data, nil
}

func MakeRequestFailure(endpoint string) ([]byte, error) {
	return nil, errors.New("something bad has happened")
}

func MakeRequestSuccessFake(endpoint string) ([]byte, error) {
	return []byte{}, nil
}
