package distributor

import (
	"bytes"
	"encoding/gob"

	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/google/go-cmp/cmp"
)

//const PORT = 41875

//var BROADCAST_ADDR = &net.UDPAddr{IP: net.ParseIP("127.255.255.255"), Port: PORT}

type Event interface {
	Ack() bool
}

type state int

const (
	idle state = iota
	distributeInternal
	distributeExternal
)

type Distributor struct {
	In  chan Event
	Out chan Event

	conn           *conn.Conn
	id             int
	eventOuts      []chan Event
	stash          Event
	peers          []int
	state          state
	datagram       datagram
	sequenceNumber int
}

type datagram struct {
	ack            bool
	origin         int
	acks           []int
	event          Event
	sequenceNumber int
}

func (d datagram) serialize() []byte {
	var buf bytes.Buffer
	utils.PanicIf(gob.NewEncoder(&buf).Encode(d))
	return buf.Bytes()
}

func parseDatagram(b []byte) (d datagram) {
	utils.PanicIf(gob.NewDecoder(bytes.NewBuffer(b)).Decode(&d))
	return
}

func New(id int) *Distributor {
	d := new(Distributor)

	go d.spin()

	return d
}

// check if we have acked

func (d *Distributor) spin() {
	for {
		switch d.state {
		case idle:
			select {
			case event := <-d.In:
				datagram := datagram{
					ack:            false,
					origin:         d.id,
					acks:           []int{d.id},
					event:          event,
					sequenceNumber: d.sequenceNumber,
				}
				d.conn.Send <- datagram.serialize()
				d.datagram = datagram
				d.state = distributeInternal

			case data := <-d.conn.Receive:
				datagram := parseDatagram(data)
				if !datagram.ack {
					datagram.acks = append(datagram.acks, d.id)
					d.conn.Send <- datagram.serialize()
					d.datagram = datagram
					d.state = distributeExternal
				}

			}

		case distributeInternal:
			select {
			case data := <-d.conn.Receive:
				datagram := parseDatagram(data)
				if datagram.sequenceNumber == d.sequenceNumber {
					if datagram.ack {
						if datagram.origin == d.id {
							d.datagram.acks = utils.Merge(d.datagram.acks, datagram.acks)
							if utils.Subset(d.peers, datagram.acks) {
								d.Out <- datagram.event
								d.sequenceNumber++
								d.state = idle
							} else {
								//do nothing, wait for timeout or more acks
							}

						} else {
							//then there is someone elses ack, i.e. we dont care
						}

					} else { //!ack
						if datagram.origin < d.id {
							d.stash = d.datagram.event
							datagram.acks = append(datagram.acks, d.id)
							d.conn.Send <- datagram.serialize()
							d.datagram = datagram
							d.state = distributeExternal
						} else {
							//someone elses datagram, but they have less priority than us, i.e. dont care
						}
					}

				} else { // not current sequence number
					if datagram.sequenceNumber < d.sequenceNumber {
						datagram.acks = append(datagram.acks, d.id)
						d.conn.Send <- datagram.serialize()
					} else {
						// ignore if seq is greater than us
					}

				}
			}
		case distributeExternal:
			select {
			case data := <-d.conn.Receive:
				datagram := parseDatagram(data)
				
			}

		}
	}
}

func (d *Distributor) distribute(event Event) {
	datagram := datagram{
		ack:    false,
		origin: d.id,
		acks:   []int{d.id},
		event:  event,
	}
	d.send(datagram)
}

func (d *Distributor) handleNew(datagram datagram) {
	datagram.ack = true
	datagram.acks = append(datagram.acks, d.id)
	if utils.Subset(d.peers, datagram.acks) {
		datagram.event.Ack()
		d.Out <- datagram.event
	}
	d.send(datagram)
}

func (d *Distributor) send(datagram datagram) {
	d.conn.Send <- datagram.serialize()
	for {
		select {
		case b := <-d.conn.Receive:
			newDatagram := parseDatagram(b)
			if !cmp.Equal(newDatagram, datagram) {

			}
		default:
		}
	}
}
