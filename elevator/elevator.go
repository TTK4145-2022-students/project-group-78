package elevator

import (
	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/door"
)

type State struct {
	Behaviour Behaviour
	Floor     int
	Direction Direction
}

type Behaviour int

const (
	Idle Behaviour = iota
	DoorOpen
	Moving
)

func Elevator(ordersC <-chan Orders, completedOrderC chan<- elevio.ButtonEvent, stateC chan<- State) {
	doorOpenC := make(chan bool, config.ChanSize)
	doorClosedC := make(chan bool, config.ChanSize)
	floorEnteredC := make(chan int)

	go door.Door(doorOpenC, doorClosedC)
	go elevio.PollFloorSensor(floorEnteredC)

	elevio.SetMotorDirection(elevio.MD_Down)
	state := State{Behaviour: Moving, Direction: Down}
	var orders Orders

	for {
		select {
		case <-doorClosedC:
			switch state.Behaviour {
			case DoorOpen:
				switch {
				case orders.inDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(state.Direction.toMd())
					state.Behaviour = Moving
					stateC <- state

				case orders[state.Floor][state.Direction.opposite()]:
					doorOpenC <- true
					state.Direction = state.Direction.opposite()
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					stateC <- state

				case orders.inDirection(state.Floor, state.Direction.opposite()):
					state.Direction = state.Direction.opposite()
					elevio.SetMotorDirection(state.Direction.toMd())
					state.Behaviour = Moving
					stateC <- state

				default:
					state.Behaviour = Idle
					stateC <- state
				}
			default:
				panic(state)
			}

		case state.Floor = <-floorEnteredC:
			elevio.SetFloorIndicator(state.Floor)
			switch state.Behaviour {
			case Moving:
				switch {
				case orders[state.Floor][state.Direction]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					state.Behaviour = DoorOpen

				case orders[state.Floor][elevio.BT_Cab] && orders.inDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					state.Behaviour = DoorOpen

				case orders[state.Floor][elevio.BT_Cab] && !orders[state.Floor][state.Direction.opposite()]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					state.Behaviour = DoorOpen

				case orders.inDirection(state.Floor, state.Direction):

				case orders[state.Floor][state.Direction.opposite()]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					state.Direction = state.Direction.opposite()
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					state.Behaviour = DoorOpen

				case orders.inDirection(state.Floor, state.Direction.opposite()):
					state.Direction = state.Direction.opposite()
					elevio.SetMotorDirection(state.Direction.toMd())

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
			case Idle:
				switch {
				case orders[state.Floor][state.Direction] || orders[state.Floor][elevio.BT_Cab]:
					doorOpenC <- true
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					state.Behaviour = DoorOpen
					stateC <- state

				case orders[state.Floor][state.Direction.opposite()]:
					doorOpenC <- true
					state.Direction = state.Direction.opposite()
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
					state.Behaviour = DoorOpen
					stateC <- state

				case orders.inDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(state.Direction.toMd())
					state.Behaviour = Moving
					stateC <- state

				case orders.inDirection(state.Floor, state.Direction.opposite()):
					state.Direction = state.Direction.opposite()
					elevio.SetMotorDirection(state.Direction.toMd())
					state.Behaviour = Moving
					stateC <- state

				default:
				}

			case DoorOpen:
				switch {
				case orders[state.Floor][elevio.BT_Cab] || orders[state.Floor][state.Direction]:
					doorOpenC <- true
					clearOrders(orders, state.Floor, state.Direction, completedOrderC)
				}

			case Moving:

			default:
				panic(state)
			}
		}
	}
}
