package elevator

import (
	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/door"
)

type State struct {
	Behaviour Behaviour
	Floor     int
}

type Behaviour int

const (
	Idle Behaviour = iota
	DoorOpenUp
	DoorOpenDown
	MovingUp
	MovingDown
)

func Elevator(ordersC <-chan Orders, completedOrderC chan<- elevio.ButtonEvent, stateC chan<- State) {
	doorOpenC := make(chan bool, config.ChanSize)
	doorClosedC := make(chan bool, config.ChanSize)
	floorEnteredC := make(chan int)

	go door.Door(doorOpenC, doorClosedC)
	go elevio.PollFloorSensor(floorEnteredC)

	elevio.SetMotorDirection(elevio.MD_Down)
	state := State{}
	var orders Orders

	for {
		select {
		case <-doorClosedC:
			switch state.Behaviour {
			case DoorOpenUp:
				switch {
				case orders.above(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Up)
					state.Behaviour = MovingUp

				case orders.below(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Down)
					state.Behaviour = MovingDown

				case orders[state.Floor][elevio.BT_HallDown]:
					completedOrderC <- elevio.ButtonEvent{Floor: state.Floor, Button: elevio.BT_HallDown}
					doorOpenC <- true
					state.Behaviour = DoorOpenDown

				default:
					state.Behaviour = Idle
				}

			case DoorOpenDown:
				switch {
				case orders.below(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Down)
					state.Behaviour = MovingDown

				case orders.above(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Up)
					state.Behaviour = MovingUp

				case orders[state.Floor][elevio.BT_HallUp]:
					completedOrderC <- elevio.ButtonEvent{Floor: state.Floor, Button: elevio.BT_HallUp}
					doorOpenC <- true
					state.Behaviour = DoorOpenUp

				default:
					state.Behaviour = Idle
				}
			default:
				panic(state)
			}
			stateC <- state

		case state.Floor = <-floorEnteredC:
			switch state.Behaviour {
			case MovingUp:
				switch {
				case orders[state.Floor][elevio.BT_HallUp] || orders[state.Floor][elevio.BT_Cab]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					state.Behaviour = DoorOpenUp
					// Clear orders

				case orders.above(state.Floor):

				case orders[state.Floor][elevio.BT_HallDown]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					state.Behaviour = DoorOpenDown
					// Clear orders
				
				case orders.below(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Down)
					state.Behaviour = MovingDown

				default:
					elevio.SetMotorDirection(elevio.MD_Stop)
					state.Behaviour = Idle
				}

			case MovingDown:
				switch {
				case orders[state.Floor][elevio.BT_HallDown] || orders[state.Floor][elevio.BT_Cab]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					state.Behaviour = DoorOpenDown
					// Clear orders

				case orders.below(state.Floor):

				case orders[state.Floor][elevio.BT_HallUp]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					state.Behaviour = DoorOpenUp
					// Clear orders

				case orders.above(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Up)
					state.Behaviour = MovingUp

				default:
					elevio.SetMotorDirection(elevio.MD_Stop)
					state.Behaviour = Idle
				}

			default:
				panic(state)
			}
			stateC <- state

		case orders = <-ordersC:
			switch state.Behaviour {
			case DoorOpenUp:
				//clear orders possibly

			case DoorOpenDown:
				//clear orders possibly

			case Idle:
				switch {
				case orders.here(state.Floor):
					doorOpenC <- true
					state.Behaviour = DoorOpenDown
					// Clear orders

				case orders.above(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Up)
					state.Behaviour = MovingUp

				case orders.below(state.Floor):
					elevio.SetMotorDirection(elevio.MD_Up)
					state.Behaviour = MovingUp
				}
			}
		}
	}
}
