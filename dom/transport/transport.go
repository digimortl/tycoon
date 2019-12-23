package transport

import (
	"time"

	"github.com/digimortl/tycoon/dom/event"
	m "github.com/digimortl/tycoon/dom/transmap"
	w "github.com/digimortl/tycoon/dom/warehouse"
	"github.com/digimortl/tycoon/simula"
)

type Transport struct {
	Name              string
	transportMap      *m.Map
	shipmentOption    m.ShipmentOption
	sim               *simula.Simulator
	home              *w.Warehouse
	arriveAtHomeAfter time.Duration
	cargoes           []*w.Cargo
	capacity          int
	timeToLoad        time.Duration
	timeToUnload      time.Duration
	log               chan event.DomainEvent
}

func (t *Transport) isFull() bool {
	return len(t.cargoes) == t.capacity
}

func (t *Transport) isEmpty() bool {
	return len(t.cargoes) == 0
}

func (t *Transport) load(aCargo *w.Cargo) {
	t.cargoes = append(t.cargoes, aCargo)
}

func (t *Transport) loadCargoesFrom(warehouse *w.Warehouse) {
	for !t.isFull() {
		if aCargo := warehouse.PickCargo(); aCargo != nil {
			t.load(aCargo)
		} else if t.isEmpty() {
			warehouse.WaitForCargo()
		} else {
			break
		}
	}
}

func (t *Transport) unload() *w.Cargo {
	var aCargo *w.Cargo
	aCargo, t.cargoes = t.cargoes[0], t.cargoes[1:]
	return aCargo
}

func (t *Transport) unloadCargoesTo(warehouse *w.Warehouse) {
	for !t.isEmpty() {
		aCargo := t.unload()
		warehouse.Bring(aCargo)
	}
}

func (t *Transport) findItineraryFrom(warehouse *w.Warehouse) m.Itinerary {
	aCargo := t.cargoes[0]
	itinerary, _ := t.transportMap.FindItinerary(warehouse.Location, aCargo.Destination)
	return itinerary.WithOption(t.shipmentOption)
}

func (t *Transport) Run() {
	t.log <- Arrived{
		OccurredAt:  t.holdFor(t.arriveAtHomeAfter),
		Transport:   t.Name,
		ShipmentOpt: t.shipmentOption,
		AtLocation:  t.home.Location,
	}

	t.loadCargoesFrom(t.home)
	t.log <- Loaded{
		OccurredAt:  t.holdFor(t.timeToLoad),
		Transport:   t.Name,
		ShipmentOpt: t.shipmentOption,
		Duration:    t.timeToLoad,
		Cargoes:     append(t.cargoes),
	}

	itinerary := t.findItineraryFrom(t.home)
	t.log <- Departed{
		OccurredAt:   t.holdFor(0),
		Transport:    t.Name,
		ShipmentOpt:  t.shipmentOption,
		FromLocation: t.home.Location,
		ToLocation:   itinerary.Destination().Location,
		Cargoes:      append(t.cargoes),
	}

	t.log <- Arrived{
		OccurredAt:  t.holdFor(itinerary.TotalTimeToTravel()),
		Transport:   t.Name,
		ShipmentOpt: t.shipmentOption,
		AtLocation:  itinerary.Destination().Location,
		Cargoes:     append(t.cargoes),
	}

	unloadTime := t.holdFor(t.timeToUnload)
	t.unloadCargoesTo(itinerary.Destination())
	t.log <- Unloaded{
		OccurredAt:  unloadTime,
		Transport:   t.Name,
		ShipmentOpt: t.shipmentOption,
		Duration:    t.timeToUnload,
	}

	t.log <- Departed{
		OccurredAt:   t.holdFor(0),
		Transport:    t.Name,
		ShipmentOpt:  t.shipmentOption,
		FromLocation: itinerary.Destination().Location,
		ToLocation:   t.home.Location,
	}

	t.arriveAtHomeAfter = itinerary.TotalTimeToTravel()
	t.Run()
}

func (t *Transport) StartJourneyFrom(warehouse *w.Warehouse) {
	t.home, t.arriveAtHomeAfter = warehouse, 0
	t.sim.Spawn(t)
}

func (t *Transport) holdFor(duration time.Duration) time.Time {
	return t.sim.WakeUpAfter(duration)
}

func Truck(name string, transportMap *m.Map, sim *simula.Simulator, log event.Stream) *Transport {
	return &Transport{
		Name:           name,
		transportMap:   transportMap,
		shipmentOption: m.Land,
		sim:            sim,
		capacity:       1,
		log:            log,
	}
}

func Vessel(name string, transportMap *m.Map, sim *simula.Simulator, log event.Stream) *Transport {
	return &Transport{
		Name:           name,
		transportMap:   transportMap,
		shipmentOption: m.Sea,
		sim:            sim,
		capacity:       1,
		log:            log,
	}
}

func (t *Transport) WithCapacity(cap int) *Transport {
	t.capacity = cap
	return t
}

func (t *Transport) WithLoadTime(dur time.Duration) *Transport {
	t.timeToLoad = dur
	return t
}

func (t *Transport) WithUnloadTime(dur time.Duration) *Transport {
	t.timeToUnload = dur
	return t
}
