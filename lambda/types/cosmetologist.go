package types

//API /cosmetologist/register - Body
type CosmetologistUserRegisterBody struct {
	Firstname string `json:"firstname"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Password string `json:"password"`
	Image string `json:"image"`
}

//Cosmetologist in database
type CosmetologistUser struct {
	Id string `json:"id"`
	Firstname string `json:"firstname"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	PasswordHash string `json:"password"`
	Image string `json:"image"`
}