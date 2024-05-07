package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/jwt"
	"lambda-func/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)


func (api APIHandler) LoginAdminUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid request"), err
	}

	user, err := api.dbStore.GetAdminUser(loginRequest.Email)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), err
	}

	if !jwt.ValidatePassowrd(user.PasswordHash, loginRequest.Password) {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid user credentials"), err
	}

	accessToken, err := jwt.CreateToken(user, jwt.RoleAdmin)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), err
	}

	successMsg := fmt.Sprintf(`{"accessToken": "%s"}`, accessToken)
	
	return events.APIGatewayProxyResponse{
		Body: successMsg,
		StatusCode: http.StatusOK,
	}, nil
}