package types

import (
	"fmt"
	"regexp"
)

func (c *CosmetologistUserRegisterBody) Validate() error {
	if c.Firstname == "" || c.Surname == "" || c.Email == "" || c.Password == "" {
		return fmt.Errorf("missing required fields")
	}

	// Validate the email
	if !isValidEmail(c.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Validate the password (you can adjust the rules here)
	// if len(c.Password) < 8 { Example rule: at least 8 characters
	// 	return fmt.Errorf("password must be at least 8 characters long")
	// }

	return nil
}

func (c *ClientUserRegisterBody) Validate() error {
	if c.Firstname == "" || c.Surname == "" || c.Email == "" || c.Password == "" {
		return fmt.Errorf("missing required fields")
	}

	// Validate the email
	if !isValidEmail(c.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func ValidateInitialPhotos(photos []InitialPhotosStruct) (bool, error) {
	requiredTypes := map[PhotoType]bool{
		FullFace:  false,
		Forehead:  false,
		LeftSide:  false,
		RightSide: false,
		Nose:      false,
		Chin:      false,
	}
	for _, photo := range photos {
		if _, ok := requiredTypes[photo.Type]; ok {
			if !requiredTypes[photo.Type] {
				requiredTypes[photo.Type] = true
			} else {
				return false, fmt.Errorf("duplicate photo type: %s", photo.Type)
			}
		} else {
			return false, fmt.Errorf("invalid photo type: %s", photo.Type)
		}
	}

	for t, found := range requiredTypes {
		if !found {
			return false, fmt.Errorf("missing photo type: %s", t)
		}
	}
	return true, nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}