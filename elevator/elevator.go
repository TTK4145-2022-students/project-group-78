package elevator

import (
	"fmt"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/events"
	"github.com/TTK4145-2022-students/project-group-78/order"
)

var StateUpdate chan central.CentralState
var SetTargetOrder chan order.Order
var SetOrderLight chan struct {
	order.Order
	bool
}

var id int
var buttonPressedC chan elevio.ButtonEvent
var floorEnteredC chan int
var doorObstructionC chan bool
var doorTimer *time.Timer

func Init(id_ int, port int) {
	id = id_
	elevio.Init(fmt.Sprintf("127.0.0.1:%v", port), config.NUM_FLOORS)

	go elevio.PollButtons(buttonPressedC)
	go elevio.PollFloorSensor(floorEnteredC)
	go elevio.PollObstructionSwitch(doorObstructionC)

	doorTimer = time.NewTimer(time.Hour)
	doorTimer.Stop()

	go run()
}

var obstructed bool

func run() {
	for {
		select {
		case be := <-buttonPressedC:
			order := order.Order{
				Floor:     be.Floor,
				OrderType: order.OrderType(be.Button),
			}
			emit(events.OrderReceived{order})

		case f := <-floorEnteredC:
			floorEntered(f)

		case obstructed = <-doorObstructionC:

		case <-doorTimer.C:
			doorTimedOut()

		case o := <-SetTargetOrder:
			targetOrderUpdated(o)

		case l := <-SetOrderLight:
			order, value := l.Order, l.bool
			elevio.SetButtonLamp(elevio.ButtonType(order.OrderType), order.Floor, value)
		}
	}
}

// Creates a new CentralState sends it out on StateUpdate
func emit(e central.Event) {

}
