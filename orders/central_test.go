package orders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	cs := CentralState{}
	newCs := CentralState{}
	newCs.HallOrders[0].Up.Time = time.Now()
	newCs.HallOrders[0].Up.Active = true
	cs = cs.Merge(newCs)
	assert.True(t, cs.HallOrders[0].Up.Active)

	newCs.Origin = 1
	newCs.CabOrders[0][0] = true
	newCs.CabOrders[1][0] = true
	cs = cs.Merge(newCs)
	assert.False(t, cs.CabOrders[0][0])
	assert.True(t, cs.CabOrders[1][0])
}
