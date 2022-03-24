package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/TTK4145-2022-students/Network-go-group-78/network/bcast"
	"github.com/TTK4145-2022-students/driver-go-group-78/elevio"
	"github.com/TTK4145-2022-students/project-group-78/assigner"
	"github.com/TTK4145-2022-students/project-group-78/central"
	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/elevator"
	"github.com/TTK4145-2022-students/project-group-78/lights"
	"github.com/rapidloop/skv"
)

func main() {
	idP := flag.Int("id", 0, "elevator id")
	portP := flag.Int("port", 15657, "elevator port")
	flag.Parse()

	id := *idP
	assignedOrdersC := make(chan elevator.Orders, config.ChanSize)
	stateC := make(chan elevator.State, config.ChanSize)
	newOrderC, orderCompletedC := make(chan elevio.ButtonEvent), make(chan elevio.ButtonEvent, config.ChanSize)
	sendC, receiveC := make(chan central.CentralState), make(chan central.CentralState)

	elevio.Init(fmt.Sprintf("127.0.0.1:%v", *portP), config.NumFloors)
	lights.Clear()
	go elevator.Elevator(assignedOrdersC, orderCompletedC, stateC)
	go elevio.PollButtons(newOrderC)
	go bcast.Transmitter(config.BcastPort, sendC)
	go bcast.Receiver(config.BcastPort, receiveC)

	store, err := skv.Open(fmt.Sprintf("elev%v.db", id))
	if err != nil {
		panic(err)
	}
	var cs central.CentralState
	if err = store.Get("cs", &cs); err != nil && err != skv.ErrNotFound {
		panic(err)
	}
	cs.Origin = id

	ticker := time.NewTicker(config.TransmitInterval)
	for {
		select {
		case o := <-newOrderC:
			cs = cs.AddOrder(o)
			sendC <- cs

		case o := <-orderCompletedC:
			cs = cs.RemoveOrder(o)
			sendC <- cs

		case s := <-stateC:
			cs.States[id] = s
			cs.LastUpdated[id] = time.Now()
			sendC <- cs

		case newCs := <-receiveC:
			if newCs.Origin == id {
				continue
			}
			cs = cs.Merge(newCs)

		case <-ticker.C:
			sendC <- cs
			continue
		}

		if err = store.Put("cs", cs); err != nil {
			panic(err)
		}
		assignedOrdersC <- assigner.Assigner(cs)
		// Delay lights so that we ensure that sufficent attemps have been done to send new orders to the other nodes.
		// Orders are anyway stored in persistant storage, so no orders can be lost anyway, but this is to ensure that the spec is satisfied
		go func() {
			time.Sleep(config.LightDelay)
			lights.Set(cs)
		}()
	}
}
