package orders

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
)

type hraInput struct {
	HallRequests [config.NUM_FLOORS][2]bool `json:"hallRequests"`
	States       map[string]hraiState       `json:"states"`
}

type hraiState struct {
	Behaviour   string                  `json:"behaviour"`
	Floor       int                     `json:"floor"`
	Direction   string                  `json:"direction"`
	CabRequests [config.NUM_FLOORS]bool `json:"cabRequests"`
}

func csBehaviour2hraiBehaviour(b elevator.Behaviour) string {
	return map[elevator.Behaviour]string{
		elevator.Idle:                      "idle",
		elevator.DoorOpen:                  "doorOpen",
		elevator.DoorOpenWithPendingTarget: "doorOpen",
		elevator.ServingOrder:              "moving",
	}[b]
}

func csDirection2hraiDirection(d elevio.MotorDirection) string {
	return map[elevio.MotorDirection]string{
		elevio.MD_Up:   "up",
		elevio.MD_Down: "down",
		elevio.MD_Stop: "stop",
	}[d]
}

func newHraInput(cs CentralState) hraInput {
	hrai := hraInput{}
	for i := range cs.HallOrders {
		hrai.HallRequests[i] = [2]bool{cs.HallOrders[i].Up.Active, cs.HallOrders[i].Down.Active}
	}
	hrai.States = make(map[string]hraiState)
	for i, es := range cs.States {
		hrai.States[fmt.Sprint(i)] = hraiState{
			Behaviour:   csBehaviour2hraiBehaviour(es.Behaviour),
			Floor:       es.Floor,
			Direction:   csDirection2hraiDirection(es.Direction),
			CabRequests: cs.CabOrders[i],
		}
	}
	return hrai
}

func hallRequestAssigner(cs CentralState) map[int][config.NUM_FLOORS][2]bool {
	b, err := json.Marshal(newHraInput(cs))
	if err != nil {
		log.Panic(err)
	}

	output, err := exec.Command("./hall_request_assigner", "-i", string(b)).Output()
	if err != nil {
		log.Panic(err)
	}

	hraiOutput := make(map[string][config.NUM_FLOORS][2]bool)
	err = json.Unmarshal(output, &hraiOutput)
	if err != nil {
		log.Panic(err)
	}

	orders := make(map[int][config.NUM_FLOORS][2]bool)
	for key, value := range hraiOutput {
		id, err := strconv.Atoi(key)
		if err != nil {
			log.Panic(err)
		}
		orders[id] = value
	}
	return orders
}

//hallRequestAssigner(newHraInput(cs))

func CalculateTarget(cs CentralState) (int, bool) {
	return 0, false
}

func DeactivateOrders(cs CentralState, floor int) CentralState {
	return CentralState{}
}
