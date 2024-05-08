package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


func (u DynamoDBClient) DoesUserExist(identifier, identifierType string) (bool, error) {
	var queryInput *dynamodb.QueryInput

	if identifierType == "email" {
			queryInput = &dynamodb.QueryInput{
					TableName: aws.String(USER_TABLE),
					IndexName: aws.String("EmailIndex"),
					KeyConditions: map[string]*dynamodb.Condition{
							"email": {
									ComparisonOperator: aws.String("EQ"),
									AttributeValueList: []*dynamodb.AttributeValue{
											{
													S: aws.String(identifier),
											},
									},
							},
					},
			}
	} else if identifierType == "id" {
			queryInput = &dynamodb.QueryInput{
					TableName: aws.String(USER_TABLE),
					KeyConditions: map[string]*dynamodb.Condition{
							"id": {
									ComparisonOperator: aws.String("EQ"),
									AttributeValueList: []*dynamodb.AttributeValue{
											{
													S: aws.String(identifier),
											},
									},
							},
					},
			}
	} else {
			return false, fmt.Errorf("invalid identifier type: %s", identifierType)
	}

	result, err := u.databaseStore.Query(queryInput)
	if err != nil {
			return false, err
	}

	return *result.Count > 0, nil
}


func (u DynamoDBClient) InsertUser(user types.ClientUser) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	if item["email"] == nil || item["password"] == nil {
		return fmt.Errorf("email and password are required fields")
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE),
		Item:      item,
	}

	_, err = u.databaseStore.PutItem(putItemInput)
	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) GetUser(email string) (types.ClientUser, error) {
	var user types.ClientUser

	result, err := u.databaseStore.Query(&dynamodb.QueryInput{
			TableName: aws.String(USER_TABLE),
			IndexName: aws.String("EmailIndex"), 
			KeyConditionExpression: aws.String("email = :email"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":email": {
							S: aws.String(email),
					},
			},
	})

	if err != nil {
			return user, err
	}

	if *result.Count == 0 {
			return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
			return user, err
	}

	return user, nil
}


func (u DynamoDBClient) GetAdminUser(email string) (types.User, error) {
	var user types.User

	result, err := u.databaseStore.Query(&dynamodb.QueryInput{
			TableName: aws.String(ADMIN_TABLE),
			IndexName: aws.String("EmailIndex"), 
			KeyConditionExpression: aws.String("email = :email"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":email": {
							S: aws.String(email),
					},
			},
	})

	if err != nil {
			return user, err
	}

	if *result.Count == 0 {
			return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
			return user, err
	}

	return user, nil
}

func (u DynamoDBClient) GetCosmetologistUser(email string) (types.CosmetologistUser, error) {
	var user types.CosmetologistUser

	result, err := u.databaseStore.Query(&dynamodb.QueryInput{
			TableName: aws.String(COSMETOLOGIST_TABLE),
			IndexName: aws.String("EmailIndex"),
			KeyConditionExpression: aws.String("email = :email"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":email": {
							S: aws.String(email),
					},
			},
	})

	if err != nil {
			return user, err
	}

	if *result.Count == 0 {
			return user, fmt.Errorf("cosmetologist user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &user)
	if err != nil {
			return user, err
	}

	return user, nil
}

func (u DynamoDBClient) InsertCosmetologistUser(user types.CosmetologistUser) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	if item["email"] == nil || item["password"] == nil {
		return fmt.Errorf("email and password are required fields")
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(COSMETOLOGIST_TABLE),
		Item:      item,
	}

	_, err = u.databaseStore.PutItem(putItemInput)
	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) DoesCosmetologistExist(identifier, identifierType string) (bool, error) {
	var queryInput *dynamodb.QueryInput

	if identifierType == "email" {
			// Query using the EmailIndex if the identifier is an email
			queryInput = &dynamodb.QueryInput{
					TableName: aws.String(COSMETOLOGIST_TABLE),
					IndexName: aws.String("EmailIndex"),
					KeyConditions: map[string]*dynamodb.Condition{
							"email": {
									ComparisonOperator: aws.String("EQ"),
									AttributeValueList: []*dynamodb.AttributeValue{
											{
													S: aws.String(identifier),
											},
									},
							},
					},
			}
	} else if identifierType == "id" {
			// Query using the primary key if the identifier is an ID
			queryInput = &dynamodb.QueryInput{
					TableName: aws.String(COSMETOLOGIST_TABLE),
					KeyConditions: map[string]*dynamodb.Condition{
							"id": {
									ComparisonOperator: aws.String("EQ"),
									AttributeValueList: []*dynamodb.AttributeValue{
											{
													S: aws.String(identifier),
											},
									},
							},
					},
			}
	} else {
			return false, fmt.Errorf("invalid identifier type: %s", identifierType)
	}

	// Execute the query
	result, err := u.databaseStore.Query(queryInput)
	if err != nil {
			return false, err
	}

	// Return true if any records are found
	return *result.Count > 0, nil
}
