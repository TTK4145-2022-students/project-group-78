package transport

import (
	"testing"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/mocknet"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.Logger.SetLevel(logrus.DebugLevel)
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
	peers := []byte{1, 2}
	t1 := New(1, peers)
	t2 := New(2, peers)
	defer t1.Close()
	defer t2.Close()

	m1 := mocknet.New(config.ACK_PORT)
	m2 := mocknet.New(config.DATAGRAM_PORT)
	defer m1.Close()
	defer m2.Close()

	sent := "Hello"
	t1.Send <- []byte(sent)
	got := string(<-t2.Receive)

	assert.Equal(t, sent, got)
}

func TestTransportWithLoss(t *testing.T) {
	peers := []byte{1, 2}
	t1 := New(1, peers)
	t2 := New(2, peers)
	defer t1.Close()
	defer t2.Close()

	m1 := mocknet.New(config.ACK_PORT)
	m2 := mocknet.New(config.DATAGRAM_PORT)
	m1.LossPercentage <- 50
	m2.LossPercentage <- 50
	defer m1.Close()
	defer m2.Close()

	sent := "Hello"
	t1.Send <- []byte(sent)
	got := string(<-t1.Receive)

	assert.Equal(t, sent, got)
}
