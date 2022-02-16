import "Driver-go/elevio"
const (
	Above relative_position = iota
	Neutral  
	Below 
)

type LocalState struct{
	current_floor int;
	relative_position relative_position;
	target int;
}

local_state = LocalState{-1, Neutral};

func update_local_state(motor_dir_up_down bool,current_floor chan int, target int){
	current_floor := <-current_floor;
	if motor_dir_up_down{

	}
}
func start_move_towards_target(new_target int){
	int current_floor = get_current_floor();
	if ready_to_serve_new_order{
		if current_floor == target {

		}
		else if current_floor < target {
			set_motor_dir(up);
		}
		else {
			set_motor_dir(down);
		}
	}
}

ready_to_serve_new_order = true;

func run_elevator(){
	current_floor_chan = make(chan int);
	target_chan = make(chan int);

	go PollFloorSensors(receiver_current_floor);
	go update_local_state(receiver_current_floor);

}
