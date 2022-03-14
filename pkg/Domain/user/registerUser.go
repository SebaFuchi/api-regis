package user

type RegisterUser struct {
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	NameTag  string `json:"nameTag"`
	Birthday string `json:"birthday"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
