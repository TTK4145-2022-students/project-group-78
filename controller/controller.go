package controller

import "github.com/TTK4145-2022-students/project-group-78/elevio"

var TargetC chan int

type relativePosition int

const (
	above   relativePosition = 1
	below                    = -1
	atFloor                  = 0
)

func Controller(floorC chan int, targetReachedC chan int, directionC chan elevio.MotorDirection) {
	TargetC = make(chan int)
	floorEnteredC := make(chan int)
	go elevio.PollFloorSensor(floorEnteredC)

	var relPos relativePosition = below
	floor := -1
	target := -1
	var direction relativePosition = below

	for {
		select {
		case floor = <-floorEnteredC:
			floorC <- floor
			if floor == target {
				elevio.SetMotorDirection(elevio.MD_Stop)
				directionC <- elevio.MD_Stop
				targetReachedC <- target
				relPos = atFloor
			} else {
				relPos = relativePosition(direction)
			}

		case target = <-TargetC:
			if floor == target && relPos == atFloor {
				targetReachedC <- target
			} else {
				//TODO: Ulrik fix?
				/*
					Suggestion:
						Assume that you are not at target
						Create a pure function that calculates the direction
						Make sure to emit direction at directionC if it changes
						Also update relPos if needed
				*/
			}
		}
	}
}

/*
func startMotorTowardsTarget(es ElevatorState) ElevatorState {
	if es.Target.Floor == es.Floor {
		switch es.RelativePosition {
		case Above:
			elevio.SetMotorDirection(elevio.MD_Down)
			es.Direction = elevio.MD_Down

		case AtFloor:
			log.Panic("Tried to start motor while at target")

		case Below:
			elevio.SetMotorDirection(elevio.MD_Up)
			es.Direction = elevio.MD_Up
		}

	} else if es.Floor < es.Target.Floor {
		elevio.SetMotorDirection(elevio.MD_Up)
		es.Direction = elevio.MD_Up

	} else {
		elevio.SetMotorDirection(elevio.MD_Down)
		es.Direction = elevio.MD_Down

	}
	es.State = Moving
	return es
}
*/
