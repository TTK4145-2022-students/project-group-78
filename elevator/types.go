package elevator

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type CentralState struct {
	Origin       int
	Orders       map[Order]time.Time
	ServedOrders map[Order]time.Time
	Elevators    [config.NUM_ELEVATORS]ElevatorState
}

type ElevatorState struct {
	State     State
	Direction elevio.MotorDirection
	Floor     int
	CabOrders [config.NUM_FLOORS]bool
}

type Order struct {
	Type  elevio.ButtonType
	Floor int
}

func NewCentralState(id int, es ElevatorState) CentralState {
	cs := CentralState{
		Origin:       id,
		Orders:       make(map[Order]time.Time, 1),
		ServedOrders: make(map[Order]time.Time, 1),
	}
	cs.Elevators[id] = es
	return cs
}

// Merge cs2 onto cs1
func MergeCentralState(cs1 CentralState, cs2 CentralState) {
	cs1.Elevators[cs2.Origin] = cs2.Elevators[cs2.Origin]

	for o, _ := range cs2.Orders {
		if cs2.Orders[o].After(cs1.Orders[o]) {
			cs1.Orders[o] = cs2.Orders[o]
		}
	}

	for o, _ := range cs2.ServedOrders {
		if cs2.ServedOrders[o].After(cs1.ServedOrders[o]) {
			cs1.ServedOrders[o] = cs2.ServedOrders[o]
		}
	}
}
