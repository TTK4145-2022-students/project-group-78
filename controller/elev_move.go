
func start_motor_towards_target(new_target int){
	int current_floor = get_current_floor();
	if ready_to_serve_new_order{
		if current_floor == target {
			switch relative_position
				case Above:
					setMotorDirection(MD_Down);
				case Neutral:
					SetMotorDirection(MD_Stop);
					
				case Below:
					SetMotorDirection(MD_Up);

		}
		else if current_floor < target {
			SetMotorDirection(MD_Up);
		}
		else {
			SetMotorDirection(MD_Down);
		}
	}
}