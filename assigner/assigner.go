package assigner

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"time"

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

func hallRequestAssigner(hrai hraInput) map[string]elevator.Orders {
	b, err := json.Marshal(hrai)
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
	hrai := newHraInput(cs)

	for c := 0; ; c++ {
		e, ok := otherFaultyElevator(cs, hrai, time.Duration(c+1)*config.ORDER_TIMEOUT)
		if ok {
			delete(hrai.States, strconv.Itoa(e))
		} else {
			return hallRequestAssigner(hrai)[strconv.Itoa(cs.Origin)]
		}
	}
}

func otherFaultyElevator(cs central.CentralState, hrai hraInput, orderTimeout time.Duration) (e int, ok bool) {
	for elevator, orders := range hallRequestAssigner(hrai) {
		e, err := strconv.Atoi(elevator)
		if err != nil {
			log.Panic(err)
		}
		if e == cs.Origin {
			continue
		}
		for f := range orders {
			for bt := range orders[f] {
				if orders[f][bt] &&
					bt != elevio.BT_Cab &&
					time.Since(cs.HallOrders[f][bt].Time) > orderTimeout &&
					time.Since(cs.LastUpdated[e]) > config.ELEV_TIMEOUT {
					// If we have an old hall order and the assigned elevator has not responeded, conclude that it is faulty
					return e, true
				}
			}
		}
	}
	return
}
