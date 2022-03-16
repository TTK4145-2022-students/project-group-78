package central

import (
	"time"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
)

type CentralState struct {
	Origin     int
	HallOrders [config.NUM_FLOORS][2]Order
	CabOrders  [config.NUM_ELEVS][config.NUM_FLOORS]bool
	States     [config.NUM_ELEVS]elevator.State
}

type Order struct {
	Active bool
	Time   time.Time
}

func New(origin int, state elevator.State) (cs CentralState) {
	cs.Origin = origin
	cs.States[origin] = state
	return cs
}

func (cs CentralState) SetOrder(be elevio.ButtonEvent, value bool) CentralState {
	if be.Button == elevio.BT_Cab {
		cs.CabOrders[cs.Origin][be.Floor] = value
	} else {
		o := Order{value, time.Now()}
		cs.HallOrders[be.Floor][be.Button] = o
	}
	return cs
}

// Merge newCs onto cs
func (cs CentralState) Merge(newCs CentralState) CentralState {
	cs.States[newCs.Origin] = newCs.States[newCs.Origin]
	cs.CabOrders[newCs.Origin] = newCs.CabOrders[newCs.Origin]
	for f := 0; f < config.NUM_FLOORS; f++ {
		for bt := 0; bt < len(cs.HallOrders[f]); bt++ {
			if cs.HallOrders[f][bt].Time.Before(newCs.HallOrders[f][bt].Time) {
				cs.HallOrders[f][bt] = newCs.HallOrders[f][bt]
			}
		}
	}
	return cs
}
