package transport

import (
	"bytes"
	"encoding/gob"

	"github.com/sirupsen/logrus"
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
		datagram := parseDatagram(data)
		t.logDatagram(datagram).Debug("received datagram")
		return datagram, true
	default:
		return datagram{}, false
	}
}

func (t *Transport) logDatagram(d datagram) *logrus.Entry {
	return t.logger.WithFields(logrus.Fields{"seq": d.Seq, "origin": d.Origin})
}
