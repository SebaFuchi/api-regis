package user

type User struct {
	Id             int    `json:"id"`
	Token          string `json:"token"`
	HashedPassword string
	Sex            string `json:"sex"`
	Name           string `json:"name"`
	NameTag        string `json:"nameTag"`
	RegisDate      string `json:"regisDate"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	Image          string `json:"image"`
}
