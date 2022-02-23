package events

type EventType int

const (
	ElevatorArrivedAtFloor EventType = iota
)

type Event struct {
	Type  EventType
	Floor int
}
