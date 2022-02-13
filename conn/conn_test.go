package conn

import (
	"net"
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	Logger.SetLevel(logrus.DebugLevel)
}

func TestSendAndReceive(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	port := 41875
	conn := New(ip, port, ip, port)

	sent := "Hello"
	conn.Send <- []byte(sent)
	got := string(<-conn.Receive)

	if sent != got {
		t.Error("sent", sent, "got", got)
	}
}
