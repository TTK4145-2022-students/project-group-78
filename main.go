package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"Network-go/network/bcast"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/lights"
	"github.com/TTK4145-2022-students/project-group-78/orders"
	"github.com/akamensky/argparse"
)

func clParams() (id int, bcastPort int, elevatorPort int) {
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
	id, bcastPort, elevatorPort := clParams()
	elevio.Init(fmt.Sprintf("127.0.0.1:%v", elevatorPort), config.NUM_FLOORS)

	targetReachedC := make(chan int)
	stateC := make(chan elevator.State)
	newCsC, bcastCsC := make(chan orders.CentralState), make(chan orders.CentralState)

	go elevator.Elevator(targetReachedC, stateC)
	go orders.Orders(id, newCsC)
	go bcast.Receiver(bcastPort, newCsC)
	go bcast.Transmitter(bcastPort, bcastCsC)

	cs := orders.CentralState{Origin: id}

	for {
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case f := <-targetReachedC:
			cs = orders.DeactivateOrders(cs, f)
			bcastCsC <- cs
			newCsC <- cs

		case s := <-stateC:
			cs.States[id] = s
			bcastCsC <- cs
			newCsC <- cs

		case newCs := <-newCsC:
			cs = cs.Merge(newCs)
			if target, ok := orders.CalculateTarget(cs); ok {
				elevator.TargetC <- target
			}
			lights.SetC <- cs

		case <-timer.C:
			bcastCsC <- cs
		}
	}
}
