package types

type UserRegisterer interface {
	GetEmail() string
	GetPassword() string
}

func (u RegisterUser) GetEmail() string {
    return u.Email
}

func (u RegisterUser) GetPassword() string {
    return u.Password
}

func (u ClientUserRegisterBody) GetEmail() string {
    return u.Email
}

func (u ClientUserRegisterBody) GetPassword() string {
    return u.Password
}

func (u CosmetologistUserRegisterBody) GetEmail() string {
    return u.Email
}

func (u CosmetologistUserRegisterBody) GetPassword() string {
    return u.Password
}

type Product struct {
	Name string `json:"name"`
	Image string `json:"image"`
}
