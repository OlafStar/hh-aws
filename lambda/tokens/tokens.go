package tokens

import "github.com/google/uuid"

func CreateResetPassToken(email string) string {
	firstUUID, err := uuid.NewUUID()

	if err != nil {
		return ""
	}
	
	secondUUID, err := uuid.NewUUID()

	if err != nil {
		return ""
	}

	token := firstUUID.String() + secondUUID.String()

	return token
}