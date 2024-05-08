package database

import (
	"fmt"
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (client DynamoDBClient) GetClients(page int64, limit int64) ([]types.ClientUserResponse, int64, int64, error) {
	queryInput := &dynamodb.ScanInput{
		TableName: aws.String(USER_TABLE),
		Limit:     aws.Int64(limit),
	}

	var accumulatedItemCount int64 = 0
	var pageStartsHere bool = page == 1
	var clients []types.ClientUserResponse
	for {
		result, err := client.databaseStore.Scan(queryInput)
		if err != nil {
			return nil, 0, 0, err
		}
		
		if !pageStartsHere {
			accumulatedItemCount += *result.Count
			if accumulatedItemCount >= (page - 1)*limit {
				pageStartsHere = true
				clients = nil 
			}
		}
		
		if pageStartsHere {
			if len(result.Items) > 0 {
				var pageClients []types.ClientUserResponse
				err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &pageClients)
				if err != nil {
					return nil, 0, 0, err
				}
				for i := range pageClients {
					if pageClients[i].Image != nil && *pageClients[i].Image == "" {
						pageClients[i].Image = nil
					}
				}
				clients = append(clients, pageClients...)
			}
		}

		if result.LastEvaluatedKey == nil || len(clients) >= int(limit) {
			break
		}
		queryInput.ExclusiveStartKey = result.LastEvaluatedKey
	}

	var nextPage int64 = 0
	if int64(len(clients)) == limit && accumulatedItemCount < (page-1)*limit+int64(len(clients)) {
		nextPage = page + 1
	}
	
	var prevPage int64 = 0
	if page > 1 {
		prevPage = page - 1
	}

	return clients, nextPage, prevPage, nil
}

func (client DynamoDBClient) CountClients() (int64, error) {
	queryInput := &dynamodb.ScanInput{
		TableName: aws.String(USER_TABLE),
		Select:    aws.String("COUNT"),
	}

	var totalCount int64
	for {
		result, err := client.databaseStore.Scan(queryInput)
		if err != nil {
			return 0, err
		}
		totalCount += *result.Count

		if result.LastEvaluatedKey == nil {
			break
		}
		queryInput.ExclusiveStartKey = result.LastEvaluatedKey
	}

	return totalCount, nil
}

func (u DynamoDBClient) AssignCosmetologistToClient(clientId, newCosmetologistId string) error {
	clientExists, err := u.DoesUserExist(clientId, "id")
	if err != nil {
		return fmt.Errorf("error checking client existence: %w", err)
	}
	if !clientExists {
		return fmt.Errorf("client with ID %s does not exist", clientId)
	}

	cosmetologistExists, err := u.DoesCosmetologistExist(newCosmetologistId, "id")
	if err != nil {
		return fmt.Errorf("error checking new cosmetologist existence: %w", err)
	}
	if !cosmetologistExists {
		return fmt.Errorf("cosmetologist with ID %s does not exist", newCosmetologistId)
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(clientId)},
		},
		ProjectionExpression: aws.String("cosmetologistID"),
	}

	result, err := u.databaseStore.GetItem(getItemInput)
	if err != nil {
		return err
	}

	currentCosmetologistId := ""
	if result.Item != nil && result.Item["cosmetologistID"] != nil {
		currentCosmetologistId = *result.Item["cosmetologistID"].S
	}

	transactWriteItems := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Update: &dynamodb.Update{
					TableName: aws.String(USER_TABLE),
					Key: map[string]*dynamodb.AttributeValue{
						"id": {S: aws.String(clientId)},
					},
					UpdateExpression: aws.String("set cosmetologistID = :newcid"),
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":newcid": {S: aws.String(newCosmetologistId)},
					},
				},
			},
			{
				Update: &dynamodb.Update{
					TableName: aws.String(COSMETOLOGIST_TABLE),
					Key: map[string]*dynamodb.AttributeValue{
						"id": {S: aws.String(newCosmetologistId)},
					},
					UpdateExpression: aws.String("ADD clients :newclient"),
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":newclient": {SS: []*string{aws.String(clientId)}},
					},
				},
			},
		},
	}

	if currentCosmetologistId != "" && currentCosmetologistId != newCosmetologistId {
		transactWriteItems.TransactItems = append(transactWriteItems.TransactItems, &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				TableName: aws.String(COSMETOLOGIST_TABLE),
				Key: map[string]*dynamodb.AttributeValue{
					"id": {S: aws.String(currentCosmetologistId)},
				},
				UpdateExpression: aws.String("DELETE clients :oldclient"),
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":oldclient": {SS: []*string{aws.String(clientId)}},
				},
			},
		})
	}

	_, err = u.databaseStore.TransactWriteItems(transactWriteItems)
	return err
}
