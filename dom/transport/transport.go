package transport

import (
	"time"

	m "github.com/digimortl/tycoon/dom/transmap"
	w "github.com/digimortl/tycoon/dom/warehouse"
	"github.com/digimortl/tycoon/simula"
)

type Transport struct {
	Name           string
	transportMap   *m.Map
	shipmentOption m.ShipmentOption
	sim            *simula.Simulator
	home           *w.Warehouse
	arriveAfter    time.Duration
	cargoes        []*w.Cargo
	capacity       int
	timeToLoad     time.Duration
	timeToUnload   time.Duration
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
	PrintEvent(Arrived{
		occurredAt:  t.holdFor(t.arriveAfter),
		transport:   t.Name,
		shipmentOpt: t.shipmentOption,
		atLocation:  t.home.Location,
	})

	t.loadCargoesFrom(t.home)
	PrintEvent(Loaded{
		occurredAt:  t.holdFor(t.timeToLoad),
		transport:   t.Name,
		shipmentOpt: t.shipmentOption,
		duration:    t.timeToLoad,
		cargoes:     append(t.cargoes),
	})

	itinerary := t.findItineraryFrom(t.home)
	PrintEvent(Departed{
		occurredAt:   t.holdFor(0),
		transport:    t.Name,
		shipmentOpt:  t.shipmentOption,
		fromLocation: t.home.Location,
		toLocation:   itinerary.Destination().Location,
		cargoes:      append(t.cargoes),
	})

	PrintEvent(Arrived{
		occurredAt:  t.holdFor(itinerary.TotalTimeToTravel()),
		transport:   t.Name,
		shipmentOpt: t.shipmentOption,
		atLocation:  itinerary.Destination().Location,
		cargoes:     append(t.cargoes),
	})

	unloadTime := t.holdFor(t.timeToUnload)
	t.unloadCargoesTo(itinerary.Destination())
	PrintEvent(Unloaded{
		occurredAt:  unloadTime,
		transport:   t.Name,
		shipmentOpt: t.shipmentOption,
		duration:    t.timeToUnload,
	})

	PrintEvent(Departed{
		occurredAt:   t.holdFor(0),
		transport:    t.Name,
		shipmentOpt:  t.shipmentOption,
		fromLocation: itinerary.Destination().Location,
		toLocation:   t.home.Location,
	})

	t.arriveAfter = itinerary.TotalTimeToTravel()
	t.Run()
}

func (t *Transport) StartJourneyFrom(warehouse *w.Warehouse) {
	t.home, t.arriveAfter = warehouse, 0
	t.sim.Spawn(t)
}

func (t *Transport) holdFor(duration time.Duration) time.Time {
	return t.sim.WakeUpAfter(duration)
}

func Truck(name string, transportMap *m.Map, sim *simula.Simulator) *Transport {
	return &Transport{
		Name:           name,
		transportMap:   transportMap,
		shipmentOption: m.Land,
		sim:            sim,
		capacity:       1,
	}
}

func Vessel(name string, transportMap *m.Map, sim *simula.Simulator) *Transport {
	return &Transport{
		Name:           name,
		transportMap:   transportMap,
		shipmentOption: m.Sea,
		sim:            sim,
		capacity:       1,
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
