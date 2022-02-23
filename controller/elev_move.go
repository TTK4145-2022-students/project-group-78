
func start_motor_towards_target(new_target int){
	if elev.current_floor == target {
		switch elev.relative_position{
			case Above:
				setMotorDirection(MD_Down);

			case Neutral:
				target_reached()

			case Below:
				SetMotorDirection(MD_Up);
		}
	}
	else if elev.current_floor < elev.target {
		SetMotorDirection(MD_Up);
	}
	else {
		SetMotorDirection(MD_Down);
	}
}
