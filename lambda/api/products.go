package api

import (
	"encoding/json"
	"lambda-func/types"
	"lambda-func/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)


func (api APIHandler) CreateProduct(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var product types.Product

	err := json.Unmarshal([]byte(request.Body), &product)
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid request: Unable to parse JSON"), err
	}

	if product.Name == "" {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid request: Product name is required"), nil
	}

	err = api.dbStore.InsertProduct(product)
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), err
	}

	return utils.CreateAPIGatewayProxyResponse(http.StatusOK, "Product created successfully"), nil
}