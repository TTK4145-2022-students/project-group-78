package central

import (
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/events"
	"github.com/stretchr/testify/assert"
)

func TestCentralState(t *testing.T) {
	cs := NewCentralState()

	id := 1
	event := events.DoorOpened{}

	t.Run("Initial event", func(t *testing.T) {
		ncs := NewCentralState()
		time_ := time.Now()
		ncs[id][event] = time_
		cs.Merge(ncs)
		assert.Equal(t, cs[id][event], time_)
	})

	t.Run("Updating event", func(t *testing.T) {
		ncs := NewCentralState()
		time_ := time.Now()
		ncs[id][event] = time_
		cs.Merge(ncs)
		assert.Equal(t, cs[id][event], time_)
	})
}
