package response

import (
	"errors"
	"fmt"
	"strings"
)

type Message struct {
	Type    string
	Content string
}

type MessageHandler struct {
	Callback func(msg *Message) error
}

func (*MessageHandler) Supports(data string) bool {
	return data[0] == '[' && data[len(data)-1] == ']'
}

func (h *MessageHandler) Handle(data string) error {
	msg := data[1 : len(data)-1]
	parts := strings.SplitN(msg, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("message: invalid message: %s", data)
	}

	if h.Callback == nil {
		return errors.New("message: no callback defined")
	}

	return h.Callback(&Message{
		Type:    parts[0],
		Content: parts[1],
	})
}
