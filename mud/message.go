package mud

type MessageType int

const (
	SystemInfoMsg MessageType = iota
	LeaveMsg
)

type Message struct {
	Type  MessageType
	Value interface{}
}

func NewMessage(t MessageType) *Message {
	return &Message{
		Type:  t,
		Value: nil,
	}
}

func NewMessageWithValue(t MessageType, value interface{}) *Message {
	return &Message{
		Type:  t,
		Value: value,
	}
}
