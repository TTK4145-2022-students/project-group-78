package elevator

import (
	"github.com/TTK4145-2022-students/project-group-78/order"
)
type relativePosition int
const (
	Above relativePositon = iota
	Neutral
	Below
)
type state int
const (
	DoorOpen state= iota
	Moving
	Idle
)
var currentFloor int = -1
var target int = -1
var relativePosition relativePosition = Below
var motorDirection MotorDirection = MD_Up
var state = Moving
var newTarget bool = false

func targetReached(doorTimer *time.Timer){
	SetMotorDirection(MD_Stop)
	SetDoorOpenLamp(true)
	doorTimer = time.NewTimer(time.Second*3)
}

func startMotorTowardsTarget(doorTimer *time.Timer){
	if (currentFloor == target) {
		switch relativePosition{
			case Above:
				setMotorDirection(MD_Down);

			case Neutral:
				state = DoorOpen
				targetReached(&doorTimer)

			case Below:
				SetMotorDirection(MD_Up);
		}

	}else if (currentFloor < target) {
		SetMotorDirection(MD_Up);

	}else {
		SetMotorDirection(MD_Down);

	}
	newTarget = false
}

func floorEntered(f int,doorTimer *time.Timer) { 
	currentFloor = f
	if currentFloor == target{
		state = DoorOpen
		targetReached(&doorTimer)
	}
	switch motorDirection{
		case MD_Up:
			relativePosition = Above
		case MD_Down:
			relativePosition = Below
		case MD_Stop:
			relativePosition = Neutral
		default:
	}
}

func closeDoor() {
	if (newTarget){
		state = Moving 
		startMotorTowardsTarget(&doorTimer)
	}else{
		state = Idle
	}
	SetDoorOpenLamp(false)
}

func targetOrderUpdated(o order.Order) {
	newTarget = true
	target = o
	if state != DoorOpen{
		if target == currentFloor && relativePosition == Neutral{
			state = DoorOpen
			targetReached()
		}else{
			startMotorTowardsTarget()
		}
	}
}

