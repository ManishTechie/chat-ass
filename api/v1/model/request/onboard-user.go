package request

type User struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
}

func NewUser() *User {
	return new(User)
}
