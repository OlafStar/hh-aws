package database

import (
	"fmt"
	"lambda-func/types"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

func (client DynamoDBClient) InsertInitialPhotos(clientId string, photos []types.InitialPhotosStruct) error {
	if len(photos) != 6 {
		return fmt.Errorf("exactly 6 photos are required, received %d", len(photos))
	}

	var imageDetails []types.InitialPhotosStruct
	for _, photo := range photos {
		imageDetail := types.InitialPhotosStruct{
			Image: photo.Image,
			Type:  photo.Type,
		}
		imageDetails = append(imageDetails, imageDetail)
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("failed to create uuid: %w", err)
	}

	initialPhotosID := newUUID.String()

	photoRecord := types.InitialPhotoRecord{
		ClientID: clientId,
		Images:   imageDetails,
		ID:       initialPhotosID,
	}

	item, err := dynamodbattribute.MarshalMap(photoRecord)
	if err != nil {
		return fmt.Errorf("error marshalling photo record: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(INITIAL_PHOTOS_TABLE),
		Item:      item,
	}

	_, err = client.databaseStore.PutItem(input)
	if err != nil {
		log.Printf("Error putting photo item in DynamoDB: %v", err)
		return fmt.Errorf("error putting photo item in DynamoDB: %w", err)
	}

	err = client.updateClientUserWithInitialPhotosID(clientId, initialPhotosID)
	if err != nil {
			return fmt.Errorf("error updating user with initial photos ID: %w", err)
	}

	return nil
}

func (client DynamoDBClient) updateClientUserWithInitialPhotosID(clientId, initialPhotosID string) error {
	key, err := dynamodbattribute.MarshalMap(map[string]string{"id": clientId})
	if err != nil {
		return fmt.Errorf("error marshalling user ID: %w", err)
	}

	updateExpr, err := expression.NewBuilder().
		WithUpdate(expression.Set(expression.Name("initialPhotosId"), expression.Value(initialPhotosID))).
		Build()
	if err != nil {
		return fmt.Errorf("error building expression for updating user: %w", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(USER_TABLE),
		Key:                       key,
		UpdateExpression:          aws.String(*updateExpr.Update()),
		ExpressionAttributeNames:  updateExpr.Names(),
		ExpressionAttributeValues: updateExpr.Values(),
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	_, err = client.databaseStore.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("error updating user in DynamoDB: %w", err)
	}

	return nil
}

func (client DynamoDBClient) HasUserSubmittedPhotos(clientId string) (bool, error) {
	keyCond := expression.Key("clientId").Equal(expression.Value(clientId))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return false, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(INITIAL_PHOTOS_TABLE),
		IndexName:                 aws.String("ClientIdIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := client.databaseStore.Query(queryInput)
	if err != nil {
		return false, err
	}

	return *result.Count > 0, nil
}