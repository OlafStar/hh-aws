package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


func (u DynamoDBClient) DoesUserExists(email string) (bool, error) {
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u DynamoDBClient) InsertUser(user types.User) error {
	//asseble item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(user.Email),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}
	//insert
	_, err := u.databaseStore.PutItem(item)

	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) GetUser(email string) (types.User, error) {
	var user types.User

	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u DynamoDBClient) GetAdminUser(email string) (types.User, error) {
	var user types.User

	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ADMIN_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u DynamoDBClient) GetCosmetologistUser(email string) (types.User, error) {
	var user types.User

	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(COSMETOLOGIST_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u DynamoDBClient) InsertCosmetologistUser(user types.User) error {
	//asseble item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(COSMETOLOGIST_TABLE),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(user.Email),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}
	//insert
	_, err := u.databaseStore.PutItem(item)

	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) DoesCosmetologistUserExists(email string) (bool, error) {
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(COSMETOLOGIST_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}