package simula

import (
	"container/heap"
	"sync"
	"time"

	"github.com/digimortl/tycoon/msgbox"
)

type Interface interface {
	Run()
}

type Simulator struct {
	currentTime     time.Time
	activeProcesses sync.WaitGroup
	scheduledEvents PriorityEventQueue
	schedule        msgbox.MessageBox
	stop            chan bool
}

func NewSimulator() *Simulator {
	sim := Simulator{
		currentTime:     time.Time{},
		activeProcesses: sync.WaitGroup{},
		scheduledEvents: nil,
		schedule:        make(msgbox.MessageBox),
		stop:            make(chan bool),
	}
	go sim.run()
	return &sim
}

func (s *Simulator) NewEvent() *DelayedEvent {
	return s.NewEventAt(time.Time{})
}

func (s *Simulator) NewEventAt(occurrenceTime time.Time) *DelayedEvent {
	return &DelayedEvent{
		createdAt: time.Now(),
		occurAt:   occurrenceTime,
		sim:       s,
		block:     make(chan bool),
	}
}

func (s *Simulator) activateProcess() {
	s.activeProcesses.Add(1)
}

func (s *Simulator) deactivateProcess() {
	s.activeProcesses.Done()
}

func (s *Simulator) waitTillThereAreNoActiveProcesses() {
	s.activeProcesses.Wait()
}

func (s *Simulator) hasNoEvents() bool {
	return s.scheduledEvents.Len() == 0
}

func (s *Simulator) pushAnEvent(anEvent *DelayedEvent) {
	heap.Push(&s.scheduledEvents, anEvent)
}
func (s *Simulator) popNextEvent() *DelayedEvent {
	elem := heap.Pop(&s.scheduledEvents)
	return elem.(*DelayedEvent)
}

func (s *Simulator) terminate() {
	close(s.schedule)
	close(s.stop)
}

func (s *Simulator) run() {
	defer s.terminate()
	for {
		select {
		case msg := <-s.schedule:
			anEvent := msg.Body.(*DelayedEvent)
			s.pushAnEvent(anEvent)
			msg.Ack()
		case <-s.stop:
			return
		}
	}
}

func (s *Simulator) WakeUpAfter(duration time.Duration) time.Time {
	anEvent := s.NewEventAt(s.currentTime.Add(duration))
	defer anEvent.Close()
	s.schedule.SendWithAck(anEvent)
	anEvent.Suspend()
	return s.currentTime
}

func (s *Simulator) Spawn(simulationObject Interface) {
	s.activateProcess()
	go func() {
		defer s.deactivateProcess()
		simulationObject.Run()
	}()
}

func (s *Simulator) Stop() {
	s.stop <- true
}

type StopCondition func() bool

func (s *Simulator) Proceed(hasStopConditionReached StopCondition) time.Duration {
	for {
		s.waitTillThereAreNoActiveProcesses()

		if s.hasNoEvents() || hasStopConditionReached() {
			break
		}

		anEvent := s.popNextEvent()
		s.currentTime = anEvent.occurAt
		anEvent.Resume()
	}
	return s.currentTime.Sub(time.Time{})
}
