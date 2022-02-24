package events

import (
	"encoding/gob"

	"github.com/TTK4145-2022-students/project-group-78/order"
)

func init() {
	gob.Register(DoorOpened{})
	gob.Register(DoorOpened{})
	gob.Register(OrderReceived{})
	gob.Register(OrderServed{})
	gob.Register(EnteredFloor{})
}

type DoorOpened struct{}

type DoorClosed struct{}

type OrderReceived struct {
	Order order.Order
}

type OrderServed struct {
	Order order.Order
}

type EnteredFloor struct{}
