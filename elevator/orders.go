package elevator

import (
	"fmt"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/config"
)

type Orders [config.NumFloors][config.NumOrderTypes]bool

func (o Orders) above(floor int) bool {
	for f := floor + 1; f < len(o); f++ {
		for bt := range o[f] {
			if o[f][bt] {
				return true
			}
		}
	}
	return false
}

func (o Orders) below(floor int) bool {
	for f := 0; f < floor; f++ {
		for bt := range o[f] {
			if o[f][bt] {
				return true
			}
		}
	}
	return false
}

func (o Orders) here(f int) bool {
	for bt := range o[f] {
		if o[f][bt] {
			return true
		}
	}
	return false
}

func nextAction(orders Orders, floor int, d elevio.MotorDirection) (elevio.MotorDirection, Behaviour) {
	switch d {
	case elevio.MD_Up:
		switch {
		case orders.above(floor):
			return elevio.MD_Up, Moving

		case orders.here(floor):
			return elevio.MD_Down, DoorOpen

		case orders.below(floor):
			return elevio.MD_Down, Moving
		}

	case elevio.MD_Down:
		switch {
		case orders.below(floor):
			return elevio.MD_Down, Moving

		case orders.here(floor):
			return elevio.MD_Up, DoorOpen

		case orders.above(floor):
			return elevio.MD_Up, Moving
		}

	case elevio.MD_Stop:
		switch {
		case orders.here(floor):
			return elevio.MD_Stop, DoorOpen

		case orders.above(floor):
			return elevio.MD_Up, Moving

		case orders.below(floor):
			return elevio.MD_Down, Moving
		}
	}
	return elevio.MD_Stop, Idle
}

func shouldStop(orders Orders, floor int, direction elevio.MotorDirection) bool {
	switch direction {
	case elevio.MD_Down:
		return orders[floor][elevio.BT_HallDown] || orders[floor][elevio.BT_Cab] || !orders.below(floor)
	case elevio.MD_Up:
		return orders[floor][elevio.BT_HallUp] || orders[floor][elevio.BT_Cab] || !orders.above(floor)
	case elevio.MD_Stop:
		panic("elevator: direction should not be stop")
	default:
		panic(fmt.Sprintf("elevator: unknown direction %v", direction))
	}
}

func clearOrders(orders Orders, floor int, direction elevio.MotorDirection, completedOrderC chan<- elevio.ButtonEvent) {
	if orders[floor][elevio.BT_Cab] {
		completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}

	switch direction {
	case elevio.MD_Up:
		if !orders.above(floor) && !orders[floor][elevio.BT_HallUp] && orders[floor][elevio.BT_HallDown] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallDown}
		} else if orders[floor][elevio.BT_HallUp] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallUp}
		}

	case elevio.MD_Down:
		if !orders.above(floor) && !orders[floor][elevio.BT_HallDown] && orders[floor][elevio.BT_HallUp] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallUp}
		} else if orders[floor][elevio.BT_HallDown] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallDown}
		}

	case elevio.MD_Stop:
		if orders[floor][elevio.BT_HallUp] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallUp}
		} else if orders[floor][elevio.BT_HallDown] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallDown}
		}

	default:
		panic(fmt.Sprintf("elevator: unknown direction %v", direction))
	}
}

func shouldClear(orders Orders, f int, direction elevio.MotorDirection) bool {
	return orders[f][elevio.BT_Cab] ||
		direction == elevio.MD_Up && orders[f][elevio.BT_HallUp] ||
		direction == elevio.MD_Down && orders[f][elevio.BT_HallDown]
}
