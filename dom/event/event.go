package event

type (
	DomainEvent interface{}
	Stream      chan DomainEvent
)
