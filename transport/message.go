package transport

import "encoding/json"

type Message struct {
	Msg string `json:"message"`
}

func NewMessage(msg string) *Message {
	return &Message{Msg: msg}
}

func (m Message) String() string {
	j, _ := json.Marshal(m)
	return string(j)
}
