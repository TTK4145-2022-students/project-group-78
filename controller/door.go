import (
	"Driver-go/elevio"
	"time"
)

func door_open(){
	local_state.doorOpen =true
	SetDoorOpenLamp(true)
	timer := time.NewTimer(3*time.Second)
	<-timer.C
	local_state.doorOpen = false
	SetDoorOpenLamp(false)
}