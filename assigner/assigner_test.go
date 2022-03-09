package assigner

import (
	"testing"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	cs := central.New("1", elevator.State{Behaviour: elevator.Idle})
	cs.HallOrders[0][0].Active = true
	assert.True(t, Assigner(cs)[0][0])
}
