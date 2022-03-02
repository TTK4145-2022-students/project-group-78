package elevator

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type relativePosition int

const (
	Above relativePositon = iota
	Neutral
	Below
)

type State int

const (
	DoorOpen State = iota
	Moving
	Idle
)

var currentFloor int = -1
var target int = -1
var relativePosition relativePosition = Below
var motorDirection MotorDirection = MD_Up
var state = Moving
var newTarget bool = false

var es = ElevatorState{State: Moving, Direction: elevio.MD_Up, Floor: -1}

func targetReached(doorTimer *time.Timer) {
	SetMotorDirection(MD_Stop)
	SetDoorOpenLamp(true)
	doorTimer = time.NewTimer(time.Second * 3)
}

func startMotorTowardsTarget(doorTimer *time.Timer) {
	if currentFloor == target {
		switch relativePosition {
		case Above:
			setMotorDirection(MD_Down)

		case Neutral:
			State = DoorOpen
			targetReached(&doorTimer)

		case Below:
			SetMotorDirection(MD_Up)
		}

	} else if currentFloor < target {
		SetMotorDirection(MD_Up)

	} else {
		SetMotorDirection(MD_Down)

	}
	newTarget = false
}

func floorEntered(f int, doorTimer *time.Timer) State {
	currentFloor = f
	if currentFloor == target {
		State = DoorOpen
		targetReached(&doorTimer)
	}
	switch motorDirection {
	case MD_Up:
		relativePosition = Above
	case MD_Down:
		relativePosition = Below
	case MD_Stop:
		relativePosition = Neutral
	default:
	}
}

func doorTimedOut() {
	if newTarget { // TODO: Calculate new target instead
		State = Moving
		startMotorTowardsTarget(&doorTimer)
	} else {
		State = Idle
	}
	SetDoorOpenLamp(false)
}

func targetOrderUpdated(o order.Order) {
	newTarget = true
	target = o
	if State != DoorOpen {
		if target == currentFloor && relativePosition == Neutral {
			State = DoorOpen
			targetReached()
		} else {
			startMotorTowardsTarget()
		}
	}
}
