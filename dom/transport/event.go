package transport

import (
	"fmt"
	"github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/warehouse"
	"time"
)

type DomainEvent interface {}

type Arrived struct {
	occurredAt  time.Time
	transport   string
	shipmentOpt transmap.ShipmentOption
	atLocation  warehouse.LocationCode
	cargoes     []*warehouse.Cargo
}

type Departed struct {
	occurredAt   time.Time
	transport    string
	shipmentOpt  transmap.ShipmentOption
	fromLocation warehouse.LocationCode
	toLocation   warehouse.LocationCode
	cargoes      []*warehouse.Cargo
}

type Loaded struct {
	occurredAt  time.Time
	transport   string
	shipmentOpt transmap.ShipmentOption
	cargoes     []*warehouse.Cargo
	duration    time.Duration
}

type Unloaded struct {
	occurredAt  time.Time
	transport   string
	shipmentOpt transmap.ShipmentOption
	duration    time.Duration
}

func PrintEvent(anEvent DomainEvent) {
	switch anEvent.(type) {
	case Arrived:
		fmt.Printf("Arrived%+v\n", anEvent)
	case Departed:
		fmt.Printf("Departed%+v\n", anEvent)
	case Loaded:
		fmt.Printf("Loaded%+v\n", anEvent)
	case Unloaded:
		fmt.Printf("Unloaded%+v\n", anEvent)
	default:
		fmt.Printf("Unknown%+v\n", anEvent)
	}
}