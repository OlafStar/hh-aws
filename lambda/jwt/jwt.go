package jwt

import (
	"fmt"
	"lambda-func/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_SECRET = "secret-user"
	COSMETOLOGIST_SECRET = "secret-cosmetologist"
	ADMIN_SECRET = "secret-admin"
)

type Role int

const (
	RoleUser Role = iota
	RoleCosmetologist
	RoleAdmin
)

var roleSecrets = map[Role]string{
	RoleUser:          USER_SECRET,
	RoleCosmetologist: COSMETOLOGIST_SECRET,
	RoleAdmin:         ADMIN_SECRET,
}

func NewUser(user types.UserRegisterer) (types.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.GetPassword()), 10)
	if err != nil {
			return types.User{}, err
	}

	return types.User{
			Email: user.GetEmail(),
			PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatePassowrd(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func CreateToken(user types.User, role Role) (string, error) {
	secret, ok := roleSecrets[role]
	if !ok {
		return "", fmt.Errorf("invalid role specified")
	}

	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, role Role) (jwt.MapClaims, error) {
	secret, ok := roleSecrets[role]
	if !ok {
		return nil, fmt.Errorf("invalid role specified")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unauthorized: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid - unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims of unauthorized type")
	}

	return claims, nil
}