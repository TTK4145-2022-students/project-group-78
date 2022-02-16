package transport

import (
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/utils"
)

type state int

const (
	idle state = iota
	sending
)

func (t *Transport) idle() {
	if message, got := t.getMessage(); got {
		t.sendMessage(message, t.id)
		t.log().Debug("transmitted orginal message")
		t.state = sending		

	} else if ack, got := t.getAck(); got {
		// Acks when idling are irrelevant
		t.logAck(ack).Debug("ignored ack")

	} else if datagram, got := t.getDatagram(); got {
		if datagram.Seq < t.seq {
			// We have already acked this one before, so we will blindly ack it again
			t.sendAck(datagram)
			t.logDatagram(datagram).Debug("re-acked datagram")
		} else if datagram.Seq > t.seq {
			// We do not handle the future
			t.logDatagram(datagram).Debug("ignored datagram")
		} else {
			t.sendAck(datagram)
			t.logDatagram(datagram).Debug("acked")
			if datagram.Origin != t.id {
				t.sendMessage(datagram.Message, datagram.Origin)
				t.logDatagram(datagram).Debug("retransmitted as own")
				t.state = sending				
			}
		}
	}
}

func (t *Transport) sending() {
	if ack, got := t.getAck(); got {
		if ack.Seq == t.seq && ack.Origin == t.messageOrigin && !utils.Member(ack.From, t.messageAcks) {
			t.messageAcks = append(t.messageAcks, ack.From)
			t.logAck(ack).Debug("registered ack")
			if utils.Subset(t.peers, t.messageAcks) {
				t.Receive <- t.message
				t.seq++
				t.messageAcks = []byte{}
				t.log().Debug("finished transmission")
				t.state = idle
			}
		} else {
			t.logAck(ack).Debug("ignored ack")
		}

	} else if datagram, got := t.getDatagram(); got {
		if datagram.Seq < t.seq {
			// We have already acked this one before, so we will blindly ack it again
			t.sendAck(datagram)
			t.logDatagram(datagram).Debug("re-acked datagram")
		} else if datagram.Seq > t.seq {
			// We do not handle the future
			t.logDatagram(datagram).Debug("ignored datagram")
		} else {
			if datagram.Origin < t.messageOrigin {
				t.logDatagram(datagram).Debug("yielding for datagram")
				// We must yield for the lower origin datagram, i.e. drop the message we are currently sending
				if t.messageOrigin == t.id {
					// If we are currently sending one of our own messages, we need to stash it first before yielding
					if t.stash != nil {
						t.log().Panic("stash was not empty")
					}
					t.stash = datagram.Message
					t.log().Debug("stashed message")
				}
				t.messageAcks = []byte{}
				t.log().Debug("resat acks")

				t.sendAck(datagram)
				t.logDatagram(datagram).Debug("acked")

				t.sendMessage(datagram.Message, datagram.Origin)
				t.logDatagram(datagram).Debug("retransmitted as own")
			} else if datagram.Origin == t.messageOrigin {
				t.sendAck(datagram)
				t.logDatagram(datagram).Debug("acked")
			} else {
				t.log().Debug("ignored datagram")
			}
		}

		// TODO: add timeout as param
	} else if time.Now().Sub(t.messageSent) >= config.RETRANSMIT_INTERVAL {
		t.sendMessage(t.message, t.messageOrigin)
		t.log().Debug("retransmitted")
	}
}
