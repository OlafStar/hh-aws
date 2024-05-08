package types

type RegisterUser struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Email string `json:"email"`
	PasswordHash string `json:"password"`
}

type Product struct {
	Name string `json:"name"`
	Image string `json:"image"`
}