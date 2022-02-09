package main

import (
	"net"
	"os"
	"strconv"

	"github.com/TTK4145-2022-students/project-group-78/distributor"
	log "github.com/sirupsen/logrus"
)

const PORT = 41875
const BUF_LENGTH = 2048

type Distributor struct {
	EventsIn chan string
	Conn     *net.UDPConn
}

func main() {
	log.SetLevel(log.DebugLevel)

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Panic(err)
	} else if id < 1 || id > 255 {
		log.Panic("Invalid id, must be between 1 and 255")
	}

	dist := distributor.New(id)
	datagram := &distributor.Datagram{1, 1, false, "Hello"}
	dist.SendDatagram(*datagram)
	datagram2 := dist.ReceiveDatagram()

	log.Infof("Sent %v - recieved %v", datagram.Message, datagram2.Message)
}
