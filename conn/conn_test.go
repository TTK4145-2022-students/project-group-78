package conn

import (
	"net"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.Logger.SetLevel(logrus.DebugLevel)
}

func TestSingleConn(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	port := 41875
	conn := New(ip, port, ip, port)
	defer conn.Close()

	val := []byte{1}
	assert.Nil(t, conn.Send(val))
	assert.Equal(t, <-conn.Receive, val)
}

func TestDoubleConn(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	conn := New(ip, 30001, ip, 30002)
	conn2 := New(ip, 30002, ip, 30001)
	defer conn.Close()
	defer conn2.Close()

	val := []byte{1}
	val2 := []byte{2}
	assert.Nil(t, conn.Send(val))
	assert.Equal(t, <-conn2.Receive, val)
	
	assert.Nil(t, conn2.Send(val2))
	assert.Equal(t, <-conn.Receive, val2)
}
