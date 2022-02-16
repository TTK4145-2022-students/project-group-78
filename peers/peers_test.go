package peers

import (
	"testing"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/mocknet"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	Logger.SetLevel(logrus.DebugLevel)
}

func TestPeer(t *testing.T) {
	p1 := New(1)
	p2 := New(2)
	defer p1.Close()
	defer p2.Close()

	broadcaster := mocknet.New(config.HEARTBEAT_PORT)
	defer broadcaster.Close()

	peers := make(chan []byte, 10)
	p1.Subscribe(peers)

	time.Sleep(10 * time.Millisecond)

	var p []byte
	for len(peers) != 0 {
		select {
		case p = <-peers:
		default:
		}
	}
	assert.True(t, utils.Subset(p, []byte{1, 2}) && utils.Subset([]byte{1, 2}, p))
}

func TestPeerWithPacketLoss(t *testing.T) {
	p1 := New(1)
	p2 := New(2)
	defer p1.Close()
	defer p2.Close()

	broadcaster := mocknet.New(config.HEARTBEAT_PORT)
	broadcaster.DropPercentage <- 50
	defer broadcaster.Close()

	peers := make(chan []byte, 10)
	p1.Subscribe(peers)

	time.Sleep(1000 * time.Millisecond)

	var p []byte
	for len(peers) != 0 {
		select {
		case p = <-peers:
		default:
		}
	}
	assert.True(t, utils.Subset(p, []byte{1, 2}) && utils.Subset([]byte{1, 2}, p))
}
