package chat

type Message struct {
	User      string `json:"user"`
	Content   string `json:"content"`
	Role      string `json:"role"`
	Timestamp string `json:"timestamp"`
}

func SystemMessage(content string) Message {
	return Message{
		User:    "",
		Content: content,
		Role:    "system",
	}
}
