package utils

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func CreateAPIGatewayProxyErrorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       fmt.Sprintf(`{"message": "%s"}`, message),
	}
}