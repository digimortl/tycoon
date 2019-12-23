package event

type (
	DomainEvent interface{}
	Stream      chan DomainEvent
)


func NullEventStream() Stream {
	es := make(Stream)
	go func() {
		for range es {}
	}()
	return es

}

func NewEventStream(handler func(DomainEvent)) Stream {
	es := make(Stream)
	go func() {
		defer close(es)
		for ev := range es {
			handler(ev)
		}
	}()
	return es
}
