package elevator

import "github.com/TTK4145-2022-students/project-group-78/elevio"

type Elevator struct {
	button      chan elevio.ButtonEvent
	floor       chan int
	obstruction chan bool
}

func New() *Elevator {
	elevio.Init()

	go e.run()

	return &Elevator{}
}

func (e *Elevator) run() {
	for {
		select {
		case button := <-e.button:
			e.handleButton(button)

		case floor := <-e.floor:
			e.handleFloorSensor(floor)
		}
	}
}

