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
var SetOrderLight chan []order.OrderLight

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

var obstructed bool  = false

func run() {
	for {
		select {
		case be := <-buttonPressedC:
			order := order.Order{
				Floor:     be.Floor,
				OrderType: order.OrderType(be.Button),
			}
			emit(events.OrderReceived{Order: order})

		case f := <-floorEnteredC:
			floorEntered(f,&doorTimer)

		case obstructed := <-doorObstructionC:
			if !obstructed && state = DoorOpen{
				closeDoor()
			}

		case <-doorTimer.C:
			if !obstructed{
				closeDoor()
			}
		case o := <-SetTargetOrder:
			targetOrderUpdated(o)

		case orderLights := <-SetOrderLight:
			for _, orderLight := range orderLights {
				order := orderLight.Order
				elevio.SetButtonLamp(elevio.ButtonType(order.OrderType), order.Floor, orderLight.Value)
			}
		}
	}
}

// Creates a new CentralState sends it out on StateUpdate
func emit(e central.Event) {
	cs := central.NewCentralState()
	cs[id][e] = time.Now()
	StateUpdate <- cs
}
