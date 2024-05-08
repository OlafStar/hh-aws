package database

import (
	"lambda-func/types"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	USER_TABLE = "userTable"
	ADMIN_TABLE = "adminUserTable"
	COSMETOLOGIST_TABLE = "cosmetologistUserTable"
	PRODUCTS_TABLE = "productsTable"
)

type Store interface {
	DoesUserExist(identifier, identifierType string) (bool, error)
	GetUser(email string) (types.ClientUser, error)
	GetAdminUser(email string) (types.User, error)
	GetCosmetologistUser(email string) (types.CosmetologistUser, error)
	DoesCosmetologistExist(identifier, identifierType string) (bool, error)
	InsertUser(event types.ClientUser) error
	InsertCosmetologistUser(event types.CosmetologistUser) error
	InsertProduct(event types.Product) error
	GetClients(page int64, limit int64) ([]types.ClientUserResponse, int64, int64, error)
	CountClients() (int64, error)
	GetCosmetologists(page int64, limit int64) ([]types.CosmetologistUserSecure, int64, int64, error)
	CountCosmetologists() (int64, error)
	AssignCosmetologistToClient(clientId, newCosmetologistId string) error
}

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)


	return DynamoDBClient{
		databaseStore: db,
	}
}