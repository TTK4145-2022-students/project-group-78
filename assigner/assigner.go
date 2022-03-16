// Assigns order to the elevator running on the node using hall_requst_assigner (hra).
// Faulty elevators are detected and their orders are reassigned.
// A faulty elevator is one that has not moved in a while, and has an order which is not entirly fresh.
// The rationale for this is that an elevator should move when it gets a new order, otherwise it is faulty.

package assigner

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"time"

	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
)

type hraInput struct {
	HallRequests [config.NumFloors][2]bool `json:"hallRequests"`
	States       map[string]hraState       `json:"states"`
}

type hraState struct {
	Behaviour   string                 `json:"behaviour"`
	Floor       int                    `json:"floor"`
	Direction   string                 `json:"direction"`
	CabRequests [config.NumFloors]bool `json:"cabRequests"`
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
		hrai.HallRequests[f] = [2]bool{cs.HallOrders[f][elevio.BT_HallUp].Active, cs.HallOrders[f][elevio.BT_HallDown].Active}
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
		panic(err)
	}

	output, err := exec.Command("hall_request_assigner", "-i", "--includeCab", string(b)).Output()
	if err != nil {
		panic(err)
	}

	orders := make(map[string]elevator.Orders)
	err = json.Unmarshal(output, &orders)
	if err != nil {
		panic(err)
	}
	return orders
}

func Assigner(cs central.CentralState) elevator.Orders {
	hrai := newHraInput(cs)
	for c := 0; ; c++ {
		// Increase order timeout for each removed elevator, so that not multiple elevators times out the same order
		key, faulty := findOtherFaultyElevator(cs, hrai, time.Duration(c+1)*config.OrderTimout)
		if faulty {
			delete(hrai.States, key)
		} else {
			return hallRequestAssigner(hrai)[strconv.Itoa(cs.Origin)]
		}
	}
}

func findOtherFaultyElevator(cs central.CentralState, hrai hraInput, orderTimeout time.Duration) (key string, faulty bool) {
	for key, orders := range hallRequestAssigner(hrai) {
		if key == strconv.Itoa(cs.Origin) {
			continue
		}
		for f := range orders {
			for bt := range orders[f] {
				id, err := strconv.Atoi(key)
				if err != nil {
					panic(err)
				}
				if orders[f][bt] &&
					bt != elevio.BT_Cab &&
					time.Since(cs.HallOrders[f][bt].Time) > orderTimeout &&
					time.Since(cs.LastUpdated[id]) > config.ElevTimeout {
					// If we have an old hall order and the assigned elevator has not responeded, conclude that it is faulty
					return key, true
				}
			}
		}
	}
	return
}
