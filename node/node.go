package node

import (
	"fmt"

	"github.com/TTK4145-2022-students/project-group-78/assigner"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

func Node(id string, port int, inC <-chan central.CentralState, outC chan<- central.CentralState) {
	newOrderC, orderCompletedC := make(chan elevio.ButtonEvent), make(chan elevio.ButtonEvent)
	stateC := make(chan elevator.State)
	assignedOrdersC := make(chan elevator.Orders)

	elevio.Init(fmt.Sprintf("127.0.0.1:%v", port), config.NUM_FLOORS)
	go elevator.Elevator(assignedOrdersC, orderCompletedC, stateC)
	go elevio.PollButtons(newOrderC)

	cs := central.New(id)
	for {
		select {
		case o := <-newOrderC:
			cs = cs.SetOrder(o, true)
			outC <- cs

		case o := <-orderCompletedC:
			cs = cs.SetOrder(o, false)
			outC <- cs

		case s := <-stateC:
			cs.States[id] = s
			outC <- cs

		case newCs := <-inC:
			cs = cs.Merge(newCs)
		}
		assignedOrdersC <- assigner.Assigner(cs)
		setLights(cs)
	}
}
