package conn

import (
	"net"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "15:04:05"})
}

func TestSendAndReceive(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	port := 41875
	conn := New(ip, port, ip, port)

	sent := "Hello"

	conn.Send([]byte(sent))
	got := string(conn.Receive())

	if sent != got {
		t.Error("sent", sent, "got", got)
	}
}
