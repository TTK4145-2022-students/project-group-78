package main

import "github.com/TTK4145-2022-students/project-group-78/central"

func main() {
	id := 1

	central := central.New()
	defer central.Stop()

	elevator := elevator.New(id, central.StateIn)
	distributor := distributor.New(id, central.StateIn)

	for {
		select {
		case state := <-central.StateOut:
			distributor.StateIn <- state
			elevator.Lights <- orders.SetOrderBoard(state)
			elevator.TargetOrder <- orders.CalculateOrder(state)			
		}
	}
}
