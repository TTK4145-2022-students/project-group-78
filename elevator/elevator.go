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

func Init(id_ int, port int) {
	id = id_
	elevio.Init(fmt.Sprintf("127.0.0.1:%v", port), config.NUM_FLOORS)

	go elevio.PollButtons(buttonPressedC)
	go elevio.PollFloorSensor(floorEnteredC)
	go elevio.PollObstructionSwitch(doorObstructionC)
	go run()
}

func run() {
	obstructed := false
	elevatorState := ElevatorState{State: Moving, Direction: elevio.MD_Up}
	doorTimer := time.NewTimer(time.Hour)
	doorTimer.Stop()

	for {
		select {
		case be := <-buttonPressedC:
			cs := buttonPressed(be, elevatorState)
			elevatorState = cs.Elevators[id]
			StateOut <- cs

		case f := <-floorEnteredC:
			cs := floorEntered(f, elevatorState, doorTimer)
			elevatorState = cs.Elevators[id]
			StateOut <- cs

		case obstructed = <-doorObstructionC:
			if !obstructed && elevatorState.State == DoorOpen {
				elevatorState = closeDoor(elevatorState)
				StateOut <- NewCentralState(id, elevatorState)
			}

		case <-doorTimer.C:
			if !obstructed && elevatorState.State == DoorOpen {
				elevatorState = closeDoor(elevatorState)
				StateOut <- NewCentralState(id, elevatorState)
			}

		case cs := <-StateIn:
			target, notEmpty := calculateTargetOrder(cs)
			if !notEmpty && target != elevatorState.Target {
				elevatorState.Target = target
				cs = targetOrderUpdated(elevatorState, doorTimer)
				elevatorState = cs.Elevators[id]
				StateOut <- cs
			}

			//TODO: add delay to lights
			setLights(cs)
		}
	}
}

var lastLightValues map[Order]bool

func init() {
	lastLightValues = make(map[Order]bool, 1)
}

func setLights(cs CentralState) {
	for o := range cs.HallOrders {
		if value := cs.HallOrders[o].After(cs.ServedHallOrders[o]); value != lastLightValues[o] {
			elevio.SetButtonLamp(o.Type, o.Floor, value)
			lastLightValues[o] = value
		}
	}
}
