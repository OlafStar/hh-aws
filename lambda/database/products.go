package database

import (
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

func (u DynamoDBClient) InsertProduct(product types.Product) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
			return err
	}
	//asseble item
	item := &dynamodb.PutItemInput{
		TableName: aws.String(PRODUCTS_TABLE),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(newUUID.String()),
			},
			"name": {
				S: aws.String(product.Name),
			},
		},
	}

	if product.Image != "" {
		item.Item["image"] = &dynamodb.AttributeValue{
			S: aws.String(product.Image),
		}
	}
	
	_, err = u.databaseStore.PutItem(item)

	if err != nil {
		return err
	}

	return nil
}