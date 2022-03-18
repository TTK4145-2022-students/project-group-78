package door

import (
	"log"
	"time"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/config"
)

type state int

const (
	closed state = iota
	inCountdown
	stuck
)

func Door(openC <-chan bool, closedC chan<- bool) {
	obstructionC := make(chan bool)
	go elevio.PollObstructionSwitch(obstructionC)

	obstructed := false
	state := closed
	var timer *time.Timer

	for {
		select {
		case obstructed = <-obstructionC:
			if state == stuck && !obstructed {
				elevio.SetDoorOpenLamp(false)
				closedC <- true
				state = closed
			}

		case <-openC:
			switch state {
			case closed:
				elevio.SetDoorOpenLamp(true)
				timer = time.NewTimer(config.DoorOpenTime)
				state = inCountdown

			case inCountdown:
				timer = time.NewTimer(config.DoorOpenTime)

			case stuck:
				timer = time.NewTimer(config.DoorOpenTime)
				state = inCountdown

			default:
				log.Panicf("door: unknown state %v", state)
			}

		case <-timer.C:
			if state != inCountdown {
				log.Panicf("door: timer expired while in %v", state)
			}
			if obstructed {
				state = stuck
			} else {
				elevio.SetDoorOpenLamp(false)
				closedC <- true
				state = closed
			}
		}
	}
}
