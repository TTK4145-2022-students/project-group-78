package distributor

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

const PORT = 41875
const BUF_LENGTH = 2048

type Distributor struct {
	EventsIn chan string
	Conn     *net.UDPConn
}

func New(id int) (distributor Distributor) {
	ip := net.ParseIP(fmt.Sprintf("127.0.0.%v", id))
	addr := &net.UDPAddr{IP: ip, Port: PORT}

	conn, err := net.ListenUDP("udp", addr)
	if err == nil {
		log.Debugf("Listening on %v", conn.LocalAddr().String())
	} else {
		log.Panic(err)
	}

	go processUDP(conn)

	return
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

func processEvents() {

}
