package orders

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

func Orders(id int, csC chan CentralState) {
	buttonPressedC := make(chan elevio.ButtonEvent)
	for {
		select {
		case be := <-buttonPressedC:
			cs := CentralState{Origin: id}
			switch be.Button {
			case elevio.BT_Cab:
				cs.CabOrders[id][be.Floor] = true

			case elevio.BT_HallUp:
				cs.HallOrders[be.Floor].Up.Active = true
				cs.HallOrders[be.Floor].Up.Time = time.Now()

			case elevio.BT_HallDown:
				cs.HallOrders[be.Floor].Down.Active = true
				cs.HallOrders[be.Floor].Down.Time = time.Now()
			}
			csC <- cs
		}
	}
}
