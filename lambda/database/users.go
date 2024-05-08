package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)


func (u DynamoDBClient) DoesUserExists(email string) (bool, error) {
	result, err := u.databaseStore.Query(&dynamodb.QueryInput{
			TableName: aws.String(USER_TABLE),
			IndexName: aws.String("EmailIndex"),
			KeyConditions: map[string]*dynamodb.Condition{
					"email": {
							ComparisonOperator: aws.String("EQ"),
							AttributeValueList: []*dynamodb.AttributeValue{
									{
											S: aws.String(email),
									},
							},
					},
			},
	})

	if err != nil {
			return false, err
	}

	// If Count is more than 0, the user exists
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

func (u DynamoDBClient) DoesCosmetologistUserExists(email string) (bool, error) {
	result, err := u.databaseStore.Query(&dynamodb.QueryInput{
		TableName: aws.String(COSMETOLOGIST_TABLE),
		IndexName: aws.String("EmailIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
				"email": {
						ComparisonOperator: aws.String("EQ"),
						AttributeValueList: []*dynamodb.AttributeValue{
								{
										S: aws.String(email),
								},
						},
				},
		},
	})

	if err != nil {
			return false, err
	}

	return *result.Count > 0, nil
}