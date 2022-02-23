package distributor

import (
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/events"
	"github.com/TTK4145-2022-students/project-group-78/mocknet"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.SetLevel(logrus.DebugLevel)
}

func TestDistributor(t *testing.T) {
	distributor := New(1)
	defer distributor.Stop()

	mocknet := mocknet.New(config.PORT)
	defer mocknet.Close()

	cs := central.NewCentralState()
	id := 1
	event := events.DoorOpened{}
	time_ := time.Now()
	cs[id][event] = time_
	distributor.Send(cs)

	ncs := <-distributor.StateUpdate
	assert.True(t, ncs[id][event].Equal(time_))
}
