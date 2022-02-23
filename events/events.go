package events

import "encoding/gob"

func init() {
	gob.Register(DoorOpened{})
	gob.Register(DoorOpened{})
	gob.Register(OrderReceived{})
	gob.Register(OrderServed{})
	gob.Register(EnteredFloor{})
}

type DoorOpened struct{}

type DoorClosed struct{}

type OrderReceived struct{}

type OrderServed struct{}

type EnteredFloor struct{}
