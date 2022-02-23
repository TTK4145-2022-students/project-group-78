import (
	"Driver-go/elevio"
	"time"
)

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