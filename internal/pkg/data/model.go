package data

type Ad struct {
	ID    uint   `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Color string `json:"color"`
	Price int    `json:"price"`
}

//Geter`s
func (a *Ad) GetID() uint {
	return a.ID
}

func (a *Ad) SetID(id uint) {
	a.ID = id
}

func (a *Ad) GetBrand() string {
	return a.Brand
}

func (a *Ad) GetModel() string {
	return a.Model
}

func (a *Ad) GetColor() string {
	return a.Color
}

func (a *Ad) GetPrice() int {
	return a.Price
}

type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"-"`
}

func (a *Account) GetID() int {
	return a.ID
}

func (a *Account) GetUserName() string {
	return a.Username
}

func (a *Account) GetPassword() string {
	return a.Password
}

func (a *Account) GetToken() string {
	return a.Token
}

func (a *Account) SetToken(token string) {
	a.Token = token
}
