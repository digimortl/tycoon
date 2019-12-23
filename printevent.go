package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/digimortl/tycoon/dom/event"
	"github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/transport"
	"github.com/digimortl/tycoon/dom/warehouse"
)

func printEventToStdout(anEvent event.DomainEvent) {
	toRelTime := func(tm time.Time) float64 {
		return tm.Sub(time.Time{}).Hours()
	}
	type cargo struct {
		CargoId     string `json:"cargo_id"`
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
	}
	toCargoes := func(cargoes []*warehouse.Cargo) []cargo {
		cs := make([]cargo, 0)
		for _, c := range cargoes {
			cs = append(cs, cargo{c.TrackNumber, c.Origin, c.Destination})
		}
		return cs
	}
	toKind := func(opt transmap.ShipmentOption) string {
		switch opt {
		case transmap.Land:
			return "TRUCK"
		case transmap.Sea:
			return "VESSEL"
		default:
			return "UNKNOWN"
		}
	}

	var s interface{}
	switch anEvent.(type) {
	case transport.Arrived:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Location    string  `json:"location"`
			Cargo       []cargo `json:"cargo"`
		}{
			toRelTime(anEvent.(transport.Arrived).OccurredAt),
			"ARRIVE",
			toKind(anEvent.(transport.Arrived).ShipmentOpt),
			anEvent.(transport.Arrived).Transport,
			anEvent.(transport.Arrived).AtLocation,
			toCargoes(anEvent.(transport.Arrived).Cargoes),
		}
	case transport.Departed:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Location    string  `json:"location"`
			Destination string  `json:"destination"`
			Cargo       []cargo `json:"cargo"`
		}{
			toRelTime(anEvent.(transport.Departed).OccurredAt),
			"DEPART",
			toKind(anEvent.(transport.Departed).ShipmentOpt),
			anEvent.(transport.Departed).Transport,
			anEvent.(transport.Departed).FromLocation,
			anEvent.(transport.Departed).ToLocation,
			toCargoes(anEvent.(transport.Departed).Cargoes),
		}
	case transport.Loaded:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Duration    float64 `json:"duration"`
			Cargo       []cargo `json:"cargo"`
		}{
			toRelTime(anEvent.(transport.Loaded).OccurredAt.Add(
				-anEvent.(transport.Loaded).Duration)),
			"LOAD",
			toKind(anEvent.(transport.Loaded).ShipmentOpt),
			anEvent.(transport.Loaded).Transport,
			anEvent.(transport.Loaded).Duration.Hours(),
			toCargoes(anEvent.(transport.Loaded).Cargoes),
		}
	case transport.Unloaded:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Duration    float64 `json:"duration"`
		}{
			toRelTime(anEvent.(transport.Unloaded).OccurredAt.Add(-anEvent.(transport.Unloaded).Duration)),
			"UNLOAD",
			toKind(anEvent.(transport.Unloaded).ShipmentOpt),
			anEvent.(transport.Unloaded).Transport,
			anEvent.(transport.Unloaded).Duration.Hours(),
		}
	default:
		s = nil
	}
	jsonData, err := json.Marshal(s)
	if err == nil {
		fmt.Fprintln(os.Stdout, string(jsonData))
	}
}
