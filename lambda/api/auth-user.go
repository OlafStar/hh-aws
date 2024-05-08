package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/jwt"
	"lambda-func/tokens"
	"lambda-func/types"
	"lambda-func/utils"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (api APIHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body types.ClientUserRegisterBody

	err := json.Unmarshal([]byte(request.Body), &body)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid request"), err
	}

	if err := body.Validate(); err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid request - %s", err.Error())), err
	}

	userExists, err := api.dbStore.DoesUserExist(body.Email, "email")

	if err != nil {
		
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), fmt.Errorf("there an error checking id user exists %w", err)
	}

	if userExists {
		
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusConflict,"User already exists"), fmt.Errorf("this user exists")
	}

	user, err := jwt.NewUser(body)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), fmt.Errorf("could not create new user - %w", err)
	}

	newUUID, err := uuid.NewUUID()

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), err
	}

	err = api.dbStore.InsertUser(types.ClientUser{
		Id: newUUID.String(),
		Firstname: body.Firstname,
		Surname: body.Surname,
		Email: body.Email,
		PasswordHash: user.PasswordHash,
		Image: body.Image,
	})

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

	accessToken, err := jwt.CreateToken(types.User{Email: user.Email, PasswordHash: user.PasswordHash}, jwt.RoleUser)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), err
	}

	successMsg := fmt.Sprintf(`{"accessToken": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body: successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

func (api APIHandler) ResetPasswordUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type Body struct {
		Email string `json:"email"`
	}

	var body Body

	err := json.Unmarshal([]byte(request.Body), &body)

	if err != nil {
		return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest,"Invalid request"), err
	}

	resetToken, err := api.dbStore.GetResetPassTokenByEmail(body.Email)

	if resetToken == nil {
		token := tokens.CreateResetPassToken(body.Email)

		if token == "" {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Error with creating token"), err
		}

		//TODO: Implement email sending in the future
		_, err = api.dbStore.CreateResetPassToken(body.Email, token)

		if err != nil {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError,"Internal server error"), err
		}
	}

	return events.APIGatewayProxyResponse{
		Body: `{"message": "Success"}`,
		StatusCode: http.StatusOK,
	}, nil
}

func (api APIHandler) ResetPasswordUserValidate(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token, tokenPresent := request.QueryStringParameters["token"]
	if !tokenPresent || token == "" {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Token is required"), nil
	}

	var body struct {
			Password string `json:"password"`
	}
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Invalid request body"), err
	}
	if body.Password == "" {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Password is required"), nil
	}

	resetToken, err := api.dbStore.GetResetPassByToken(token)
	if err != nil {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Internal server error"), err
	}
	if resetToken == nil {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusNotFound, "No user associated with this token"), nil
	}

	currentTime := time.Now().Unix()
	if currentTime > resetToken.Expires {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusBadRequest, "Token expired"), nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Failed to hash password"), err
	}

	err = api.dbStore.UpdateUserPassword(resetToken.Email, string(hashedPassword))
	if err != nil {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Failed to update password"), err
	}

	err = api.dbStore.ExpireResetToken(resetToken.Email, token)
	if err != nil {
			return utils.CreateAPIGatewayProxyErrorResponse(http.StatusInternalServerError, "Failed to expire reset token"), err
	}

	return events.APIGatewayProxyResponse{
			Body:       `{"message": "Password reset successfully"}`,
			StatusCode: http.StatusOK,
	}, nil
}