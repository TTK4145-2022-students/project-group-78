package distributor

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

const PORT = 41875

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

	d := new(Distributor)
	var err error

	d.Conn, err = net.ListenUDP("udp", localAddr)
	if err == nil {
		log.Debugf("Listening on %v", d.Conn.LocalAddr().String())
	} else {
		log.Panic(err)
	}

	d.Id = id
	return d
}

func (d *Distributor) SendDatagram(datagram Datagram) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(datagram)

	ip := net.ParseIP("127.255.255.255")
	addr := &net.UDPAddr{IP: ip, Port: PORT}

	n, err := d.Conn.WriteToUDP(buf.Bytes(), addr)
	if err != nil {
		log.Panic(err)
	} else {
		log.Debugf("Successfully sent %v bytes to %v", n, addr.String())
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

	decoder := gob.NewDecoder(bytes.NewBuffer(buf[0:n]))
	err = decoder.Decode(&datagram)
	if err != nil {
		log.Panic(err)
	}
	return
}
