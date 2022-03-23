package assigner

import (
	"fmt"
	"testing"
	"time"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	cs := central.CentralState{Origin: 2}

	cs.States[0] = elevator.State{
		Behaviour: elevator.DoorOpen,
		Direction: elevio.MD_Down,
		Floor:     0,
	}

	cs.States[1] = cs.States[0]

	cs.States[2] = elevator.State{
		Behaviour: elevator.DoorOpen,
		Direction: elevio.MD_Up,
		Floor:     1,
	}

	cs = cs.AddOrder(elevio.ButtonEvent{Floor: 1, Button: elevio.BT_HallUp}).AddOrder(elevio.ButtonEvent{Floor: 1, Button: elevio.BT_HallUp})
	fmt.Print(Assigner(cs))

	//cs.LastUpdated[0] = time.Now()
	cs.LastUpdated[2] = time.Now()

	fmt.Print(Assigner(cs))
	//cs.Origin = 1


	assert.True(t, false)
}

func BenchmarkAddigner(b *testing.B) {
	Assigner(central.CentralState{})
}
