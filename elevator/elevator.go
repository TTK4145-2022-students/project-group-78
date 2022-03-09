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
			if state.Behaviour != DoorOpen {
				log.Panicf("door closed while %v", state.Behaviour)
			}
			state.Direction, state.Behaviour = nextAction(orders, state.Floor, state.Direction)
			switch state.Behaviour {
			case Idle:
			case DoorOpen:
				doorOpenC <- true
				clearOrders(orders, state.Floor, state.Direction, completedOrderC)

			case Moving:
				elevio.SetMotorDirection(state.Direction)
			}
			stateC <- state

		case state.Floor = <-floorEnteredC:
			if state.Behaviour != Moving {
				log.Panicf("elevator entered floor while %v", state.Behaviour)

			}
			if shouldStop(orders, state.Floor, state.Direction) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				doorOpenC <- true
				clearOrders(orders, state.Floor, state.Direction, completedOrderC)
				state.Behaviour = DoorOpen
			}
			stateC <- state

		case orders = <-ordersC:
			if state.Behaviour != Idle {
				continue
			}

			state.Direction, state.Behaviour = nextAction(orders, state.Floor, state.Direction)
			if state.Direction == elevio.MD_Stop {
				doorOpenC <- true
				state.Behaviour = DoorOpen
			} else {
				elevio.SetMotorDirection(state.Direction)
				state.Behaviour = Moving
			}
			stateC <- state

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
	switch direction {
	case elevio.MD_Down:
		return orders[floor][elevio.BT_HallDown] || orders[floor][elevio.BT_Cab] || !orders.Below(floor)
	case elevio.MD_Up:
		return orders[floor][elevio.BT_HallUp] || orders[floor][elevio.BT_Cab] || !orders.Above(floor)
	case elevio.MD_Stop:
		log.Panicf("Direction is  %v, when expected to be up or down", direction)
		return true
	default:
		log.Panicf("Direction is corrupted %v", direction)
		return false
	}
}

func clearOrders(orders Orders, floor int, direction elevio.MotorDirection, completedOrderC chan<- elevio.ButtonEvent) {
	if orders[floor][elevio.BT_Cab] {
		completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}

	if orders[floor][direction] {
		completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(direction)}
	}
}
