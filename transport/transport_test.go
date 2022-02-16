package transport

import (
	"testing"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/mocknet"
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
	t1 := New(1, []int{2})
	t2 := New(2, []int{1})

	b1 := mocknet.New(config.ACK_PORT)
	b2 := mocknet.New(config.DATAGRAM_PORT)
	defer b1.Close()
	defer b2.Close()

	sent := "Hello"
	t1.Send <- []byte(sent)
	got := string(<-t2.Receive)

	if sent != got {
		t.Error("sent", sent, "got", got)
	}
}
