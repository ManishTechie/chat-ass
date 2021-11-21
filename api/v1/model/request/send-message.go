package request

type UserMessages struct {
	Message string `json:"message"`
}

func NewUserMessages() *UserMessages {
	return new(UserMessages)
}
