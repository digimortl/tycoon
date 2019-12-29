package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	m "github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/warehouse"
	"github.com/digimortl/tycoon/usecase"
)

func format(trackNumber int) string {
	return fmt.Sprintf("%05d", trackNumber)
}

func UseCase1(destinationCodes ...warehouse.LocationCode) time.Duration {
	ctx := usecase.NewContext()
	defer ctx.Close()

	factory := ctx.WarehouseOf("Factory")
	defer factory.Stop()

	port := ctx.WarehouseOf("Port")
	defer port.Stop()

	warehouseA := ctx.WarehouseOf("A")
	defer warehouseA.Stop()

	warehouseB := ctx.WarehouseOf("B")
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

	ctx.Truck("Truck 1", transportMap).
			StartJourneyFrom(factory)
	ctx.Truck("Truck 2", transportMap).
			StartJourneyFrom(factory)
	ctx.Vessel("Vessel 1", transportMap).
			StartJourneyFrom(port)

	tillCargoesHaveBeenDelivered := func() bool {
		return warehouseA.Fullness()+warehouseB.Fullness() == len(cargoesToDeliver)
	}
	return ctx.Simulator().Proceed(tillCargoesHaveBeenDelivered)
}

func UseCase2(destinationCodes ...warehouse.LocationCode) time.Duration {
	ctx := usecase.NewContext().
		WithEventHandler(printEventToStdout)
	defer ctx.Close()

	factory := ctx.WarehouseOf("Factory")
	defer factory.Stop()

	port := ctx.WarehouseOf("Port")
	defer port.Stop()

	warehouseA := ctx.WarehouseOf("A")
	defer warehouseA.Stop()

	warehouseB := ctx.WarehouseOf("B")
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



	ctx.Truck("Truck 1", transportMap).
		StartJourneyFrom(factory)

	ctx.Truck("Truck 2", transportMap).
		StartJourneyFrom(factory)

	ctx.Vessel("Vessel 1", transportMap).
		WithCapacity(4).
		WithLoadTime(1 * time.Hour).
		WithUnloadTime(1 * time.Hour).
		StartJourneyFrom(port)

	tillCargoesHaveBeenDelivered := func() bool {
		return warehouseA.Fullness()+warehouseB.Fullness() == len(cargoesToDeliver)
	}
	// NOTE The last event can be unhandled due to this return can cause an immediate return from the main program.
	// NOTE So it'd be nice to implement graceful termination of event stream goroutine.
	return ctx.Simulator().Proceed(tillCargoesHaveBeenDelivered)
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
