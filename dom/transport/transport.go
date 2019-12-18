package transport

import (
	"fmt"
	"strings"
	"time"

	m "github.com/digimortl/tycoon/dom/transmap"
	w "github.com/digimortl/tycoon/dom/warehouse"
	"github.com/digimortl/tycoon/simula"
)

type Transport struct {
	Name string
	transportMap *m.Map
	shipmentOption m.ShipmentOption
	sim *simula.Simulator
	home *w.Warehouse
	arriveAfter time.Duration
	cargoes []*w.Cargo
	capacity int
	timeToLoad time.Duration
	timeToUnload time.Duration
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

func formatCargoes(cargoes []*w.Cargo) string {
	var xs []string
	for _, aCargo := range cargoes {
		xs = append(xs, fmt.Sprintf("%s(%s->%s)", aCargo.TrackNumber, aCargo.Origin, aCargo.Destination))
	}
	return strings.Join(xs, ",")
}

func (t *Transport) Run() {
	arrivalTime := t.holdFor(t.arriveAfter)
	fmt.Printf("%s arrived at %s at %s\n", t.Name, t.home.Location, arrivalTime)

	t.loadCargoesFrom(t.home)
	loadTime := t.holdFor(t.timeToLoad)
	fmt.Printf("%s loaded cargoes %s at %s from %s\n", t.Name, formatCargoes(t.cargoes), loadTime, t.home.Location)

	itinerary := t.findItineraryFrom(t.home)

	departureTime := t.holdFor(0)
	fmt.Printf("%s departed at %s from %s to %s\n", t.Name, departureTime, t.home.Location, itinerary.Destination().Location)

	arrivalTime = t.holdFor(itinerary.TotalTimeToTravel())
	fmt.Printf("%s arrived at %s at %s\n", t.Name, arrivalTime, itinerary.Destination().Location)

	unloadTime := t.holdFor(t.timeToUnload)
	t.unloadCargoesTo(itinerary.Destination())
	fmt.Printf("%s unloaded cargoes %s at %s to %s\n", t.Name, formatCargoes(t.cargoes), unloadTime, itinerary.Destination().Location)

	departureTime = t.holdFor(0)
	fmt.Printf("%s Departed at %s from %s to %s \n", t.Name, departureTime, itinerary.Destination().Location, t.home.Location)

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
		Name: name,
		transportMap: transportMap,
		shipmentOption: m.Land,
		sim: sim,
		home: nil,
		arriveAfter: 0,
		cargoes: nil,
		capacity: 1,
		timeToLoad: 0 * time.Hour,
		timeToUnload: 0 * time.Hour,

	}
}

func Vessel(name string, transportMap *m.Map, sim *simula.Simulator) *Transport {
	return &Transport{
		Name: name,
		transportMap: transportMap,
		shipmentOption: m.Sea,
		sim: sim,
		home: nil,
		arriveAfter: 0,
		cargoes: nil,
		capacity: 1,
		timeToLoad: 0 * time.Hour,
		timeToUnload: 0 * time.Hour,
	}
}
