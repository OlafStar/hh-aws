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
)

type UserStore interface {
	DoesUserExists(email string) (bool, error)
	InsertUser(event types.User) error
	GetUser(email string) (types.User, error)
	GetAdminUser(email string) (types.User, error)
	GetCosmetologistUser(email string) (types.User, error)
	DoesCosmetologistUserExists(email string) (bool, error)
	InsertCosmetologistUser(event types.User) error
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