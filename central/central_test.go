package central

import (
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/stretchr/testify/assert"
)

func TestCentral(t *testing.T) {
	c := New()

	time_ := time.Time{}
	event := elevator.DoorOpened{}
	in := NetworkState{
		1: ElevatorState{
			event: time_,
		}}
	c.StateIn <- in

	t.Run("Initial event", func(t *testing.T) {
		out := <-c.StateOut
		assert.Equal(t, out[1][event], time_)
	})

	t.Run("Updating event", func(t *testing.T) {
		time_ := time.Now()
		in[1][event] = time_
		c.StateIn <- in

		out := <-c.StateOut
		assert.Equal(t, out[1][event], time_)
	})
}
