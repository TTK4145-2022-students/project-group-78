package orders

import (
	"testing"

	"github.com/TTK4145-2022-students/project-group-78/elevator"
)

func TestCalculateTarget(t *testing.T) {
	cs := CentralState{}
	cs.HallOrders[0].Up.Active = true
	cs.States[0].Behaviour = elevator.Idle
	cs.States[1].Behaviour = elevator.Idle
	cs.States[2].Behaviour = elevator.Idle
	hallRequestAssigner(cs)
	t.Fail()
}
