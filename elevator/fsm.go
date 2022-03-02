package elevator

import (
	"log"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

func buttonPressed(be elevio.ButtonEvent, es ElevatorState) CentralState {
	if be.Button == elevio.BT_Cab {
		es.CabOrders[be.Floor] = true
		return NewCentralState(id, es)
	} else {
		cs := NewCentralState(id, es)
		o := Order{
			Type:  be.Button,
			Floor: be.Floor,
		}
		cs.HallOrders[o] = time.Now()
		return cs
	}
}

func floorEntered(f int, es ElevatorState, doorTimer *time.Timer) CentralState {
	if f == es.Target.Floor {
		return targetReached(es, doorTimer)
	} else {
		es.RelativePosition = RelativePosition(es.Direction)
		return NewCentralState(id, es)
	}
}

func closeDoor(es ElevatorState) ElevatorState {
	es.State = Idle
	elevio.SetDoorOpenLamp(false)
	return es
}

func targetOrderUpdated(es ElevatorState, doorTimer *time.Timer) CentralState {
	switch es.State {
	case Idle:
		if es.Target.Floor == es.Floor {
			return targetReached(es, doorTimer)
		} else {
			return NewCentralState(id, startMotorTowardsTarget(es))
		}

	case Moving:
		return NewCentralState(id, startMotorTowardsTarget(es))

	case DoorOpen:
		if es.Target.Floor == es.Floor {
			return targetReached(es, doorTimer)
		} else {
			return NewCentralState(id, es)
		}

	default:
		log.Panicf("Invalid state %v", es.State)
		return CentralState{}
	}
}

func targetReached(es ElevatorState, doorTimer *time.Timer) CentralState {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	doorTimer = time.NewTimer(config.DOOR_OPEN_TIME)
	es.Direction = elevio.MD_Stop
	es.State = DoorOpen

	if es.Target.Type == elevio.BT_Cab {
		es.CabOrders[es.Target.Floor] = false
		return NewCentralState(id, es)
	} else {
		cs := NewCentralState(id, es)
		cs.ServedHallOrders[es.Target] = time.Now()
		return cs
	}
}

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
