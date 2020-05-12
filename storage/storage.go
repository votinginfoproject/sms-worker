package storage

import (
	"os"
	"strings"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/dynamodb"
)

type ExternalStorageService interface {
	GetItem(key string) (map[string]string, error)
	CreateItem(key string, attrs map[string]string) error
	UpdateItem(key string, attrs map[string]string) error
}

type Dynamo struct {
	db *dynamodb.Table
}

func New() *Dynamo {
	accessKey := os.Getenv("ACCESS_KEY_ID")
	secretKey := os.Getenv("SECRET_ACCESS_KEY")

	auth := aws.Auth{AccessKey: accessKey, SecretKey: secretKey}
	server := dynamodb.Server{auth, aws.USWest2}

	primary := dynamodb.NewStringAttribute("phone_number", "")
	key := dynamodb.PrimaryKey{primary, nil}

	tableName := os.Getenv("DB_PREFIX") + "-" + strings.ToLower(os.Getenv("ENVIRONMENT"))

	return &Dynamo{server.NewTable(tableName, key)}
}

func (s *Dynamo) GetItem(key string) (map[string]string, error) {
	item, err := s.db.GetItemConsistent(&dynamodb.Key{HashKey: key}, true)
	if err != nil {
		return make(map[string]string), err
	}

	attrs := make(map[string]string)
	for key, dynamoAttr := range item {
		attrs[key] = dynamoAttr.Value
	}

	return attrs, nil
}

func (s *Dynamo) CreateItem(key string, attrs map[string]string) error {
	dynamoAttrs := []dynamodb.Attribute{}
	for key, value := range attrs {
		dynamoAttrs = append(dynamoAttrs, *dynamodb.NewStringAttribute(key, value))
	}

	_, err := s.db.PutItem(key, "", dynamoAttrs)
	return err
}

func (s *Dynamo) UpdateItem(key string, attrs map[string]string) error {
	dynamoAttrs := []dynamodb.Attribute{}
	for key, value := range attrs {
		dynamoAttrs = append(dynamoAttrs, *dynamodb.NewStringAttribute(key, value))
	}

	_, err := s.db.UpdateAttributes(&dynamodb.Key{HashKey: key}, dynamoAttrs)
	return err
}
