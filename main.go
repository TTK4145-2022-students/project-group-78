package main

import (
	"Network-go/network/bcast"
	"flag"
	"fmt"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/assigner"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/lights"
)

func clParams() (id string, bcastPort int, elevPort int) {
	idP := flag.String("id", "elevator", "elevator id")
	bcastPortP := flag.Int("bcastPort", 56985, "broadcast port")
	elevPortP := flag.Int("elevPort", 15657, "elevator port")
	flag.Parse()
	return *idP, *bcastPortP, *elevPortP
}

func main() {
	id, bcastPort, elevPort := clParams()
	newOrderC, orderCompletedC := make(chan elevio.ButtonEvent), make(chan elevio.ButtonEvent, 16)
	stateC := make(chan elevator.State, 16)
	assignedOrdersC := make(chan elevator.Orders, 16)
	sendC, receiveC := make(chan central.CentralState), make(chan central.CentralState)

	elevio.Init(fmt.Sprintf("127.0.0.1:%v", elevPort), config.NUM_FLOORS)
	lights.Clear()
	go elevator.Elevator(assignedOrdersC, orderCompletedC, stateC)
	go elevio.PollButtons(newOrderC)
	go bcast.Transmitter(bcastPort, sendC)
	go bcast.Receiver(bcastPort, receiveC)

	cs := central.New(id, <-stateC)
	for {
		select {
		case o := <-newOrderC:
			cs = cs.SetOrder(o, true)
			sendC <- cs

		case o := <-orderCompletedC:
			cs = cs.SetOrder(o, false)
			sendC <- cs

		case s := <-stateC:
			cs.States[id] = s
			sendC <- cs

		case newCs := <-receiveC:
			if newCs.Origin == id {
				continue
			}
			cs = cs.Merge(newCs)

		case <-time.After(config.TRANSMIT_INTERVAL):
			sendC <- cs
			continue
		}
		assignedOrdersC <- assigner.Assigner(cs)
		lights.Set(cs)
	}
}
