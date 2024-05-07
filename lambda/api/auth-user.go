package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/jwt"
	"lambda-func/types"
	"lambda-func/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (api APIHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest,"Invalid request"), err
	}

	if registerUser.Email == "" || registerUser.Password == "" {
		
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest,"Invalid request - empty parameters"), fmt.Errorf("request has empty parameters")
	}

	userExists, err := api.dbStore.DoesUserExists(registerUser.Email)

	if err != nil {
		
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), fmt.Errorf("there an error checking id user exists %w", err)
	}

	if userExists {
		
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusConflict,"User already exists"), fmt.Errorf("this user exists")
	}

	user, err := jwt.NewUser(registerUser)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), fmt.Errorf("could not create new user - %w", err)
	}

	err = api.dbStore.InsertUser(user)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), fmt.Errorf("error registering user %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body: "Succesfuly registered user",
		StatusCode: http.StatusOK,
	}, nil
}

func (api APIHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest,"Invalid request"), err
	}

	user, err := api.dbStore.GetUser(loginRequest.Email)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), err
	}

	if !jwt.ValidatePassowrd(user.PasswordHash, loginRequest.Password) {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest,"Invalid user credentials"), nil
	}

	accessToken, err := jwt.CreateToken(user, jwt.RoleUser)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), err
	}

	successMsg := fmt.Sprintf(`{"accessToken": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body: successMsg,
		StatusCode: http.StatusOK,
	}, nil
}
