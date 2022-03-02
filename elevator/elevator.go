package elevator

import (
	"fmt"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

var StateOut chan CentralState
var StateIn chan CentralState

var id int
var buttonPressedC chan elevio.ButtonEvent
var floorEnteredC chan int
var doorObstructionC chan bool
var doorTimer *time.Timer

func Init(id int, port int) {
	id = id
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
			if be.Button == elevio.BT_Cab {
				es.CabOrders[be.Floor] = true
				StateOut <- NewCentralState(id, es)
			} else {
				cs := NewCentralState(id, es)
				o := Order{
					Type:  be.Button,
					Floor: be.Floor,
				}
				cs.Orders[o] = time.Now()
				StateOut <- cs
			}

		case f := <-floorEnteredC:
			floorEntered(f)
			// TODO: must send new cs and order served

		case obstructed := <-doorObstructionC:
			if !obstructed && state = DoorOpen{
				closeDoor()
			}

		case <-doorTimer.C:
			doorTimedOut()
			// TODO: send new cs

		case cs := <-StateIn:
			order = assigner.CalculateTarget(cs)
			targetOrderUpdated(order)
			// TODO: must send new cs and order served
			setLights(cs)
		}
	}
}

var lastLightValues map[Order]bool

func setLights(cs CentralState) {
	for o, _ := range cs.Orders {
		if value := cs.Orders[o].After(cs.ServedOrders[o]); value != lastLightValues[o] {
			elevio.SetButtonLamp(o.Type, o.Floor, value)
			lastLightValues[o] = value
		}
	}
}
