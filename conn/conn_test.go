package conn

import (
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.SetLevel(logrus.DebugLevel)
}

func TestConn(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")

	conn := New(ip, 2001)
	conn2 := New(ip, 2002)
	defer conn.Close()
	defer conn2.Close()

	msg := []byte{1}

	t.Run("Sending to self", func(t *testing.T) {		
		conn.SendTo(msg, ip, 2001)
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case r := <-conn.Receive:
			assert.Equal(t, msg, r)
		case <-timer.C:
			t.Error("timed out")
		}
	})

	t.Run("Sending to other", func(t *testing.T) {
		conn.SendTo(msg, ip, 2002)
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case r := <-conn2.Receive:
			assert.Equal(t, msg, r)
		case <-timer.C:
			t.Error("timed out")
		}
	})
}
