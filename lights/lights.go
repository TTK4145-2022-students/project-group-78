package lights

import (
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/orders"
)

var SetC chan orders.CentralState

func Lights() {
	prevCs := orders.CentralState{}
	for {
		select {
		case cs := <-SetC:
			for f := 0; f < len(cs.CabOrders); f++ {
				if cs.CabOrders[cs.Origin][f] != prevCs.CabOrders[cs.Origin][f] {
					elevio.SetButtonLamp(elevio.BT_Cab, f, cs.CabOrders[cs.Origin][f])
				}
			}

			for f := 0; f < len(cs.HallOrders); f++ {
				if cs.HallOrders[f].Up.Active != prevCs.HallOrders[f].Up.Active {
					elevio.SetButtonLamp(elevio.BT_HallUp, f, cs.HallOrders[f].Up.Active)
				}
				if cs.HallOrders[f].Down.Active != prevCs.HallOrders[f].Down.Active {
					elevio.SetButtonLamp(elevio.BT_HallDown, f, cs.HallOrders[f].Down.Active)
				}
			}
			prevCs = cs
		}
	}
}
