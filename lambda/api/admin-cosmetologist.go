package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (api APIHandler) GetCosmetologists(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	limit := int64(10)
	page := int64(1)

	if val, ok := request.QueryStringParameters["limit"]; ok {
		fmt.Sscanf(val, "%d", &limit)
	}
	if val, ok := request.QueryStringParameters["page"]; ok {
		fmt.Sscanf(val, "%d", &page)
	}

	totalCount, err := api.dbStore.CountCosmetologists()
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), nil
	}

	cosmetologists, nextPage, prevPage, err := api.dbStore.GetCosmetologists(page, limit)
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), nil
	}

	lastPage := totalCount / limit
	if totalCount%limit != 0 {
		lastPage++
	}

	responseBody, err := json.Marshal(map[string]interface{}{
		"page":      page,
		"lastPage":  lastPage,
		"next":      nextPage,
		"previous":  prevPage,
		"cosmetologists":   cosmetologists,
	})
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}