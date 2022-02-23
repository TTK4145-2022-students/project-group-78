package events

import (
	"testing"
)

func TestTopic(t *testing.T) {
	topic := Topic{}

	c := make(chan Event)
	topic.Subscribe(c)
	go topic.Publish(Event{Type: ElevatorArrivedAtFloor})
	t.Error(<-c)
}
