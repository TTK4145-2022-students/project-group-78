package transport

import (
	"net"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	Logger.SetLevel(log.DebugLevel)
}

func TestAckSerialization(t *testing.T) {
	serialized := ack{1, 2, 3}
	gotBack := parseAck(serialized.serialize())

	if serialized != gotBack {
		t.Error("serialized", serialized, "got back", gotBack)
	}
}

func TestDatagramSerialization(t *testing.T) {
	serialized := datagram{1, 2, []byte("Hello")}
	gotBack := parseDatagram(serialized.serialize())

	if string(serialized.Message) != string(gotBack.Message) {
		t.Error("serialized", serialized, "got back", gotBack)
	}
}

func TestTransport(t *testing.T) {
	// Must have broadcast on 127.255.255.255
	ip := net.ParseIP("127.255.255.255")
	t1 := New(1, ip, []int{2})
	t2 := New(2, ip, []int{1})

	sent := "Hello"
	t1.Send <- []byte(sent)
	got := string(<-t2.Receive)

	if sent != got {
		t.Error("sent", sent, "got", got)
	}
}
