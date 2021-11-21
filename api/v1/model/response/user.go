package response

type Users struct {
	UserID   string `json:"user_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Gender   string `json:"gender,omitempty"`
	StreamID string `json:"stream_id,omitempty"`
}
