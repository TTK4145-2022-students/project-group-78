package lights

import (
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

var SetC chan central.CentralState

func Lights() {
	prevCs := central.CentralState{}
	for {
		select {
		case cs := <-SetC:
			for i := 0; i < len(cs.Elevators[cs.Origin].CabOrders); i++ {
				if cs.Elevators[cs.Origin].CabOrders[i] != prevCs.Elevators[cs.Origin].CabOrders[i] {
					elevio.SetButtonLamp(elevio.BT_Cab, i, cs.Elevators[cs.Origin].CabOrders[i])
				}
			}

			for i := 0; i < len(cs.HallOrders); i++ {
				if cs.HallOrders[i].Up.Active != prevCs.HallOrders[i].Up.Active {
					elevio.SetButtonLamp(elevio.BT_HallUp, i, cs.HallOrders[i].Up.Active)
				}
				if cs.HallOrders[i].Down.Active != prevCs.HallOrders[i].Down.Active {
					elevio.SetButtonLamp(elevio.BT_HallDown, i, cs.HallOrders[i].Down.Active)
				}
			}
			prevCs = cs
		}
	}
}
