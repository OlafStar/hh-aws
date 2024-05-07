package api

import (
	"lambda-func/database"
)

type APIHandler struct {
	dbStore database.UserStore
}

func NewAPIHandler(dbStore database.UserStore) APIHandler {
	return APIHandler{
		dbStore: dbStore,
	}
}

