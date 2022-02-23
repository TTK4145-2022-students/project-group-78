package controller

import (
	"time"
)
type relative_position int

const (
	Above relative_position = iota
	Neutral  
	Below 
)
type state int

const (
	Moving state = iota
	DoorOpen  
	AtRest 
)
type LocalElevator struct{
	target int;
	current_floor int;
	relative_position relative_position;
	motor_dir MotorDirection;
	state State;
}



func target_reached(){
	SetMotorDirection(MD_Stop);
	go door_open()
}
func start_motor_towards_target(){
	if (elev.state != DoorOpen){
		if (elev.current_floor == target) {
			switch elev.relative_position{
				case Above:
					setMotorDirection(MD_Down);

				case Neutral:
					target_reached()

				case Below:
					SetMotorDirection(MD_Up);
			}
		}else if (elev.current_floor < elev.target) {
			SetMotorDirection(MD_Up);
		}else {
			SetMotorDirection(MD_Down);
		}
		elev.new_target = false
	}
}


func door_open(){
	elev.state = DoorOpen
	SetDoorOpenLamp(true)
	timer := time.NewTimer(3*time.Second)
	<-timer.C
	if elev.new_target{
		elev.state = AtRest
		start_motor_towards_target(target)
		elev.state = Moving
	}
	else{
		elev.state = AtRest
	}
	
	SetDoorOpenLamp(false)
	
}
func Run_elevator(){
	elev = LocalElevator{-1,-1, Neutral, MD_Stop, AtRest};
	addr = "11111"
	numFloors = 5

	Init(addr, numFloors);

	for {
		event = EM_listen_for_event()
		if (event){
			switch event.type{
				case floorSensorTriggered:
					elev.current_floor = event.floor
					if current_floor == target{
						target_reached()
					}
					switch elev.motor_dir{
						case MD_Up:
							elev.relative_position = Above
						case MD_Down:
							elev.relative_position = Below
						case MD_Stop:
							elev.relative_position = Neutral
						default:
					}
					

				case newTarget:
					elev.new_target = true
					elev.target = event.target
					start_motor_towards_target(elev.target)

				default:
			}
		}
	}
}
