package clouditor

import "sync"

type Channel chan Event

type EventBus struct {
	subscribers map[string][]Channel
	mutex       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: map[string][]Channel{},
	}
}

type Event struct {
	Data  interface{}
	Topic string
}

func (e *EventBus) Subscribe(topic string, ch Channel) {
	e.mutex.Lock()
	var (
		channels []Channel
		exists   bool
	)

	// retrieve previous channels, if they exist, otherwise create new slice
	if channels, exists = e.subscribers[topic]; !exists {
		channels = []Channel{}
	}

	// append our new channel
	e.subscribers[topic] = append(channels, ch)
	e.mutex.Unlock()
}

func (e *EventBus) Publish(topic string, data interface{}) {
	e.mutex.Lock()
	if channels, exists := e.subscribers[topic]; exists {
		go func(event Event, channels []Channel) {
			for _, channel := range channels {
				channel <- event
			}
		}(Event{Data: data, Topic: topic}, channels)
	}
	e.mutex.Unlock()
}
