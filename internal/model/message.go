package model

type Message struct {
	Type string `json:"type"`
	User string `json:"user,omitempty"`
	Chat Chat   `json:"chat,omitempty"`
}
