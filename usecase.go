package main

import (
	"fmt"
	"github.com/digimortl/tycoon/dom/event"
	"os"
	"strings"
	"time"

	m "github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/transport"
	"github.com/digimortl/tycoon/dom/warehouse"
	"github.com/digimortl/tycoon/simula"
)

func format(trackNumber int) string {
	return fmt.Sprintf("%05d", trackNumber)
}

func newEventStream(handler func(event.DomainEvent)) event.Stream {
	es := make(event.Stream)
	go func() {
		defer close(es)
		for ev := range es {
			handler(ev)
		}
	}()
	return es
}

func UseCase1(destinationCodes ...warehouse.LocationCode) time.Duration {
	if len(destinationCodes) == 0 {
		return time.Duration(0)
	}

	sim := simula.NewSimulator()
	defer sim.Stop()

	factory := warehouse.Of("Factory", sim)
	defer factory.Stop()

	port := warehouse.Of("Port", sim)
	defer port.Stop()

	warehouseA := warehouse.Of("A", sim)
	defer warehouseA.Stop()

	warehouseB := warehouse.Of("B", sim)
	defer warehouseB.Stop()

	isDestinationValid := func(dest warehouse.LocationCode) bool {
		switch dest {
		case warehouseA.Location, warehouseB.Location:
			return true
		}
		return false
	}
	var cargoesToDeliver []*warehouse.Cargo = nil
	for trackNumber, destCode := range destinationCodes {
		if isDestinationValid(destCode) {
			cargoesToDeliver = append(cargoesToDeliver, &warehouse.Cargo{
				TrackNumber: format(trackNumber + 1),
				Origin:      factory.Location,
				Destination: destCode,
			})
		} else {
			panic(fmt.Sprintf("Invalid destination code %s", destCode))
		}
	}

	factory.Bring(cargoesToDeliver...)

	transportMap :=
		m.NewMap().
			ByLand(factory, port, time.Hour).
			BySea(port, warehouseA, 4*time.Hour).
			ByLand(factory, warehouseB, 5*time.Hour)

	nullEventStream := newEventStream(func(_ event.DomainEvent) {})

	transport.Truck("Truck 1", transportMap, sim, nullEventStream).StartJourneyFrom(factory)
	transport.Truck("Truck 2", transportMap, sim, nullEventStream).StartJourneyFrom(factory)
	transport.Vessel("Vessel 1", transportMap, sim, nullEventStream).StartJourneyFrom(port)

	tillCargoesHaveBeenDelivered := func() bool {
		return warehouseA.Fullness()+warehouseB.Fullness() == len(cargoesToDeliver)
	}
	return sim.Proceed(tillCargoesHaveBeenDelivered)
}

func UseCase2(destinationCodes ...warehouse.LocationCode) time.Duration {
	sim := simula.NewSimulator()
	defer sim.Stop()

	factory := warehouse.Of("Factory", sim)
	defer factory.Stop()

	port := warehouse.Of("Port", sim)
	defer port.Stop()

	warehouseA := warehouse.Of("A", sim)
	defer warehouseA.Stop()

	warehouseB := warehouse.Of("B", sim)
	defer warehouseB.Stop()

	isDestinationValid := func(dest warehouse.LocationCode) bool {
		switch dest {
		case warehouseA.Location, warehouseB.Location:
			return true
		}
		return false
	}
	var cargoesToDeliver []*warehouse.Cargo = nil
	for trackNumber, destCode := range destinationCodes {
		if isDestinationValid(destCode) {
			cargoesToDeliver = append(cargoesToDeliver, &warehouse.Cargo{
				TrackNumber: format(trackNumber + 1),
				Origin:      factory.Location,
				Destination: destCode,
			})
		} else {
			panic(fmt.Sprintf("Invalid destination code %s", destCode))
		}
	}

	factory.Bring(cargoesToDeliver...)

	transportMap :=
		m.NewMap().
			ByLand(factory, port, time.Hour).
			BySea(port, warehouseA, 6*time.Hour).
			ByLand(factory, warehouseB, 5*time.Hour)

	printingEventStream := newEventStream(printEvent)

	transport.Truck("Truck 1", transportMap, sim, printingEventStream).
		StartJourneyFrom(factory)

	transport.Truck("Truck 2", transportMap, sim, printingEventStream).
		StartJourneyFrom(factory)

	transport.Vessel("Vessel 1", transportMap, sim, printingEventStream).
		WithCapacity(4).
		WithLoadTime(1 * time.Hour).
		WithUnloadTime(1 * time.Hour).
		StartJourneyFrom(port)

	tillCargoesHaveBeenDelivered := func() bool {
		return warehouseA.Fullness()+warehouseB.Fullness() == len(cargoesToDeliver)
	}
	return sim.Proceed(tillCargoesHaveBeenDelivered)
}

func destCodes(args []string) []string {
	return strings.Split(strings.Join(args, ""), "")
}

func main() {
	switch strings.ToLower(os.Args[1]) {
	case "exercise-1":
		UseCase1(destCodes(os.Args[2:])...)
	case "exercise-2":
		UseCase2(destCodes(os.Args[2:])...)
	default:
		fmt.Println("expected 'Exercise-1' or 'Exercise-2'")
		os.Exit(1)
	}
}
