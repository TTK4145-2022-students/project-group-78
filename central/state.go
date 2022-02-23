package central

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
)

type Event interface{}

type NetworkState map[int]ElevatorState

type ElevatorState map[Event]time.Time

func MakeNetworkState() NetworkState {
	ns := make(NetworkState, config.NUM_ELEVATORS)
	for id := 1; id <= config.NUM_ELEVATORS; id++ {
		ns[id] = make(ElevatorState)
	}
	return ns
}
