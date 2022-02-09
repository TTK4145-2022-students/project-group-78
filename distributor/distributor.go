package distributor

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

const PORT = 41875
var BROADCAST_ADDR = &net.UDPAddr{IP: net.ParseIP("127.255.255.255"), Port: PORT}

type Distributor struct {
	Conn *net.UDPConn
	Id int
}

type Datagram struct {
	SequenceNum int
	From        int
	Ack         bool
	Message     string
}

func New(id int) *Distributor {
	localIp := net.ParseIP(fmt.Sprintf("127.0.0.%v", id))
	localAddr := &net.UDPAddr{IP: localIp, Port: PORT}

	conn, err := net.ListenUDP("udp", localAddr)
	if err == nil {
		log.Debugf("Listening on %v", conn.LocalAddr().String())
	} else {
		log.Panic(err)
	}

	return &Distributor{conn, id}
}

func (d *Distributor) SendDatagram(datagram Datagram) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(datagram)
	if err != nil {
		log.Panic(err)
	}

	n, err := d.Conn.WriteToUDP(buf.Bytes(), BROADCAST_ADDR)
	if err != nil {
		log.Panic(err)
	} else {
		log.Debugf("Successfully sent %v bytes to %v", n, BROADCAST_ADDR.String())
	}
}

func (d *Distributor) ReceiveDatagram() (datagram Datagram) {
	buf_size := 2048
	buf := make([]byte, buf_size)

	n, addr, err := d.Conn.ReadFrom(buf)
	if err != nil {
		log.Panic(err)
	} else if n == buf_size {
		log.Panic("Read max number of bytes (%v) into buffer. Consider increasing the buffer size", buf_size)
	} else {
		log.Debugf("Received %v bytes from %v", n, addr.String())
	}

	err = gob.NewDecoder(bytes.NewBuffer(buf[0:n])).Decode(&datagram)
	if err != nil {
		log.Panic(err)
	}
	return
}
