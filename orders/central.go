package orders

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
)

type CentralState struct {
	Origin     int
	HallOrders [config.NUM_FLOORS]struct {
		Up   Order
		Down Order
	}
	CabOrders [config.NUM_ELEVATORS][config.NUM_FLOORS]bool
	States    [config.NUM_ELEVATORS]elevator.State
}

type Order struct {
	Active bool
	Time   time.Time
}

// Merge newCs onto cs. If cs happend before newCs, overwrite with newCs
func (cs CentralState) Merge(newCs CentralState) CentralState {
	cs.States[newCs.Origin] = newCs.States[newCs.Origin]
	cs.CabOrders[newCs.Origin] = newCs.CabOrders[newCs.Origin]
	for i := 0; i < len(cs.HallOrders); i++ {
		if cs.HallOrders[i].Up.Time.Before(newCs.HallOrders[i].Up.Time) {
			cs.HallOrders[i].Up = newCs.HallOrders[i].Up
		}
		if cs.HallOrders[i].Down.Time.Before(newCs.HallOrders[i].Down.Time) {
			cs.HallOrders[i].Down = newCs.HallOrders[i].Down
		}
	}
	return cs
}
