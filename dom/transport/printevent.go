package transport

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/warehouse"
)

func PrintEvent(anEvent DomainEvent) {
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
	case Arrived:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Location    string  `json:"location"`
			Cargo       []cargo `json:"cargo"`
		}{
			toRelTime(anEvent.(Arrived).occurredAt),
			"ARRIVE",
			toKind(anEvent.(Arrived).shipmentOpt),
			anEvent.(Arrived).transport,
			anEvent.(Arrived).atLocation,
			toCargoes(anEvent.(Arrived).cargoes),
		}
	case Departed:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Location    string  `json:"location"`
			Destination string  `json:"destination"`
			Cargo       []cargo `json:"cargo"`
		}{
			toRelTime(anEvent.(Departed).occurredAt),
			"DEPART",
			toKind(anEvent.(Departed).shipmentOpt),
			anEvent.(Departed).transport,
			anEvent.(Departed).fromLocation,
			anEvent.(Departed).toLocation,
			toCargoes(anEvent.(Departed).cargoes),
		}
	case Loaded:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Duration    float64 `json:"duration"`
			Cargo       []cargo `json:"cargo"`
		}{
			toRelTime(anEvent.(Loaded).occurredAt.Add(
				-anEvent.(Loaded).duration)),
			"LOAD",
			toKind(anEvent.(Loaded).shipmentOpt),
			anEvent.(Loaded).transport,
			anEvent.(Loaded).duration.Hours(),
			toCargoes(anEvent.(Loaded).cargoes),
		}
	case Unloaded:
		s = struct {
			Time        float64 `json:"time"`
			Event       string  `json:"event"`
			Kind        string  `json:"kind"`
			TarnsportId string  `json:"transport_id"`
			Duration    float64 `json:"duration"`
		}{
			toRelTime(anEvent.(Unloaded).occurredAt.Add(-
				anEvent.(Unloaded).duration)),
			"UNLOAD",
			toKind(anEvent.(Unloaded).shipmentOpt),
			anEvent.(Unloaded).transport,
			anEvent.(Unloaded).duration.Hours(),
		}
	default:
		s = nil
	}
	jsonData, err := json.Marshal(s)
	if err == nil {
		fmt.Println(string(jsonData))
	}
}
