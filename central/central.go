package central

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type State string

const (
	DoorOpen     State = "doorOpen"
	ServingOrder       = "servingOrder"
	Idle               = "idle"
)

type CentralState struct {
	Origin     int
	HallOrders [config.NUM_FLOORS]struct {
		Up   Order
		Down Order
	}
	Elevators [config.NUM_ELEVATORS]ElevatorState
}

type ElevatorState struct {
	State     State
	Direction elevio.MotorDirection
	Floor     int
	CabOrders [config.NUM_FLOORS]bool
}

type Order struct {
	Active bool
	Time   time.Time
}

func (cs *CentralState) Merge(newCs CentralState) {
	cs.Elevators[newCs.Origin] = newCs.Elevators[newCs.Origin]
	for i := 0; i < len(cs.HallOrders); i++ {
		if cs.HallOrders[i].Up.Time.Before(newCs.HallOrders[i].Up.Time) {
			cs.HallOrders[i].Up = newCs.HallOrders[i].Up
		}
		if cs.HallOrders[i].Down.Time.Before(newCs.HallOrders[i].Down.Time) {
			cs.HallOrders[i].Down = newCs.HallOrders[i].Down
		}
	}
}
