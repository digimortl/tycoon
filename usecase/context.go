package usecase

import (
	"github.com/digimortl/tycoon/dom/event"
	m "github.com/digimortl/tycoon/dom/transmap"
	"github.com/digimortl/tycoon/dom/transport"
	"github.com/digimortl/tycoon/dom/warehouse"
	"github.com/digimortl/tycoon/simula"
)

type Context struct {
	sim         *simula.Simulator
	eventStream event.Stream
}

func NewContext() *Context {
	return &Context{
		sim:         simula.NewSimulator(),
		eventStream: event.NullEventStream(),
	}
}

func (ctx *Context) Close() {
	ctx.sim.Stop()
}

func (ctx *Context) Simulator() *simula.Simulator {
	return ctx.sim
}

func (ctx *Context) WithEventHandler(handler func(domainEvent event.DomainEvent)) *Context {
	ctx.eventStream = event.NewEventStream(handler)
	return ctx
}

func (ctx *Context) Truck(name string, transportMap *m.Map) *transport.Transport {
	return transport.Truck(name, transportMap, ctx.sim, ctx.eventStream)
}

func (ctx *Context) Vessel(name string, transportMap *m.Map) *transport.Transport {
	return transport.Vessel(name, transportMap, ctx.sim, ctx.eventStream)
}

func (ctx *Context) WarehouseOf(location warehouse.LocationCode) *warehouse.Warehouse {
	return warehouse.Of(location, ctx.sim)
}
