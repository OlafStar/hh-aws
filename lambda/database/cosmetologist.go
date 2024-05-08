package database

import (
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (client DynamoDBClient) GetCosmetologists(page int64, limit int64) ([]types.CosmetologistUserSecure, int64, int64, error) {
	queryInput := &dynamodb.ScanInput{
		TableName: aws.String(COSMETOLOGIST_TABLE),
		Limit:     aws.Int64(limit),
	}

	var accumulatedItemCount int64 = 0
	var pageStartsHere bool = page == 1
	var cosmetologists []types.CosmetologistUserSecure
	for {
		result, err := client.databaseStore.Scan(queryInput)
		if err != nil {
			return nil, 0, 0, err
		}
		
		if !pageStartsHere {
			accumulatedItemCount += *result.Count
			if accumulatedItemCount >= (page - 1)*limit {
				pageStartsHere = true
				cosmetologists = nil 
			}
		}
		
		if pageStartsHere {
			if len(result.Items) > 0 {
				var pageClients []types.CosmetologistUserSecure
				err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &pageClients)
				if err != nil {
					return nil, 0, 0, err
				}
				cosmetologists = append(cosmetologists, pageClients...)
			}
		}

		if result.LastEvaluatedKey == nil || len(cosmetologists) >= int(limit) {
			break
		}
		queryInput.ExclusiveStartKey = result.LastEvaluatedKey
	}

	var nextPage int64 = 0
	if int64(len(cosmetologists)) == limit && accumulatedItemCount < (page-1)*limit+int64(len(cosmetologists)) {
		nextPage = page + 1
	}
	
	var prevPage int64 = 0
	if page > 1 {
		prevPage = page - 1
	}

	return cosmetologists, nextPage, prevPage, nil
}

func (client DynamoDBClient) CountCosmetologists() (int64, error) {
	queryInput := &dynamodb.ScanInput{
		TableName: aws.String(COSMETOLOGIST_TABLE),
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