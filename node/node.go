package node

import (
	"fmt"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/lights"
	"github.com/TTK4145-2022-students/project-group-78/orders"
)

var InC chan orders.CentralState

func Node(id int, port int, outC chan orders.CentralState) {
	InC = make(chan orders.CentralState)
	targetReachedC, stateC := make(chan int), make(chan elevator.State)
	csOrderC := make(chan orders.CentralState)

	elevio.Init(fmt.Sprintf("127.0.0.1:%v", port), config.NUM_FLOORS)
	go elevator.Elevator(targetReachedC, stateC)
	go orders.Orders(id, csOrderC)

	cs := orders.CentralState{Origin: id}
	for {
		select {
		case f := <-targetReachedC:
			cs = orders.DeactivateOrders(cs, f)
			outC <- cs

		case cs.States[id] = <-stateC:
			outC <- cs

		case newCs := <-csOrderC:
			cs = cs.Merge(newCs)
			outC <- cs

		case newCs := <-InC:
			cs = cs.Merge(newCs)
		}
		if target, ok := orders.CalculateTarget(cs); ok {
			elevator.TargetC <- target
		}
		lights.SetC <- cs
	}
}
