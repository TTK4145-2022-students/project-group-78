package central

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockEvent struct{}

func TestCentralState(t *testing.T) {
	cs := NewCentralState()

	id := 1
	event := MockEvent{}

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
