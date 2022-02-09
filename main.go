package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const PORT = 41875
const BUF_LENGTH = 2048
var MULTICAST_ADDR = &net.UDPAddr{IP: net.ParseIP("127.255.255.255"), Port: PORT}

type Distributor struct {
	EventsIn chan string
	Conn     *net.UDPConn
}

func main() {
	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Panic(err)
	}

	ip := net.ParseIP(fmt.Sprintf("127.0.0.%v", id))
	addr := &net.UDPAddr{IP: ip, Port: PORT}

	conn, err := net.ListenUDP("udp", addr)
	if err == nil {
		log.Debugf("Listening on %v", conn.LocalAddr().String())
	} else {
		log.Panic(err)
	}

	_, err = conn.WriteToUDP([]byte(fmt.Sprintf("%v hello", id)), MULTICAST_ADDR)
	if err != nil {
		log.Panic(err)
	}

	processUDP(conn)
}

func logConn(conn *net.UDPConn) {
	buf := make([]byte, BUF_LENGTH)
	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		log.Panic(err)
	} else if n == BUF_LENGTH {
		log.Panic("Read max number of bytes (%v) into buffer. Consider increasing the buffer size", BUF_LENGTH)
	} else {
		log.Infof("%v from %v", string(buf), addr.String())
		//log.Debugf("Received %v bytes from %v", n, addr.String())
	}
}

func processUDP(conn *net.UDPConn) {
	for {
		logConn(conn)
	}
}
