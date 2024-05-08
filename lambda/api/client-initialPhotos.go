package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/jwt"
	"lambda-func/middleware"
	"lambda-func/types"
	"lambda-func/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (api APIHandler) InitialPhotosOfClient(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var photosPostBody types.InitialPhotosOfClientPostBody
	if err := json.Unmarshal([]byte(request.Body), &photosPostBody); err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid request body"), fmt.Errorf("error unmarshalling request body: %w", err)
	}

	if valid, err := types.ValidateInitialPhotos(photosPostBody.Images); !valid {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, err.Error()), fmt.Errorf("validation error: %w", err)
	}

	token := middleware.ExtractTokenFromHeaders(request.Headers)

	parsedToken, err := jwt.ParseToken(token, jwt.RoleUser)
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusUnauthorized, "Unauthenticated"), fmt.Errorf("error parsing token: %w", err)
	}

	email, ok := parsedToken["email"].(string)
	if !ok {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusUnauthorized, "Email is required in token"), fmt.Errorf("email not found in token")
	}

	clientID, err := api.dbStore.GetClientIDByEmail(email)
	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Failed to get client ID"), fmt.Errorf("error retrieving client ID: %w", err)
	}

	submited, err := api.dbStore.HasUserSubmittedPhotos(clientID)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Errow with checking is user submited photos"), fmt.Errorf("errow has user submitted photos: %w", err)
	}

	if submited {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "User already submited initial photos"), fmt.Errorf("user already submited initial photos: %w", err)
	}

	err = api.dbStore.InsertInitialPhotos(clientID, photosPostBody.Images)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Error with insering images"), fmt.Errorf("error with inserting images: %w", err)
	}
	
	return events.APIGatewayProxyResponse{
		Body:       `{"message":"Success"}`,
		StatusCode: http.StatusOK,
	}, nil
}
