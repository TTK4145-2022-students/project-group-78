package main

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/distributor"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
)

func main() {
	var id int = 1

	state := central.NewCentralState()

	elevator := elevator.New(id)
	distributor := distributor.New(id)

	for {
		select {
		case s := <-elevator.StateUpdate:
			state.Merge(s)

		case s := <-distributor.StateUpdate:
			state.Merge(s)

		default:
			time.Sleep(10 * time.Millisecond)
		}
		distributor.Send(state)
		elevator.Lights <- orders.SetOrderBoard(state)
		elevator.TargetOrder <- orders.CalculateOrder(state)
	}
}
