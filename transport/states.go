package transport

import (
	"time"

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
		t.state = sending

	} else if _, got := t.getAck(); got {
		// Acks when idling are irrelevant

	} else if datagram, got := t.getDatagram(); got {
		if datagram.Seq < t.seq {
			// We have already acked this one before, so we will blindly ack it again
			t.sendAck(datagram)
		} else if datagram.Seq > t.seq {
			// We do not handle the future
		} else {
			t.sendAck(datagram)
			if datagram.Origin != t.id {
				t.sendMessage(datagram.Message, datagram.Origin)
				t.state = sending
			}
		}
	}
}

func (t *Transport) sending() {
	if ack, got := t.getAck(); got {
		if ack.Seq == t.seq && ack.Origin == t.messageOrigin && !utils.Member(ack.From, t.messageAcks) && ack.From != t.id {
			t.messageAcks = append(t.messageAcks, ack.From)
			if utils.Subset(t.peers, t.messageAcks) {
				t.Receive <- t.message
				t.seq++
				t.messageAcks = []int{}
				t.state = idle
			}
		}

	} else if datagram, got := t.getDatagram(); got {
		if datagram.Seq < t.seq {
			// We have already acked this one before, so we will blindly ack it again
			t.sendAck(datagram)
		} else if datagram.Seq > t.seq {
			// We do not handle the future
		} else {
			if datagram.Origin < t.messageOrigin {
				// We must yield for the lower origin datagram, i.e. drop the message we are currently sending
				if t.messageOrigin == t.id {
					// If we are currently sending one of our own messages, we need to stash it first before yielding
					if t.stash != nil {
						Logger.Panicf("Stash contained %v when trying to stash %v", t.stash, datagram.Message)
					}
					t.stash = datagram.Message
					Logger.Debugf("Stashed %v", t.stash)

					t.messageAcks = []int{}
				}
				t.sendAck(datagram)
				t.sendMessage(datagram.Message, datagram.Origin)
			} else if datagram.Origin == t.messageOrigin {
				t.sendAck(datagram)
			}
		}

		// TODO: add timeout as param
	} else if time.Now().Sub(t.messageSent) >= 10*time.Millisecond {
		t.sendMessage(t.message, t.messageOrigin)
	}
}
