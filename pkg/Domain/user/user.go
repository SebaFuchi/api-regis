package user

type User struct {
	Id             int    `json:"id"`
	Token          string `json:"token"`
	HashedPassword string `json:"hashedPassword,omitempty"`
	Name           string `json:"name"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
}
