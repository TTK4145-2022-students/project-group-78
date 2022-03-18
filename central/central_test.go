package central

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	cs := CentralState{Origin: 1}
	newCs := CentralState{Origin: 2}
	o := Order{true, time.Now()}
	newCs.HallOrders[0][0] = o
	cs = cs.Merge(newCs)
	assert.True(t, cs.HallOrders[0][0].Active)

	newCs.CabOrders[1] = [4]bool{true, false, false, false}
	newCs.CabOrders[2] = [4]bool{false, true, false, false}
	cs = cs.Merge(newCs)
	assert.False(t, cs.CabOrders[1][0])
	assert.True(t, cs.CabOrders[2][1])
}
