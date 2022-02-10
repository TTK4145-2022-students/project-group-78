package distributor

import (
	"bytes"
	"encoding/gob"
	"net"

	log "github.com/sirupsen/logrus"
)

const PORT = 41875

var BROADCAST_ADDR = &net.UDPAddr{IP: net.ParseIP("127.255.255.255"), Port: PORT}

type Event struct {
	Ack bool
}

type Distributor struct {
	EventIn chan Event

	conn              *net.UDPConn
	id                int
	eventOuts         []chan Event
	datagramIn        chan Datagram
	timedoutDatagrams chan Datagram
}

type Datagram struct {
	//SequenceNum int
	From  int
	Acks  []int
	Event Event
}

func New(id int) *Distributor {
	d := new(Distributor)

	go d.spin()

	return d
}

func (d *Distributor) spin() {
	for {
		select {
		case event := <-d.EventIn:
			d.handleEvent(event)
		case datagram := <-d.datagramIn:
			d.handleDatagram(datagram)
		case timedoutDatagram := <-d.timedoutDatagrams:
			d.handleTimedoutDatagrams(timedoutDatagram)
		}
	}
}

func (d *Distributor) handleEvent(event Event) {
	d.emit(event)
	d.datagramIn <- Datagram{d.id, []int{d.id}, event}
}

func (d *Distributor) handleDatagram(datagram Datagram) {

}

func (d *Distributor) handleTimedoutDatagrams(datagram Datagram) {

}

func (d *Distributor) emit(event Event) {}

func (d *Distributor) SendDatagram(datagram Datagram) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(datagram)
	if err != nil {
		log.Panic(err)
	}

	n, err := d.conn.WriteToUDP(buf.Bytes(), BROADCAST_ADDR)
	if err != nil {
		log.Panic(err)
	} else {
		log.Debugf("Successfully sent %v bytes to %v", n, BROADCAST_ADDR.String())
	}
}

func (d *Distributor) ReceiveDatagram() (datagram Datagram) {
	buf_size := 2048
	buf := make([]byte, buf_size)

	n, addr, err := d.conn.ReadFrom(buf)
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
