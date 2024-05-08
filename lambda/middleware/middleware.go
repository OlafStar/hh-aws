package middleware

import (
	"lambda-func/jwt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func ExtractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return ""
	}
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}
	return splitToken[1]
}

func ValidateJWTMiddleware(role jwt.Role) func(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			tokenString := ExtractTokenFromHeaders(request.Headers)
			if tokenString == "" {
				return events.APIGatewayProxyResponse{
					Body:       "Missing auth token",
					StatusCode: http.StatusUnauthorized,
				}, nil
			}

			claims, err := jwt.ParseToken(tokenString, role)
			if err != nil {
				return events.APIGatewayProxyResponse{
					Body:       "Unauthorized",
					StatusCode: http.StatusUnauthorized,
				}, err
			}

			expires := int64(claims["exp"].(float64))
			if time.Now().Unix() > expires {
				return events.APIGatewayProxyResponse{
					Body:       "Token expired",
					StatusCode: http.StatusUnauthorized,
				}, nil
			}

			return next(request)
		}
	}
}

// Middleware instances for each role
var ValidateUserJWT = ValidateJWTMiddleware(jwt.RoleUser)
var ValidateCosmetologistJWT = ValidateJWTMiddleware(jwt.RoleCosmetologist)
var ValidateAdminJWT = ValidateJWTMiddleware(jwt.RoleAdmin)
