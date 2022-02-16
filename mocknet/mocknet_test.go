package mocknet

import (
	"net"
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.SetLevel(logrus.DebugLevel)
}

func TestNoDrop(t *testing.T) {
	port := 41875
	conn := conn.New(net.ParseIP("127.0.0.1"), port, net.ParseIP("127.255.255.255"), port)
	defer conn.Close()

	broadcaster := New(port)
	defer broadcaster.Close()

	val := []byte{1}
	assert.Nil(t, conn.Send(val))
	assert.Equal(t, <-conn.Receive, val)
}

func TestDrop(t *testing.T) {
	port := 41875
	conn := conn.New(net.ParseIP("127.0.0.1"), port, net.ParseIP("127.255.255.255"), port)
	defer conn.Close()

	broadcaster := New(port)
	broadcaster.DropPercentage <- 25
	defer broadcaster.Close()

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
	broadcaster.logger.Info(c)
	assert.False(t, c < 65 || c > 85)
}
