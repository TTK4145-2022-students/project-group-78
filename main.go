package main

import (
	"log"
	"os"
	"time"

	"Network-go/network/bcast"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/node"
	"github.com/TTK4145-2022-students/project-group-78/orders"
	"github.com/akamensky/argparse"
)

func clParams() (id int, bcastPort int, elevatorPort int) {
	parser := argparse.NewParser("lifty", "lifty.")
	id = *parser.Int("i", "id", &argparse.Options{Default: 0})
	bcastPort = *parser.Int("b", "broadcast-port", &argparse.Options{Default: 46952})
	elevatorPort = *parser.Int("e", "elevator-port", &argparse.Options{Default: 15657})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Panic(err)
	}
	return
}

func main() {
	id, bcastPort, elevatorPort := clParams()
	nodeOutC := make(chan orders.CentralState)
	bcastTransmitC, bcastReceiveC := make(chan orders.CentralState), make(chan orders.CentralState)

	node.Node(id, elevatorPort, nodeOutC)
	go bcast.Receiver(bcastPort, bcastReceiveC)
	go bcast.Transmitter(bcastPort, bcastTransmitC)

	cs := orders.CentralState{Origin: id}

	for {
		select {
		case newCs := <-bcastReceiveC:
			cs = cs.Merge(newCs)
			node.InC <- cs

		case newCs := <-nodeOutC:
			cs = cs.Merge(newCs)
			bcastTransmitC <- cs

		case <-time.After(config.TRANSMIT_INTERVAL):
			bcastTransmitC <- cs
		}
	}
}
