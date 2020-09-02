package event

import (
	"errors"
	"sync"
)

//BUS public var of evnet bus
var BUS = NewBus()

//Listener define listener handler
type Listener func(e Event, payload ...interface{})

//Event interface
type Event interface {
	GetEventName() string
}

//SimpleEvent struct
type SimpleEvent struct {
	name string
}

//GetEventName return event name
func (e *SimpleEvent) GetEventName() string {
	return e.name
}

//Bus event bus struct
type Bus struct {
	listeners map[string][]Listener
	m         *sync.Mutex
	wg        *sync.WaitGroup
}

//Dispatch event bus
func (b *Bus) Dispatch(event interface{}, payload ...interface{}) error {
	e, err := b.resolveEvent(event)
	b.m.Lock()
	defer b.m.Unlock()
	if err == nil {
		if listeners, exists := b.listeners[e.GetEventName()]; exists {
			b.wg.Add(len(listeners))
			for _, listener := range listeners {
				go func(listener Listener) {
					listener(e, payload...)
					defer b.wg.Done()
				}(listener)
			}
			b.wg.Wait()
		}
	}

	return err
}

//Listen event
func (b *Bus) Listen(name string, listener Listener) {
	b.m.Lock()
	defer b.m.Unlock()
	b.listeners[name] = append(b.listeners[name], listener)
}

func (b *Bus) resolveEvent(input interface{}) (Event, error) {
	if e, ok := input.(Event); ok {
		return e, nil
	}

	if name, ok := input.(string); ok {
		return NewSimpleEvent(name), nil
	}

	return nil, errors.New("input event should an instance of event.Event or string")
}

// NewSimpleEvent return simple event instance
func NewSimpleEvent(name string) *SimpleEvent {
	return &SimpleEvent{name: name}
}

//NewBus constructor of event bus
func NewBus() *Bus {
	return &Bus{
		listeners: map[string][]Listener{},
		m:         &sync.Mutex{},
		wg:        &sync.WaitGroup{},
	}
}
