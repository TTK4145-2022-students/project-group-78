package transport

import (
	"bytes"
	"encoding/gob"

	"github.com/sirupsen/logrus"
)

type ack struct {
	Seq    int
	Origin byte
	From   byte
}

func (t *Transport) sendAck(datagram datagram) {
	ack := ack{datagram.Seq, datagram.Origin, t.id}
	t.ackConn.Send(ack.serialize())
	t.logAck(ack).Trace("sent ack")
}

func parseAck(b []byte) (a ack) {
	PanicIf(gob.NewDecoder(bytes.NewBuffer(b)).Decode(&a))
	return a
}

func (a ack) serialize() []byte {
	var buf bytes.Buffer
	PanicIf(gob.NewEncoder(&buf).Encode(a))
	return buf.Bytes()
}

func (t *Transport) getAck() (ack, bool) {
	select {
	case data := <-t.ackConn.Receive:
		ack := parseAck(data)
		t.logAck(ack).Debug("received ack")
		return ack, true
	default:
		return ack{}, false
	}
}

func (t *Transport) logAck(ack ack) *logrus.Entry {
	return t.logger.WithFields(logrus.Fields{"seq": ack.Seq, "origin": ack.Origin, "from": ack.From})
}

func PanicIf(err error) {
	if err != nil {
		Logger.Panic(err)
	}
}
