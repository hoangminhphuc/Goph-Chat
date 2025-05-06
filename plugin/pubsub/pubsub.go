package pubsub

import (
	"time"
	"fmt"
)
type PubSub interface {
	Publish()
	Subscribe()
}

type Topic string

type Message struct {
	id    string
	data  interface{}
	topic Topic
}

func NewMessage(topic Topic, data interface{}) *Message {
	now := time.Now().UTC()
	return &Message{
		id:    fmt.Sprintf("%d", now.UnixNano()),
		data:  data,
		topic: topic,
	}
}

func (m *Message) GetData() interface{} {
	return m.data
}
func (m *Message) GetTopic() Topic {
	return m.topic
}

func (m *Message) SetTopic(topic Topic) {
	m.topic = topic
}