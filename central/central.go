package central

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
)

type Event interface{}

type CentralState map[int]ElevatorState

type ElevatorState map[Event]time.Time

func NewCentralState() CentralState {
	c := make(CentralState, config.NUM_ELEVATORS)
	for id := 1; id <= config.NUM_ELEVATORS; id++ {
		c[id] = make(ElevatorState)
	}
	return c
}

// Merge s into c
func (c CentralState) Merge(s CentralState) {
	for id, es := range s {
		for event, time := range es {
			if time.After(c[id][event]) {
				c[id][event] = time
			}
		}
	}
}
