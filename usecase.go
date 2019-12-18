package main

import (
	"fmt"
	"time"

	m "github.com/digimortl/tycoon/dom/transmap"
	t "github.com/digimortl/tycoon/dom/transport"
	"github.com/digimortl/tycoon/dom/warehouse"
	"github.com/digimortl/tycoon/simula"
)

func format(trackNumber int) string {
	return fmt.Sprintf("%05d", trackNumber)
}

func UseCase(destinationCodes ...warehouse.LocationCode) time.Duration {
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
				TrackNumber: format(trackNumber+1),
				Origin: factory.Location,
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
			BySea(port, warehouseA, 4 * time.Hour).
			ByLand(factory, warehouseB, 5 * time.Hour)

	t.Truck("Truck 1", transportMap, sim).StartJourneyFrom(factory)
	t.Truck("Truck 2", transportMap, sim).StartJourneyFrom(factory)
	t.Vessel("Vessel 1", transportMap, sim).StartJourneyFrom(port)

	tillCargoesHaveBeenDelivered := func () bool {
		return warehouseA.Fullness() + warehouseB.Fullness() == len(cargoesToDeliver)
	}
	sim.Proceed(tillCargoesHaveBeenDelivered)

	return sim.CurrentTime.Sub(time.Time{})
}

func main() {
	//fmt.Println(UseCase("A")) // 5
	//fmt.Println(UseCase("A", "B")) // 5
	//fmt.Println(UseCase("B", "B")) // 5
	//fmt.Println(UseCase("A", "B", "B")) // 7
	//fmt.Println(UseCase("A", "A", "B", "A", "B", "B", "A", "B")) // 29
	fmt.Println(UseCase("A", "B", "B", "B", "A", "B", "A", "A", "A", "B", "B", "B")) // 41
}


