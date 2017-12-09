package mandrill

type MessageRequest struct {
	Key     string   `json:"key"`
	Message *Message `json:"message"`
}

func NewMessageRequest(key string, message *Message) *MessageRequest {
	return &MessageRequest{
		Key:     key,
		Message: message,
	}
}
