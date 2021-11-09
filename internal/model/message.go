package model

type Message struct {
	User      string `json:"user"`
	Content   string `json:"content"`
	Role      string `json:"role"`
	Timestamp string `json:"timestamp"`
}
