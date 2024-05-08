package api

import (
	"lambda-func/database"
)

type APIHandler struct {
	dbStore database.Store
}

func NewAPIHandler(dbStore database.Store) APIHandler {
	return APIHandler{
		dbStore: dbStore,
	}
}

