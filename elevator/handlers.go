package elevator

import "github.com/TTK4145-2022-students/project-group-78/elevio"

func (e *Elevator) handleButtonPress(button elevio.ButtonEvent) {
	switch e.state {
	case idle:
	case moving:

	}
}

func (e *Elevator) handleFloorSensor(floor int) {
	elev.current_floor = event.floor
	if current_floor == target {
		target_reached()
	}
	switch elev.motor_dir {
	case MD_Up:
		elev.relative_position = Above
	case MD_Down:
		elev.relative_position = Below
	case MD_Stop:
		elev.relative_position = Neutral
	default:
	}
}
