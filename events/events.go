package events

import (
	"github.com/TTK4145-2022-students/project-group-78/order"
)

type DoorOpened struct{}

type DoorClosed struct{}

type OrderReceived struct {
	Order order.Order
}

type OrderServed struct {
	Order order.Order
}

type FloorEntered struct {
	Floor int
}
