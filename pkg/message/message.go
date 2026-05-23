package message

import (
	"fmt"
	"time"
)

type Message struct {
	Time string
	From string
	Text string
}

func New(from, text string) Message {
	return Message{
		Time: time.Now().Format("2006-01-02 15:04:05"),
		From: from,
		Text: text,
	}
}

func (m Message) Format() string {
	if m.From == "System" {
		return fmt.Sprintf("[%s][System]: %s\n", m.Time, m.Text)
	}

	return fmt.Sprintf("[%s][%s]: %s\n", m.Time, m.From, m.Text)
}