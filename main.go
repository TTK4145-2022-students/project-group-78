package main

import (
	"log"
	"os"
	"time"

	"Network-go/network/bcast"

	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("lifty", "lifty.")
	id := *parser.Int("i", "id", &argparse.Options{Default: 0})
	bcastPort := *parser.Int("b", "broadcast-port", &argparse.Options{Default: 46952})
	elevatorPort := *parser.Int("e", "elevator-port", &argparse.Options{Default: 15657})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Panic(err)
	}

	//
	state := elevator.NewCentralState(id, elevator.ElevatorState{})
	elevator.Init(id, elevatorPort)

	bcastReceive, bcastSend := make(chan elevator.CentralState), make(chan elevator.CentralState)
	go bcast.Receiver(bcastPort, bcastReceive)
	go bcast.Transmitter(bcastPort, bcastSend)

	for {
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case s := <-elevator.StateOut:
			merge(state, s)
			bcastSend <- state
			//delay to ensure that package are sent before turning on lights etc...
			elevator.StateIn <- state

		case s := <-bcastReceive:
			merge(state, s)
			elevator.StateIn <- state

		case <-timer.C:
			bcastSend <- state
		}
	}
}

// Merge cs2 onto cs1
func merge(cs1 elevator.CentralState, cs2 elevator.CentralState) {
	cs1.Elevators[cs2.Origin] = cs2.Elevators[cs2.Origin]

	for o := range cs2.HallOrders {
		if cs2.HallOrders[o].After(cs1.HallOrders[o]) {
			cs1.HallOrders[o] = cs2.HallOrders[o]
		}
	}

	for o := range cs2.ServedHallOrders {
		if cs2.ServedHallOrders[o].After(cs1.ServedHallOrders[o]) {
			cs1.ServedHallOrders[o] = cs2.ServedHallOrders[o]
		}
	}
}
