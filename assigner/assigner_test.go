package assigner

import (
	"testing"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	cs := central.CentralState{}
	cs.HallOrders[1][elevio.BT_HallUp].Active = true
	assert.True(t, Assigner(cs)[1][0])
}
