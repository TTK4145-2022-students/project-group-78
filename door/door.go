package door

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

func Door(openC <-chan bool, closedC chan<- bool) {
	obstructionC := make(chan bool)
	go elevio.PollObstructionSwitch(obstructionC)

	obstructed := false
	doorOpen := false

	timer := time.NewTimer(time.Hour)
	timer.Stop()

	for {
		select {
		case obstructed = <-obstructionC:
			if !obstructed && doorOpen {
				elevio.SetDoorOpenLamp(false)
				doorOpen = false
				closedC <- true
			}

		case <-openC:
			elevio.SetDoorOpenLamp(true)
			doorOpen = true
			timer = time.NewTimer(config.DOOR_OPEN_TIME)

		case <-timer.C:
			if !obstructed && doorOpen {
				elevio.SetDoorOpenLamp(false)
				doorOpen = false
				closedC <- true
			}
		}
	}
}
