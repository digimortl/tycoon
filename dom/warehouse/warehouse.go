package warehouse

import (
	"github.com/digimortl/tycoon/msgbox"
	"github.com/digimortl/tycoon/simula"
)

type LocationCode = string

type Cargo struct {
	TrackNumber string
	Origin      LocationCode
	Destination LocationCode
}

type Warehouse struct {
	Location     LocationCode
	cargoes      []*Cargo
	sim          *simula.Simulator
	bring        msgbox.MessageBox
	pick         msgbox.MessageBox
	waitForCargo msgbox.MessageBox
	waiters      []*simula.DelayedEvent
	stop         chan bool
}

func Of(location LocationCode, sim *simula.Simulator) *Warehouse {
	w := Warehouse{
		Location:     location,
		cargoes:      nil,
		sim:          sim,
		bring:        make(msgbox.MessageBox),
		pick:         make(msgbox.MessageBox),
		waitForCargo: make(msgbox.MessageBox),
		waiters:      nil,
		stop:         make(chan bool),
	}
	go w.run()
	return &w
}

func (w *Warehouse) putCargo(aCargo *Cargo) {
	w.cargoes = append(w.cargoes, aCargo)
}

func (w *Warehouse) takeFirstCargo() *Cargo {
	var aCargo *Cargo = nil
	if len(w.cargoes) > 0 {
		aCargo, w.cargoes = w.cargoes[0], w.cargoes[1:]
	}
	return aCargo
}

func (w *Warehouse) addWaiter(waiter *simula.DelayedEvent) {
	w.waiters = append(w.waiters, waiter)
}

func (w *Warehouse) removeWaiter() *simula.DelayedEvent {
	var waiter *simula.DelayedEvent = nil
	if len(w.waiters) > 0 {
		waiter, w.waiters = w.waiters[0], w.waiters[1:]
	}
	return waiter
}

func (w *Warehouse) isEmpty() bool {
	return len(w.cargoes) == 0
}

func (w *Warehouse) terminate() {
	close(w.bring)
	close(w.pick)
	close(w.waitForCargo)
	close(w.stop)
}

func (w *Warehouse) run() {
	defer w.terminate()
	for {
		select {
		case msg := <-w.bring:
			aCargo := msg.Body.(*Cargo)
			w.putCargo(aCargo)
			if waiter := w.removeWaiter(); waiter != nil {
				waiter.Resume()
			}
			msg.Ack()
		case msg := <-w.pick:
			aCargo := w.takeFirstCargo()
			msg.Reply(aCargo)
		case msg := <-w.waitForCargo:
			if w.isEmpty() {
				waiter := w.sim.NewEvent()
				w.addWaiter(waiter)
				msg.Reply(waiter)
			} else {
				msg.Reply(nil)
			}
		case <-w.stop:
			return
		}
	}
}

func (w *Warehouse) Bring(cargoes ...*Cargo) {
	for _, aCargo := range cargoes {
		msgbox.SendWithAck(w.bring, aCargo)
	}
}

func (w *Warehouse) PickCargo() *Cargo {
	answer := msgbox.SendAndReceive(w.pick, msgbox.Whatever())
	return answer.(*Cargo)
}

func (w *Warehouse) WaitForCargo() {
	if answer := msgbox.SendAndReceive(w.waitForCargo, msgbox.Whatever()); answer != nil {
		waiter := answer.(*simula.DelayedEvent)
		defer waiter.Close()
		waiter.Suspend()
	}
}

func (w *Warehouse) Fullness() int {
	return len(w.cargoes)
}

func (w *Warehouse) Stop() {
	w.stop <- true
}
