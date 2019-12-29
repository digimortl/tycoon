package main

import (
	"encoding/json"
	"fmt"
	"github.com/digimortl/tycoon/dom/warehouse"
	"os"
	"time"

	"github.com/digimortl/tycoon/dom/event"
	"github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/transport"
)

type (
	eventMap = map[string]interface{}
	cargoMap = map[string]string
)

func toRelTime(tm time.Time) float64 {
	return tm.Sub(time.Time{}).Hours()
}

func toKind(opt transmap.ShipmentOption) string {
	switch opt {
	case transmap.Land:
		return "TRUCK"
	case transmap.Sea:
		return "VESSEL"
	default:
		return "UNKNOWN"
	}
}

func toCargoMaps(cargoes []*warehouse.Cargo) []cargoMap {
	cs := make([]cargoMap, 0)
	for _, c := range cargoes {
		cs = append(cs, cargoMap{
			"cargo_id":    c.TrackNumber,
			"origin":      c.Origin,
			"destination": c.Destination,
		})
	}
	return cs
}

func toEventMap(anEvent event.DomainEvent) (e eventMap) {
	switch anEvent.(type) {
	case transport.Arrived:
		e = eventMap{
			"time":         toRelTime(anEvent.(transport.Arrived).OccurredAt),
			"event":        "ARRIVE",
			"kind":         toKind(anEvent.(transport.Arrived).ShipmentOpt),
			"transport_id": anEvent.(transport.Arrived).Transport,
			"location":     anEvent.(transport.Arrived).AtLocation,
			"cargo":        toCargoMaps(anEvent.(transport.Arrived).Cargoes),
		}
	case transport.Departed:
		e = eventMap{
			"time":         toRelTime(anEvent.(transport.Departed).OccurredAt),
			"event":        "DEPART",
			"kind":         toKind(anEvent.(transport.Departed).ShipmentOpt),
			"transport_id": anEvent.(transport.Departed).Transport,
			"location":     anEvent.(transport.Departed).FromLocation,
			"destination":  anEvent.(transport.Departed).ToLocation,
			"cargo":        toCargoMaps(anEvent.(transport.Departed).Cargoes),
		}
	case transport.Loaded:
		e = eventMap{
			"time": toRelTime(anEvent.(transport.Loaded).OccurredAt.Add(
				-anEvent.(transport.Loaded).Duration)),
			"event":        "LOAD",
			"kind":         toKind(anEvent.(transport.Loaded).ShipmentOpt),
			"transport_id": anEvent.(transport.Loaded).Transport,
			"duration":     anEvent.(transport.Loaded).Duration.Hours(),
			"cargo":        toCargoMaps(anEvent.(transport.Loaded).Cargoes),
		}
	case transport.Unloaded:
		e = eventMap{
			"time":         toRelTime(anEvent.(transport.Unloaded).OccurredAt.Add(-anEvent.(transport.Unloaded).Duration)),
			"event":        "UNLOAD",
			"kind":         toKind(anEvent.(transport.Unloaded).ShipmentOpt),
			"transport_id": anEvent.(transport.Unloaded).Transport,
			"duration":     anEvent.(transport.Unloaded).Duration.Hours(),
		}
	}
	return
}

func printEventToStdout(anEvent event.DomainEvent) {
	e := toEventMap(anEvent)
	if e != nil {
		jsonData, err := json.Marshal(e)
		if err == nil {
			fmt.Fprintln(os.Stdout, string(jsonData))
		}
	}
}
