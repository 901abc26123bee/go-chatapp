package connection

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
