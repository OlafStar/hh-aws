package database

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ResetToken struct {
	Email     string `json:"email"`
	Token     string `json:"token"`
	Expires   int64  `json:"expires"`
}

func (u DynamoDBClient) GetResetPassTokenByEmail(email string) (*ResetToken, error) {
	cond := expression.Key("Email").Equal(expression.Value(email))
	expr, err := expression.NewBuilder().WithKeyCondition(cond).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression: %v", err)
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(RESET_TOKENS_TABLE),
		IndexName:                 aws.String("EmailIndex"), 
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := u.databaseStore.Query(queryInput)
	if err != nil {
		return nil, fmt.Errorf("error querying DynamoDB: %v", err)
	}

	if len(result.Items) == 0 {
		return nil, nil 
	}

	var resetToken ResetToken
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &resetToken)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling item: %v", err)
	}

	return &resetToken, nil
}

func (u DynamoDBClient) GetResetPassByToken(token string) (*ResetToken, error) {
	cond := expression.Key("token").Equal(expression.Value(token))
	expr, err := expression.NewBuilder().WithKeyCondition(cond).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression: %v", err)
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(RESET_TOKENS_TABLE),
		IndexName:                 aws.String("TokenIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := u.databaseStore.Query(queryInput)
	if err != nil {
		return nil, fmt.Errorf("error querying DynamoDB: %v", err)
	}

	if len(result.Items) == 0 {
		return nil, nil 
	}

	var resetToken ResetToken
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &resetToken)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling item: %v", err)
	}

	return &resetToken, nil
}

func (u DynamoDBClient) CreateResetPassToken(email, token string) (*ResetToken,error) {
	expires := time.Now().Add(1 * time.Hour).Unix()

	resetToken := &ResetToken{
		Email:   email,
		Token:   token,
		Expires: expires,
	}

	item, err := dynamodbattribute.MarshalMap(resetToken)
	if err != nil {
		return nil, fmt.Errorf("error marshaling reset token: %v", err)
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(RESET_TOKENS_TABLE),
		Item:      item,
	}

	_, err = u.databaseStore.PutItem(putItemInput)
	if err != nil {
		return nil, fmt.Errorf("error putting item into DynamoDB: %v", err)
	}

	return resetToken, nil
}

func (u DynamoDBClient) ExpireResetToken(email string, token string) error {
	newExpire := time.Now().Unix()
	input := &dynamodb.UpdateItemInput{
    TableName: aws.String(RESET_TOKENS_TABLE),
    Key: map[string]*dynamodb.AttributeValue{
        "email": {
            S: aws.String(email),
        },
    },
    UpdateExpression: aws.String("set expires = :e, #t = :t"),
    ExpressionAttributeNames: map[string]*string{
        "#t": aws.String("token"),
    },
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
        ":e": {
            N: aws.String(fmt.Sprintf("%d", newExpire)),
        },
        ":t": {
            S: aws.String(token), 
        },
    },
    ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := u.databaseStore.UpdateItem(input)
	if err != nil {
			return fmt.Errorf("failed to expire reset token: %v", err)
	}

	return nil
}
