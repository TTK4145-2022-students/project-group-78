package elevator

import (
	"log"

	"github.com/TTK4145-2022-students/project-group-78/config"
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

func Elevator(ordersC <-chan [config.NUM_FLOORS][3]bool, completedOrderC chan<- elevio.ButtonEvent, stateC chan<- State) {
	doorOpenC := make(chan bool)
	doorClosedC := make(chan bool)
	floorEnteredC := make(chan int)

	go door.Door(doorOpenC, doorClosedC)
	go elevio.PollFloorSensor(floorEnteredC)

	elevio.SetMotorDirection(elevio.MD_Down)
	state := State{Behaviour: Moving}
	var orders [config.NUM_FLOORS][3]bool

	for {
		select {
		case <-doorClosedC:
			switch state.Behaviour {
			case DoorOpen:
				state.Direction = chooseNextDirection(orders, state.Direction)
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
				if chooseNextDirection(orders, state.Direction) == elevio.MD_Stop {
					doorOpenC <- true
				}

			default:
				log.Panicf("received new orders while %v", state.Behaviour)
			}
		}
	}
}

func chooseNextDirection(orders [config.NUM_FLOORS][3]bool, d elevio.MotorDirection) elevio.MotorDirection {
	return elevio.MD_Stop
}

func shouldStop(orders [config.NUM_FLOORS][3]bool, floor int, direction elevio.MotorDirection) bool {
	switch direction {
	case elevio.MD_Down:
		return orders[floor][elevio.BT_HallDown] || orders[floor][elevio.BT_Cab] || !orders.below
	case elevio.MD_Up:
		return orders[floor][elevio.BT_HallUp] || orders[floor][elevio.BT_Cab] || !orders.above
	case elevio.MD_Stop:
		log.Panicf("Direction is  %v, when expected to be up or down", direction)
		return true
	default:
		log.Panicf("Direction is corrupted %v", direction)
		return false
	}
}

func clearOrders(orders [config.NUM_FLOORS][3]bool, floor int, direction elevio.MotorDirection, completedOrderC chan<- elevio.ButtonEvent) {

}
