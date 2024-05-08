package types

//API /client/register - Body
type ClientUserRegisterBody struct {
	Firstname string `json:"firstname"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Password string `json:"password"`
	Image string `json:"image"`
}

//Cosmetologist in database
type ClientUser struct {
	Id string `json:"id"`
	Firstname string `json:"firstname"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	PasswordHash string `json:"password"`
	Image string `json:"image"`
}