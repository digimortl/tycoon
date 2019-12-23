package transport

import (
	"time"

	"github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/warehouse"
)

type Arrived struct {
	OccurredAt  time.Time
	Transport   string
	ShipmentOpt transmap.ShipmentOption
	AtLocation  warehouse.LocationCode
	Cargoes     []*warehouse.Cargo
}

type Departed struct {
	OccurredAt   time.Time
	Transport    string
	ShipmentOpt  transmap.ShipmentOption
	FromLocation warehouse.LocationCode
	ToLocation   warehouse.LocationCode
	Cargoes      []*warehouse.Cargo
}

type Loaded struct {
	OccurredAt  time.Time
	Transport   string
	ShipmentOpt transmap.ShipmentOption
	Cargoes     []*warehouse.Cargo
	Duration    time.Duration
}

type Unloaded struct {
	OccurredAt  time.Time
	Transport   string
	ShipmentOpt transmap.ShipmentOption
	Duration    time.Duration
}
