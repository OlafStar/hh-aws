package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (api APIHandler) GetClients(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	limit := int64(10)
	page := int64(1)

	if val, ok := request.QueryStringParameters["limit"]; ok {
		fmt.Sscanf(val, "%d", &limit)
	}
	if val, ok := request.QueryStringParameters["page"]; ok {
		fmt.Sscanf(val, "%d", &page)
	}

	totalCount, err := api.dbStore.CountClients()
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), nil
	}

	clients, nextPage, prevPage, err := api.dbStore.GetClients(page, limit)
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
		"clients":   clients,
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

type AssignClientToCosmetologistBody struct {
	ClientID string `json:"clientId"`
	CosmetologistID string `json:"cosmetologistId"`
}

func (api APIHandler) AssignClientToCosmetologist(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body AssignClientToCosmetologistBody

	err := json.Unmarshal([]byte(request.Body), &body)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid request"), err
	}

	err = api.dbStore.AssignCosmetologistToClient(body.ClientID, body.CosmetologistID)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Problem with assigning"), err
	}

	return utils.CreateAPIGatewayProxyResponse(http.StatusOK, "Succes"), nil
}