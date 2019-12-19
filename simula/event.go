package simula

import "time"

type DelayedEvent struct {
	createdAt time.Time
	occurAt   time.Time
	sim       *Simulator
	block     chan bool
}

func (e *DelayedEvent) Close() {
	close(e.block)
}

func (e *DelayedEvent) Suspend() {
	e.sim.deactivateProcess()
	<-e.block
}

func (e *DelayedEvent) Resume() {
	e.sim.activateProcess()
	e.block <- true
}

type PriorityEventQueue []*DelayedEvent

func (q *PriorityEventQueue) Less(i, j int) bool {
	return (*q)[i].occurAt.Before((*q)[j].occurAt)
}

func (q *PriorityEventQueue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

func (q *PriorityEventQueue) Len() int {
	return len(*q)
}

func (q *PriorityEventQueue) Pop() (v interface{}) {
	*q, v = (*q)[:q.Len()-1], (*q)[q.Len()-1]
	return
}

func (q *PriorityEventQueue) Push(v interface{}) {
	*q = append(*q, v.(*DelayedEvent))
}
