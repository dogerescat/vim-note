package util

import "sync"

type EventType int

type Events map[EventType]interface{}

type EventBox struct {
	events Events
	cond   *sync.Cond
	ignore map[EventType]bool
}

func NewEventBox() *EventBox {
	return &EventBox{
		events: make(Events),
		cond:   sync.NewCond(&sync.Mutex{}),
		ignore: make(map[EventType]bool)}
}

func (b *EventBox) Wait(callback func(*Events)) {
	b.cond.L.Lock()

	if len(b.events) == 0 {
		b.cond.Wait()
	}

	callback(&b.events)
	b.cond.L.Unlock()
}

func (b *EventBox) Set(event EventType, value interface{}) {
	b.cond.L.Lock()
	b.events[event] = value
	if _, found := b.ignore[event]; !found {
		b.cond.Broadcast()
	}
	b.cond.L.Unlock()
}

func (events *Events) Clear() {
	for event := range *events {
		delete(*events, event)
	}
}
