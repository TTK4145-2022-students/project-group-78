package elevator

import (
	"log"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/door"
	"github.com/vetleras/driver-go/elevio"
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
	doorOpenC := make(chan bool, config.CHAN_SIZE)
	doorClosedC := make(chan bool, config.CHAN_SIZE)
	floorEnteredC := make(chan int)

	go door.Door(doorOpenC, doorClosedC)
	go elevio.PollFloorSensor(floorEnteredC)

	elevio.SetMotorDirection(elevio.MD_Down)
	state := State{Behaviour: Moving, Direction: elevio.MD_Down}
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
			elevio.SetFloorIndicator(state.Floor)
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
			switch state.Behaviour {
			case Idle:
			case DoorOpen:
				doorOpenC <- true
				clearOrders(orders, state.Floor, state.Direction, completedOrderC)
				stateC <- state

			case Moving:
				elevio.SetMotorDirection(state.Direction)
				stateC <- state
			}
		}
	}
}
