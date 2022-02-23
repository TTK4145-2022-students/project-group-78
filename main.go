package main

import (
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/distributor"
)

func main() {
	var id byte = 1

	central := central.New()
	defer central.Stop()

	elevator := elevator.New(id, central.StateIn)
	distributor := distributor.New(id, central.StateIn)

	for {
		select {
		case s := <-central.StateOut:
			distributor.Send(s)
			elevator.Lights <- orders.SetOrderBoard(s)
			elevator.TargetOrder <- orders.CalculateOrder(s)
		}
	}
}
