package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	APIHandler api.APIHandler
}

func NewApp() App {
	//we initialize our db store
	db := database.NewDynamoDBClient()
	apiHandler := api.NewAPIHandler(db)

	return App{
		APIHandler: apiHandler,
	}
}