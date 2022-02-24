package main

import (
	"log"
	"os"
	"time"

	"Network-go/network/bcast"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/order"
	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("lifty", "lifty.")
	id := *parser.Int("i", "id", &argparse.Options{Default: 1})
	bcastPort := *parser.Int("b", "broadcast-port", &argparse.Options{Default: 46952})
	elevatorPort := *parser.Int("e", "elevator-port", &argparse.Options{Default: 15657})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Panic(err)
	}

	//
	state := central.NewCentralState()
	elevator.Init(id, elevatorPort)

	stateIn, stateOut := make(chan central.CentralState), make(chan central.CentralState)
	go bcast.Receiver(bcastPort, stateIn)
	go bcast.Transmitter(bcastPort, stateOut)

	for {
		timer := time.NewTimer(100 * time.Millisecond)
		select {
		case s := <-elevator.StateUpdate:
			state.Merge(s)
			stateOut <- state
			elevator.SetOrderLight <- order.CalculateOrderLights(state)
			elevator.SetTargetOrder <- order.CalculateTargetOrder(state)

		case s := <-stateIn:
			state.Merge(s)
			elevator.SetOrderLight <- order.CalculateOrderLights(state)
			elevator.SetTargetOrder <- order.CalculateTargetOrder(state)

		case <-timer.C:
			stateOut <- state
		}
	}
}
