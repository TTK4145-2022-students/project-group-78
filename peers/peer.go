package peers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/elliotchance/pie/pie"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

var Logger = utils.NewLogger("peer", "id")

type Peer struct {
	Peers chan pie.Ints

	conn   *conn.Conn
	id     int
	closed *abool.AtomicBool
	times  map[int]time.Time
	peers  pie.Ints
}

func New(id int) *Peer {
	p := &Peer{
		Peers:  make(chan pie.Ints, 16),
		conn:   conn.New(config.LocalIp(id), config.HEARTBEAT_PORT),
		id:     id,
		closed: abool.New(),
		times:  make(map[int]time.Time, 1),
	}

	go p.send()
	go p.listen()

	return p
}

func (p *Peer) Close() {
	p.closed.Set()
	time.Sleep(time.Millisecond)
	p.conn.Close()
}

func (p *Peer) log() *logrus.Entry {
	return Logger.WithField("id", p.id)
}

func (p *Peer) getHeartbeat() (int, bool) {
	select {
	case b := <-p.conn.Receive:
		id, err := strconv.Atoi(string(b))
		if err != nil {
			p.log().Error(err)
			return 0, false
		} else {
			p.log().WithField("from", id).Debug("received heartbeat")
			return id, true
		}
	default:
		return 0, false
	}
}

func (p *Peer) listen() {
	for p.closed.IsNotSet() {
		id, got := p.getHeartbeat()
		if got {
			p.times[id] = time.Now()
		}

		currentPeers := pie.Ints{}
		for id, time_ := range p.times {
			if time.Now().Sub(time_) < config.TRANSMISSION_TIMEOUT {
				currentPeers = currentPeers.Append(id)
			}
		}

		add, remove := currentPeers.Diff(p.peers)
		if len(add) != 0 || len(remove) != 0 {
			p.log().WithField("peers", currentPeers).Info("peers changed")
			p.Peers <- currentPeers.Append() //Deep copy
			p.peers = currentPeers
		}
	}
}

func (p *Peer) send() {
	for p.closed.IsNotSet() {
		p.conn.SendTo([]byte(string(fmt.Sprint(p.id))), config.BROADCAST_IP, config.HEARTBEAT_PORT)
		p.log().Debug("sent heartbeat")
		time.Sleep(config.RETRANSMIT_INTERVAL)
	}
}
