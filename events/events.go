package events

import "encoding/gob"

func init() {
	gob.Register(DoorOpened{})
	gob.Register(DoorOpened{})
	gob.Register(OrderReceived{})
	gob.Register(OrderServed{})
	gob.Register(EnteredFloor{})
	gob.Register(DoorObstructed{})
	gob.Register(DoorUnobstructed{})
}

type DoorOpened struct{}

type DoorClosed struct{}

type OrderReceived struct{}

type OrderServed struct{}

type EnteredFloor struct{}

type DoorObstructed struct{}

type DoorUnobstructed struct{}
