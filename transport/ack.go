package transport

import (
	"bytes"
	"encoding/gob"
)

type ack struct {
	Seq    int
	Origin int
	From   int
}

func (t *Transport) sendAck(datagram datagram) {
	ack := ack{t.seq, datagram.Origin, t.id}	
	t.ackConn.Send <- ack.serialize()
	Logger.Debugf("Sent ack %+v", ack)
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
		return parseAck(data), true
	default:
		return ack{}, false
	}
}

func PanicIf(err error) {
	if err != nil {
		Logger.Panic(err)
	}
}
