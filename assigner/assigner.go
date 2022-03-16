package assigner

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
)

type hraInput struct {
	HallRequests [config.NUM_FLOORS][2]bool `json:"hallRequests"`
	States       map[string]hraState        `json:"states"`
}

type hraState struct {
	Behaviour   string                  `json:"behaviour"`
	Floor       int                     `json:"floor"`
	Direction   string                  `json:"direction"`
	CabRequests [config.NUM_FLOORS]bool `json:"cabRequests"`
}

func elevioMd2String(d elevio.MotorDirection) string {
	return map[elevio.MotorDirection]string{elevio.MD_Up: "up", elevio.MD_Down: "down", elevio.MD_Stop: "stop"}[d]
}

func elevatorBehaviour2String(b elevator.Behaviour) string {
	return map[elevator.Behaviour]string{elevator.Idle: "idle", elevator.DoorOpen: "doorOpen", elevator.Moving: "moving"}[b]
}

func newHraInput(cs central.CentralState) hraInput {
	hrai := hraInput{}
	for f := range cs.HallOrders {
		hrai.HallRequests[f] = [2]bool{cs.HallOrders[f][0].Active, cs.HallOrders[f][1].Active}
	}
	hrai.States = make(map[string]hraState)
	for id, state := range cs.States {
		hrai.States[strconv.Itoa(id)] = hraState{
			Behaviour:   elevatorBehaviour2String(state.Behaviour),
			Floor:       state.Floor,
			Direction:   elevioMd2String(state.Direction),
			CabRequests: cs.CabOrders[id],
		}
	}
	return hrai
}

func hallRequestAssigner(cs central.CentralState) map[string]elevator.Orders {
	b, err := json.Marshal(newHraInput(cs))
	if err != nil {
		log.Panic(err)
	}

	output, err := exec.Command("hall_request_assigner", "-i", "--includeCab", string(b)).Output()
	if err != nil {
		log.Panic(err)
	}

	orders := make(map[string]elevator.Orders)
	err = json.Unmarshal(output, &orders)
	if err != nil {
		log.Panic(err)
	}
	return orders
}

func Assigner(cs central.CentralState) elevator.Orders {
	return hallRequestAssigner(cs)[strconv.Itoa(cs.Origin)]
}
