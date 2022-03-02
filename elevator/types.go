package elevator

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type RelativePosition int

const (
	Above   RelativePosition = 1
	Below                    = -1
	AtFloor                  = 0
)

type State int

const (
	DoorOpen State = iota
	Moving
	Idle
)

type CentralState struct {
	Origin           int
	HallOrders       map[Order]time.Time
	ServedHallOrders map[Order]time.Time
	Elevators        [config.NUM_ELEVATORS]ElevatorState
}

type ElevatorState struct {
	State            State
	Direction        elevio.MotorDirection
	Floor            int
	CabOrders        [config.NUM_FLOORS]bool
	RelativePosition RelativePosition
	Target           Order
}

type Order struct {
	Type  elevio.ButtonType
	Floor int
}

func NewCentralState(id int, es ElevatorState) CentralState {
	cs := CentralState{
		Origin:           id,
		HallOrders:       make(map[Order]time.Time, 1),
		ServedHallOrders: make(map[Order]time.Time, 1),
	}
	cs.Elevators[id] = es
	return cs
}
