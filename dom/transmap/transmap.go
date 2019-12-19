package transmap

import (
	"errors"
	"time"

	w "github.com/digimortl/tycoon/dom/warehouse"
)

type ShipmentOption = string

const (
	Land ShipmentOption = "land"
	Sea                 = "sea"
)

type Segment struct {
	Origin         *w.Warehouse
	Destination    *w.Warehouse
	TimeToTravel   time.Duration
	shipmentOption ShipmentOption
}

type Itinerary struct {
	Segments []Segment
}

func (it *Itinerary) Origin() *w.Warehouse {
	return it.Segments[0].Origin
}

func (it *Itinerary) Destination() *w.Warehouse {
	return it.Segments[len(it.Segments)-1].Destination
}

func (it *Itinerary) TotalTimeToTravel() time.Duration {
	var total time.Duration
	for _, seg := range it.Segments {
		total += seg.TimeToTravel
	}
	return total
}

func (it *Itinerary) WithOption(option ShipmentOption) Itinerary {
	var segments []Segment
	for _, seg := range it.Segments {
		if seg.shipmentOption == option {
			segments = append(segments, seg)
		} else {
			break
		}
	}
	return Itinerary{Segments: segments}
}

type Map struct {
	graph map[w.LocationCode]map[w.LocationCode]Segment
}

func NewMap() *Map {
	return &Map{graph: make(map[w.LocationCode]map[w.LocationCode]Segment)}
}

func (m *Map) FindItinerary(origin, destination w.LocationCode) (Itinerary, error) {
	var find func(w.LocationCode, []Segment) []Segment

	find = func(orig w.LocationCode, path []Segment) []Segment {
		if _, ok := m.graph[orig][destination]; ok {
			return append(path, m.graph[orig][destination])
		}

		var itinerary []Segment
		for loc := range m.graph[orig] {
			if len(path) > 0 {
				if path[len(path)-1].Origin.Location == loc {
					continue
				}
			}
			itinerary = find(loc, append(path, m.graph[orig][loc]))
			if len(itinerary) > 0 {
				break
			}
		}
		return itinerary
	}

	if _, ok := m.graph[origin]; !ok {
		return Itinerary{}, errors.New("Itinerary not found")
	}

	return Itinerary{find(origin, nil)}, nil
}

func (m *Map) link(loc1, loc2 *w.Warehouse, ttt time.Duration, opt ShipmentOption) *Map {

	putSegment := func(orig, dest w.LocationCode, seg Segment) {
		if _, ok := m.graph[orig]; !ok {
			m.graph[orig] = make(map[w.LocationCode]Segment)
		}
		m.graph[orig][dest] = seg
	}

	putSegment(loc1.Location, loc2.Location, Segment{loc1, loc2, ttt, opt})
	putSegment(loc2.Location, loc1.Location, Segment{loc2, loc1, ttt, opt})

	return m
}

func (m *Map) ByLand(origin, destination *w.Warehouse, timeToTravel time.Duration) *Map {
	return m.link(origin, destination, timeToTravel, Land)
}

func (m *Map) BySea(origin, destination *w.Warehouse, timeToTravel time.Duration) *Map {
	return m.link(origin, destination, timeToTravel, Sea)
}
