package elevator

import (
	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/config"
)

type Orders [config.NumFloors][config.NumOrderTypes]bool

func (o Orders) inDirection(floor int, direction Direction) bool {
	switch direction {
	case Up:
		for f := floor + 1; f < len(o); f++ {
			for bt := range o[f] {
				if o[f][bt] {
					return true
				}
			}
		}
		return false

	case Down:
		for f := 0; f < floor; f++ {
			for bt := range o[f] {
				if o[f][bt] {
					return true
				}
			}
		}
		return false

	default:
		panic(direction)
	}
}

func clearOrders(orders Orders, floor int, direction Direction, completedOrderC chan<- elevio.ButtonEvent) {
	if orders[floor][elevio.BT_Cab] {
		completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if orders[floor][direction] {
		completedOrderC <- elevio.ButtonEvent{Floor: floor, Button: direction.toBt()}
	}
}
