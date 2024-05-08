package main

import (
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func ProtectedRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body: "This is secret path",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.APIHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.APIHandler.LoginUser(request)
		case "/reset-password/create":
			return myApp.APIHandler.ResetPasswordUser(request)
		case "/reset-password/validate":
			return myApp.APIHandler.ResetPasswordUserValidate(request)
		case "/admin/clients":
			return middleware.ValidateAdminJWT(myApp.APIHandler.GetClients)(request)
		case "/admin/cosmetologists":
			return middleware.ValidateAdminJWT(myApp.APIHandler.GetCosmetologists)(request)
		case "/admin/cosmetologist/register":
			return middleware.ValidateAdminJWT(myApp.APIHandler.RegisterCosmetologistHandler)(request)
		case "/admin/cosmetologist/assign":
			return middleware.ValidateAdminJWT(myApp.APIHandler.AssignClientToCosmetologist)(request)
		case "/admin/products/create":
			return middleware.ValidateAdminJWT(myApp.APIHandler.CreateProduct)(request)
		case "/admin/login":
			return myApp.APIHandler.LoginAdminUser(request)
		case "/cosmetologist/login":
			return myApp.APIHandler.LoginCosmetologistUser(request)
		case "/protected-user":
			return middleware.ValidateUserJWT(ProtectedRequest)(request)
		case "/protected-cosmetologist":
			return middleware.ValidateCosmetologistJWT(ProtectedRequest)(request)
		case "/protected-admin":
			return middleware.ValidateAdminJWT(ProtectedRequest)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body: "Not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}