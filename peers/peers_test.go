package peers

import (
	"errors"
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/mocknet"
	"github.com/elliotchance/pie/pie"
	"github.com/stretchr/testify/assert"
)

func init() {
	//Logger.SetLevel(logrus.DebugLevel)
}

func getPeers(p *Peer) (pie.Ints, error) {
	timer := time.NewTimer(config.TRANSMISSION_TIMEOUT)
	select {
	case p := <-p.Peers:
		return p, nil
	case <-timer.C:
		return pie.Ints{}, errors.New("timed out")
	}
}

func TestPeer(t *testing.T) {
	p1 := New(1)
	p2 := New(2)
	defer p1.Close()

	mocknet := mocknet.New(config.HEARTBEAT_PORT)
	defer mocknet.Close()

	t.Run("Without packet loss", func(t *testing.T) {
		p, err := getPeers(p1)
		assert.Nil(t, err)
		assert.True(t, p.Equals(pie.Ints{1}) || p.Equals(pie.Ints{2}))

		p, err = getPeers(p1)
		assert.Nil(t, err)
		assert.True(t, p.Equals(pie.Ints{1, 2}) || p.Equals(pie.Ints{2, 1}))
	})

	t.Run("With packet loss", func(t *testing.T) {
		mocknet.SetLossPercentage(50)
		_, err := getPeers(p1)
		assert.NotNil(t, err)
	})

	t.Run("Drop peer", func(t *testing.T) {
		p2.Close()
		time.Sleep(config.TRANSMISSION_TIMEOUT)
		p, err := getPeers(p1)
		assert.Nil(t, err)
		assert.True(t, p.Equals(pie.Ints{1}) || p.Equals(pie.Ints{2}))
	})
}
