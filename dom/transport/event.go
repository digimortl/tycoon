package transport

import (
	"time"

	"github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/warehouse"
)

type DomainEvent interface{}

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
