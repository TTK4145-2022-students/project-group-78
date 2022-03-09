package node

import (
	"fmt"

	"github.com/TTK4145-2022-students/project-group-78/assigner"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/lights"
)

func Node(id string, port int, inC <-chan central.CentralState, outC chan<- central.CentralState) {
	newOrderC := make(chan elevio.ButtonEvent)
	orderCompletedC, stateC := make(chan elevio.ButtonEvent), make(chan elevator.State)
	assignedOrdersC := make(chan [config.NUM_FLOORS][3]bool)

	elevio.Init(fmt.Sprintf("127.0.0.1:%v", port), config.NUM_FLOORS)
	go elevator.Elevator(assignedOrdersC, orderCompletedC, stateC)
	go elevio.PollButtons(newOrderC)

	cs := central.New(id)
	for {
		select {
		case o := <-newOrderC:
			cs = cs.SetOrder(o, true)

		case o := <-orderCompletedC:
			cs = cs.SetOrder(o, false)

		case s := <-stateC:
			cs.States[id] = s

		case newCs := <-inC:
			cs = cs.Merge(newCs)
		}
		assignedOrdersC <- assigner.Assigner(cs)
		outC <- cs
		lights.Set(cs)
	}
}
