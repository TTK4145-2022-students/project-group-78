package central

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/vetleras/driver-go/elevio"
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
		CabOrders: make(map[string][config.NUM_FLOORS]bool),
		States:    make(map[string]elevator.State),
	}
	cs.States[origin] = state
	return cs
}

func (cs CentralState) deepCopy() CentralState {
	copiedCabOrders := make(map[string][config.NUM_FLOORS]bool)
	for id, values := range cs.CabOrders {
		copiedCabOrders[id] = values
	}
	cs.CabOrders = copiedCabOrders

	copiedStates := make(map[string]elevator.State)
	for id, state := range cs.States {
		copiedStates[id] = state
	}
	cs.States = copiedStates

	return cs
}

func (cs CentralState) SetOrder(be elevio.ButtonEvent, value bool) CentralState {
	cs = cs.deepCopy()
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
	cs = cs.deepCopy()
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
