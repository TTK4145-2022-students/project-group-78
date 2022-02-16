import "Driver-go/elevio"


type relative_position int

const (
	Above relative_position = iota
	Neutral  
	Below 
)

type LocalState struct{
	target int;
	current_floor int;
	relative_position relative_position;
	doorOpen bool;
	motor_dir MotorDirection;
}

local_state = LocalState{-1,-1, Neutral, false, MD_Stop};

func run_elevator(){
	addr = "11111"
	numFloors = 5
	Init(addr, numFloors);

	for {
		if (event){
			switch event.type{

				case floorSensorTriggered:
					current_floor == event.floor
					if current_floor == target{
						elevator_finished_order();
						go door_open();
					}
					switch local_state.motor_dir{
						case MD_Up:
							local_state.relative_position = Above;
						case MD_Down:
							local_state.relative_position = Below;
						case MD_Stop:
							local_state.relative_position = Neutral;
						default:
					}

				case newTarget:
					start_motor_towards_target(event.target)
					local_state.target = event.target

				default:
			}
		}
	}
}
