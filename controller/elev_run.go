package controller

import "Driver-go/elevio"


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

elev = LocalElevator{-1,-1, Neutral, MD_Stop, AtRest};

func target_reached(){
	SetMotorDirection(MD_Stop);
	go door_open()
}

func Run_elevator(){
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
