package transport

import (
	"bytes"
	"encoding/gob"
)

type datagram struct {
	Seq     int
	Origin  byte
	Message []byte
}

func (d datagram) serialize() []byte {
	var buf bytes.Buffer
	PanicIf(gob.NewEncoder(&buf).Encode(d))
	return buf.Bytes()
}

func parseDatagram(b []byte) (d datagram) {
	PanicIf(gob.NewDecoder(bytes.NewBuffer(b)).Decode(&d))
	return
}

func (t *Transport) getDatagram() (datagram, bool) {
	select {
	case data := <-t.datagramConn.Receive:
		return parseDatagram(data), true
	default:
		return datagram{}, false
	}
}
