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

func TestMocknet(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	conn := conn.New(ip, 2001)
	mocknet := New(2001)
	defer conn.Close()
	defer mocknet.Close()

	msg := []byte{1}

	t.Run("Without packet loss", func(t *testing.T) {
		conn.SendTo(msg, config.BROADCAST_IP, 2001)
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case r := <-conn.Receive:
			assert.Equal(t, msg, r)
		case <-timer.C:
			t.Error("timed out")
		}
	})

	t.Run("With packet loss", func(t *testing.T) {
		mocknet.SetLossPercentage(50) // 50 % chance of loss each way from mocknet, means 75 % total loss to and from the mocknet
		for i := 0; i < 100; i++ {
			conn.SendTo(msg, config.BROADCAST_IP, 2001)
		}

		c := 0
		timeout := 1000 * time.Millisecond
		timer := time.NewTimer(timeout)
		for {
			select {
			case <-conn.Receive:
				c++
				timer.Reset(timeout)

			case <-timer.C:
				if c < 15 || c > 35 {
					t.Errorf("received %v messages, should have been around 25", c)
				}
				return
			}
		}
	})
}
