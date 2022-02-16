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
	logger := t.logger.WithField("state", "idle")
	if message, got := t.getMessage(); got {
		t.sendMessage(message, t.id)
		t.state = sending
		logger.Debug("transmitted message")

	} else if _, got := t.getAck(); got {
		logger.Debug("ignored ack")
		// Acks when idling are irrelevant

	} else if datagram, got := t.getDatagram(); got {
		if datagram.Seq < t.seq {
			// We have already acked this one before, so we will blindly ack it again
			t.sendAck(datagram)
			logger.Debug("re-acked")
		} else if datagram.Seq > t.seq {
			logger.Debug("ignored datagram")
			// We do not handle the future
		} else {			
			t.sendAck(datagram)
			logger.Debug("acked")
			if datagram.Origin != t.id {
				t.sendMessage(datagram.Message, datagram.Origin)
				logger.Debug("transmitted message")
				t.state = sending
			}
		}
	}
}

func (t *Transport) sending() {
	logger := t.logger.WithField("state", "sending")
	if ack, got := t.getAck(); got {
		if ack.Seq == t.seq && ack.Origin == t.messageOrigin && !utils.Member(ack.From, t.messageAcks) {
			t.messageAcks = append(t.messageAcks, ack.From)
			logger.Debug("registered ack")
			if utils.Subset(t.peers, t.messageAcks) {
				t.Receive <- t.message
				t.seq++
				t.messageAcks = []byte{}
				logger.WithField("seq", t.seq).Debug("finished transmission")
				t.state = idle
			}
		} else {
			logger.Debug("ignored ack")
		}

	} else if datagram, got := t.getDatagram(); got {
		if datagram.Seq < t.seq {
			// We have already acked this one before, so we will blindly ack it again
			t.sendAck(datagram)
			logger.Debug("re-acked")
		} else if datagram.Seq > t.seq {
			logger.Debug("ignored datagram")
			// We do not handle the future
		} else {
			if datagram.Origin < t.messageOrigin {
				logger.Debug("yielding")
				// We must yield for the lower origin datagram, i.e. drop the message we are currently sending
				if t.messageOrigin == t.id {
					// If we are currently sending one of our own messages, we need to stash it first before yielding
					if t.stash != nil {
						logger.Panic("stash was not empty")
					}
					t.stash = datagram.Message
					logger.Debug("stashed message")

					t.messageAcks = []byte{}
					logger.Debug("resat acks")
				}
				t.sendAck(datagram)
				logger.Debug("acked")
				t.sendMessage(datagram.Message, datagram.Origin)
				logger.Debug("transmitted")
			} else if datagram.Origin == t.messageOrigin {
				t.sendAck(datagram)
				logger.Debug("acked")
			}
		}

		// TODO: add timeout as param
	} else if time.Now().Sub(t.messageSent) >= config.RETRANSMIT_INTERVAL {		
		t.sendMessage(t.message, t.messageOrigin)
		t.logger.Trace("retransmitted")
	}
}
