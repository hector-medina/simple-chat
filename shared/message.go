package shared

// Message
// Logical design:
// nickName: Text
// message: Text
type Message struct {
	ID      int    `json:"id"`
	Channel string `json:"channel"`
	Author  string `json:"author"`
	Message string `json:"message"`
}
