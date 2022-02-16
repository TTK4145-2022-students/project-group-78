package mocknet

import (
	"net"
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.SetLevel(logrus.DebugLevel)
}

func TestNoDrop(t *testing.T) {
	port := 41875
	conn := conn.New(net.ParseIP("127.0.0.1"), port, config.BROADCAST_IP, port)
	defer conn.Close()

	mocknet := New(port)
	defer mocknet.Close()

	time.Sleep(10 * time.Microsecond)

	val := []byte{1}
	assert.Nil(t, conn.Send(val))
	assert.Equal(t, <-conn.Receive, val)
}

func TestDrop(t *testing.T) {
	port := 41875
	conn := conn.New(net.ParseIP("127.0.0.1"), port, config.BROADCAST_IP, port)
	defer conn.Close()

	mocknet := New(port)
	mocknet.LossPercentage <- 50
	defer mocknet.Close()

	val := []byte{1}
	for i := 0; i < 100; i++ {
		assert.Nil(t, conn.Send(val))
	}

	c := 0
	start := time.Now()
	for time.Now().Sub(start) < 500*time.Millisecond {
		select {
		case <-conn.Receive:
			c++
		default:
		}
	}
	assert.False(t, c < 15 || c > 35)
}
