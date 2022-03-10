package elevator

import (
	"fmt"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/vetleras/driver-go/elevio"
)

type Orders [config.NUM_FLOORS][3]bool

func (o Orders) Above(floor int) bool {
	for f := floor + 1; f < len(o); f++ {
		for bt := 0; bt < 3; bt++ {
			if o[f][bt] {
				return true
			}
		}
	}
	return false
}

func (o Orders) Below(floor int) bool {
	for f := 0; f < floor; f++ {
		for bt := 0; bt < 3; bt++ {
			if o[f][bt] {
				return true
			}
		}
	}
	return false
}

func (o Orders) Here(f int) bool {
	for bt := 0; bt < 3; bt++ {
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
		case orders.Above(floor):
			return elevio.MD_Up, Moving

		case orders.Here(floor):
			return elevio.MD_Down, DoorOpen

		case orders.Below(floor):
			return elevio.MD_Down, Moving
		}

	case elevio.MD_Down:
		switch {
		case orders.Below(floor):
			return elevio.MD_Down, Moving

		case orders.Here(floor):
			return elevio.MD_Up, DoorOpen

		case orders.Above(floor):
			return elevio.MD_Up, Moving
		}

	case elevio.MD_Stop:
		switch {
		case orders.Here(floor):
			return elevio.MD_Stop, DoorOpen

		case orders.Above(floor):
			return elevio.MD_Up, Moving

		case orders.Below(floor):
			return elevio.MD_Down, Moving
		}
	}
	return elevio.MD_Stop, Idle
}

func shouldStop(orders Orders, floor int, direction elevio.MotorDirection) bool {
	switch direction {
	case elevio.MD_Down:
		return orders[floor][elevio.BT_HallDown] || orders[floor][elevio.BT_Cab] || !orders.Below(floor)
	case elevio.MD_Up:
		return orders[floor][elevio.BT_HallUp] || orders[floor][elevio.BT_Cab] || !orders.Above(floor)
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
		if !orders.Above(floor) && !orders[floor][elevio.BT_HallUp] && orders[floor][elevio.BT_HallDown] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallDown}
		} else if orders[floor][elevio.BT_HallUp] {
			completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_HallUp}
		}

	case elevio.MD_Down:
		if !orders.Above(floor) && !orders[floor][elevio.BT_HallDown] && orders[floor][elevio.BT_HallUp] {
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
