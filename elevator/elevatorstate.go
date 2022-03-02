package elevator

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type CentralState struct {
	HallUpOrders   []Order
	HallDownOrders []Order
	Elevators      []ElevatorState
}

type ElevatorState struct {
	State     State
	Direction elevio.MotorDirection
	Floor     int
	CabOrders []Order
}

type Order struct {
	Time   time.Time
	Served bool
}
