package lights

import (
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

var lights elevator.Orders

func Set(cs central.CentralState) {
	cabOrders := cs.CabOrders[cs.Origin]
	for f := 0; f < len(cabOrders); f++ {
		if cabOrders[f] != lights[f][elevio.BT_Cab] {
			elevio.SetButtonLamp(elevio.BT_Cab, f, cabOrders[f])
			lights[f][elevio.BT_Cab] = cabOrders[f]
		}
	}
	for f := 0; f < len(cs.HallOrders); f++ {
		for bt := 0; bt < 2; bt++ {
			value := cs.HallOrders[f][bt].Active
			if value != lights[f][bt] {
				elevio.SetButtonLamp(elevio.ButtonType(bt), f, value)
				lights[f][bt] = value
			}
		}
	}
}

func Clear() {
	elevio.SetDoorOpenLamp(false)
	for f := 0; f < config.NUM_FLOORS; f++ {
		for bt := 0; bt < 3; bt++ {
			elevio.SetButtonLamp(elevio.ButtonType(bt), f, false)
		}
	}
}
