package assigner

import (
	"encoding/json"
	"log"
	"os/exec"

	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevio"
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

func elevioMD2string(d elevio.MotorDirection) string {
	return map[elevio.MotorDirection]string{elevio.MD_Up: "up", elevio.MD_Down: "down", elevio.MD_Stop: "stop"}[d]

}

func newHraInput(cs central.CentralState) hraInput {
	hrai := hraInput{}
	for i := range cs.HallOrders {
		hrai.HallRequests[i] = [2]bool{cs.HallOrders[i][0].Active, cs.HallOrders[i][1].Active}
	}
	hrai.States = make(map[string]hraState)
	for id, state := range cs.States {
		hrai.States[id] = hraState{
			Behaviour:   string(state.Behaviour),
			Floor:       state.Floor,
			Direction:   elevioMD2string(state.Direction),
			CabRequests: cs.CabOrders[id],
		}
	}
	return hrai
}

func hallRequestAssigner(cs central.CentralState) map[string][config.NUM_FLOORS][3]bool {
	b, err := json.Marshal(newHraInput(cs))
	if err != nil {
		log.Panic(err)
	}

	output, err := exec.Command("./hall_request_assigner", "-i", "--includeCab", string(b)).Output()
	if err != nil {
		log.Panic(err)
	}

	orders := make(map[string][config.NUM_FLOORS][3]bool)
	err = json.Unmarshal(output, &orders)
	if err != nil {
		log.Panic(err)
	}
	return orders
}

func Assigner(cs central.CentralState) [config.NUM_FLOORS][3]bool {
	return hallRequestAssigner(cs)[cs.Origin]
}
