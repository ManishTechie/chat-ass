package response

import stream "github.com/GetStream/stream-chat-go/v2"

type Onboard struct {
	User     stream.User `json:"user"`
	Token    string      `json:"token"`
	JWTToken string      `json:"jwt_token"`
}
