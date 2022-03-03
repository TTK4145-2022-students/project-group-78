package elevator

import (
	"log"

	"github.com/TTK4145-2022-students/project-group-78/controller"
	"github.com/TTK4145-2022-students/project-group-78/door"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

var TargetC chan int

type Behaviour string

const (
	DoorOpen                  Behaviour = "doorOpen"
	DoorOpenWithPendingTarget           = "doorOpenWithPendingTarget"
	ServingOrder                        = "servingOrder"
	Idle                                = "idle"
)

type State struct {
	Behaviour Behaviour
	Direction elevio.MotorDirection
	Floor     int
}

func Elevator(targetReachedCOut chan int, stateC chan State) {
	doorClosedC := make(chan bool)
	floorEnteredC, targetReachedCIn := make(chan int), make(chan int)
	directionSetC := make(chan elevio.MotorDirection)

	go door.Door(doorClosedC)
	go controller.Controller(floorEnteredC, targetReachedCIn, directionSetC)

	var state State
	var target int
	for {
		select {
		case <-doorClosedC:
			switch state.Behaviour {
			case DoorOpen:
				state.Behaviour = Idle
				stateC <- state

			case DoorOpenWithPendingTarget:
				controller.TargetC <- target
				state.Behaviour = ServingOrder
				stateC <- state

			default:
				log.Panicf("door closed while in %v", state.Behaviour)
			}

		case state.Floor = <-floorEnteredC:
			stateC <- state

		case t := <-targetReachedCIn:
			switch state.Behaviour {
			case ServingOrder:
				door.OpenC <- true
				targetReachedCOut <- t
				state.Behaviour = DoorOpen
				stateC <- state

			default:
				log.Panicf("target reached while in %v", state.Behaviour)
			}

		case state.Direction = <-directionSetC:
			stateC <- state

		case t := <-TargetC:
			switch state.Behaviour {
			case Idle:
				controller.TargetC <- t
				state.Behaviour = ServingOrder
				stateC <- state

			case DoorOpen:
				target = t
				state.Behaviour = DoorOpenWithPendingTarget
				stateC <- state

			case DoorOpenWithPendingTarget:
				target = t

			case ServingOrder:
				controller.TargetC <- t
			}
		}
	}

}
