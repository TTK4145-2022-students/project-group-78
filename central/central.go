package central

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type CentralState struct {
	Origin     string
	HallOrders [config.NUM_FLOORS][2]Order
	CabOrders  map[string][config.NUM_FLOORS]bool
	States     map[string]elevator.State
}

type Order struct {
	Active bool
	Time   time.Time
}

func New(origin string, state elevator.State) CentralState {
	cs := CentralState{
		Origin:    origin,
		CabOrders: make(map[string][config.NUM_FLOORS]bool, config.NUM_ELEVATORS),
		States:    make(map[string]elevator.State, config.NUM_ELEVATORS),
	}
	cs.States[origin] = state
	return cs
}

func (cs CentralState) SetOrder(be elevio.ButtonEvent, value bool) CentralState {
	if be.Button == elevio.BT_Cab {
		cabOrders := cs.CabOrders[cs.Origin]
		cabOrders[be.Floor] = value
		cs.CabOrders[cs.Origin] = cabOrders
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
	for f := 0; f < len(cs.HallOrders); f++ {
		for bt := 0; bt < 2; bt++ {
			if cs.HallOrders[f][bt].Time.Before(newCs.HallOrders[f][bt].Time) {
				cs.HallOrders[f][bt] = newCs.HallOrders[f][bt]
			}
		}
	}
	return cs
}
