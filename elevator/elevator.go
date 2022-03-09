package elevator

import (
	"log"

	"github.com/TTK4145-2022-students/project-group-78/door"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type State struct {
	Behaviour Behaviour
	Floor     int
	Direction elevio.MotorDirection
}

type Behaviour string

const (
	DoorOpen Behaviour = "doorOpen"
	Moving             = "moving"
	Idle               = "idle"
)

func Elevator(ordersC <-chan Orders, completedOrderC chan<- elevio.ButtonEvent, stateC chan<- State) {
	doorOpenC := make(chan bool)
	doorClosedC := make(chan bool)
	floorEnteredC := make(chan int)

	go door.Door(doorOpenC, doorClosedC)
	go elevio.PollFloorSensor(floorEnteredC)

	elevio.SetMotorDirection(elevio.MD_Down)
	state := State{Behaviour: Moving}
	var orders Orders

	for {
		select {
		case <-doorClosedC:
			switch state.Behaviour {
			case DoorOpen:
				d, b := nextAction(orders, state.Direction)
				if state.Direction == elevio.MD_Stop {
					state.Behaviour = Idle
				} else {
					elevio.SetMotorDirection(state.Direction)
					state.Behaviour = Moving
				}
				stateC <- state

			default:
				log.Panicf("door closed while %v", state.Behaviour)
			}

		case state.Floor = <-floorEnteredC:
			switch state.Behaviour {
			case Moving:
				if shouldStop(orders, state.Floor, state.Direction) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					state.Behaviour = DoorOpen
				}
				stateC <- state

			default:
				log.Panicf("elevator entered floor while %v", state.Behaviour)
			}

		case orders = <-ordersC:
			switch state.Behaviour {
			case Idle:
				state.Direction = chooseNextDirection(orders, state.Direction)
				if state.Direction == elevio.MD_Stop {
					doorOpenC <- true
					state.Behaviour = DoorOpen
				} else {
					elevio.SetMotorDirection(state.Direction)
					state.Behaviour = Moving
				}
				stateC <- state

			case Moving:

			case DoorOpen:

			default:
				log.Panicf("received new orders while %v", state.Behaviour)
			}
		}
	}
}

func nextAction(orders Orders, floor int, d elevio.MotorDirection) (elevio.MotorDirection, Behaviour) {
	switch d {
	case elevio.MD_Up:
		switch {
		case orders.Above(floor):
			return elevio.MD_Up, Moving

		case orders.Here(floor):
			return elevio.MD_Down, DoorOpen

		case orders.Below(floor):
			return elevio.MD_Down, Moving
		}

	case elevio.MD_Down:
		switch {
		case orders.Below(floor):
			return elevio.MD_Down, Moving

		case orders.Here(floor):
			return elevio.MD_Up, DoorOpen

		case orders.Above(floor):
			return elevio.MD_Up, Moving
		}

	case elevio.MD_Stop:
		switch {
		case orders.Here(floor):
			return elevio.MD_Stop, DoorOpen

		case orders.Above(floor):
			return elevio.MD_Up, Moving

		case orders.Below(floor):
			return elevio.MD_Down, Moving
		}
	}
	return elevio.MD_Stop, Idle
}

func shouldStop(orders Orders, floor int, direction elevio.MotorDirection) bool {
	return false
}

func clearOrders(orders Orders, floor int, completedOrderC chan<- elevio.ButtonEvent) {
	for bt := 0; bt < 3; bt++ {
		if orders[floor][bt] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(bt)}
		}
	}
}
