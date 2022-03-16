package central

import (
	"time"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
)

type CentralState struct {
	Origin      int
	HallOrders  [config.NUM_FLOORS][2]Order
	CabOrders   [config.NUM_ELEVS][config.NUM_FLOORS]bool
	States      [config.NUM_ELEVS]elevator.State
	LastUpdated [config.NUM_ELEVS]time.Time
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

func (cs CentralState) AddOrder(be elevio.ButtonEvent) CentralState {
	if be.Button == elevio.BT_Cab {
		cs.CabOrders[cs.Origin][be.Floor] = true
	} else if !cs.HallOrders[be.Floor][be.Button].Active {
		cs.HallOrders[be.Floor][be.Button] = Order{true, time.Now()}
	}
	return cs
}

func (cs CentralState) RemoveOrder(be elevio.ButtonEvent) CentralState {
	if be.Button == elevio.BT_Cab {
		cs.CabOrders[cs.Origin][be.Floor] = false
	} else {
		cs.HallOrders[be.Floor][be.Button] = Order{false, time.Now()}
	}
	return cs
}

// Merge newCs onto cs
func (cs CentralState) Merge(newCs CentralState) CentralState {
	cs.States[newCs.Origin] = newCs.States[newCs.Origin]
	cs.CabOrders[newCs.Origin] = newCs.CabOrders[newCs.Origin]
	cs.LastUpdated[newCs.Origin] = newCs.LastUpdated[newCs.Origin]
	for f := range cs.HallOrders {
		for bt := range cs.HallOrders[f] {
			if cs.HallOrders[f][bt].Time.Before(newCs.HallOrders[f][bt].Time) {
				cs.HallOrders[f][bt] = newCs.HallOrders[f][bt]
			}
		}
	}
	return cs
}
