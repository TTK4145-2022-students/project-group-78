package door

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

var OpenC chan bool

func Door(closedC chan bool) {
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
				closedC <- true
			}

		case <-OpenC:
			elevio.SetDoorOpenLamp(true)
			timer = time.NewTimer(config.DOOR_OPEN_TIME)

		case <-timer.C:
			if !obstructed && doorOpen {
				elevio.SetDoorOpenLamp(false)
				closedC <- true
			}
		}
	}
}
