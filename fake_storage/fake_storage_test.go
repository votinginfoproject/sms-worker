package fakeStorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItemSuccess(t *testing.T) {
	db := New()
	db.CreateItem("+12318384215", map[string]string{"one": "two"})
	attrs, err := db.GetItem("+12318384215")
	assert.Nil(t, err)
	assert.Equal(t, "two", attrs["one"])
}

func TestGetItemFailure(t *testing.T) {
	db := New()
	_, err := db.GetItem("+12318384215")
	assert.NotNil(t, err)
}

func TestUpdateExistingItem(t *testing.T) {
	db := New()
	db.CreateItem("+12318384215", map[string]string{"one": "two"})
	db.UpdateItem("+12318384215", map[string]string{"two": "three", "four": "five"})
	attrs, _ := db.GetItem("+12318384215")
	assert.Equal(t, "two", attrs["one"])
	assert.Equal(t, "three", attrs["two"])
	assert.Equal(t, "five", attrs["four"])
}

func TestUpdateNewItem(t *testing.T) {
	db := New()
	db.UpdateItem("+12318384215", map[string]string{"two": "three", "four": "five"})
	attrs, _ := db.GetItem("+12318384215")
	assert.Equal(t, "three", attrs["two"])
	assert.Equal(t, "five", attrs["four"])
}
