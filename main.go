package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"Network-go/network/bcast"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/door"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/lights"
	"github.com/TTK4145-2022-students/project-group-78/orders"
	"github.com/akamensky/argparse"
)

func parseArgs() (id int, bcastPort int, elevatorPort int) {
	parser := argparse.NewParser("lifty", "lifty.")
	id = *parser.Int("i", "id", &argparse.Options{Default: 0})
	bcastPort = *parser.Int("b", "broadcast-port", &argparse.Options{Default: 46952})
	elevatorPort = *parser.Int("e", "elevator-port", &argparse.Options{Default: 15657})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Panic(err)
	}
	return
}

func main() {
	id, bcastPort, elevatorPort := parseArgs()
	elevio.Init(fmt.Sprintf("127.0.0.1:%v", elevatorPort), config.NUM_FLOORS)

	doorClosedC := make(chan bool)
	floorEnteredC, targetReachedC := make(chan int), make(chan int)
	directionSetC := make(chan elevio.MotorDirection)
	newCsC, bcastCsC := make(chan central.CentralState), make(chan central.CentralState)

	go door.Door(doorClosedC)
	go elevator.Elevator(floorEnteredC, targetReachedC, directionSetC)
	go orders.Orders(id, newCsC)
	go bcast.Receiver(bcastPort, newCsC)
	go bcast.Transmitter(bcastPort, bcastCsC)

	cs := central.CentralState{Origin: id}

	for {
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case <-doorClosedC:
			if cs.Elevators[id].State == central.DoorOpen {
				cs.Elevators[id].State = central.Idle
				bcastCsC <- cs
				newCsC <- cs
			} else {
				log.Printf("error: Door closed while in %v", cs.Elevators[id].State)
			}

		case f := <-floorEnteredC:
			cs.Elevators[id].Floor = f
			bcastCsC <- cs
			newCsC <- cs

		case f := <-targetReachedC:
			if cs.Elevators[id].State == central.ServingOrder {
				door.OpenC <- true
				cs.Elevators[id].State = central.DoorOpen
				cs = orders.DeactivateOrders(cs, f)
				bcastCsC <- cs
				newCsC <- cs
			} else {
				log.Printf("error: Reached target while in %v", cs.Elevators[id].State)
			}

		case d := <-directionSetC:
			cs.Elevators[id].Direction = d
			bcastCsC <- cs
			newCsC <- cs

		case newCs := <-newCsC:
			cs.Merge(newCs)
			switch cs.Elevators[id].State {
			case central.Idle:
				target, ok := orders.CalculateTarget(cs)
				if ok {
					elevator.TargetC <- target
					cs.Elevators[id].State = central.ServingOrder
					bcastCsC <- cs
					newCsC <- cs
				}

			case central.DoorOpen:

			case central.ServingOrder:
				if target, ok := orders.CalculateTarget(cs); ok {
					elevator.TargetC <- target
				}
			}
			lights.SetC <- cs

		case <-timer.C:
			bcastCsC <- cs
		}
	}
}
