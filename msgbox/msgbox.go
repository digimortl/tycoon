package msgbox

type Message struct {
	Body  interface{}
	reply chan interface{}
}

func newMessage(body interface{}) Message {
	return Message{body, make(chan interface{})}
}

func (m *Message) Ack() {
	m.Reply(Whatever())
}

func (m *Message) Reply(data interface{}) {
	m.reply <- data
}

type MessageBox chan Message

func (box *MessageBox) SendAndReceive(body interface{}) interface{} {
	msg := newMessage(body)
	defer close(msg.reply)
	*box <- msg
	return <-msg.reply
}

func (box *MessageBox) SendWithAck(body interface{}) {
	box.SendAndReceive(body)
}

func Whatever() interface{} {
	return struct{}{}
}
