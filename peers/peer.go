package peers

import (
	"fmt"
	"net"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

var TIMEOUT = time.Second
var RESEND = 100 * time.Millisecond

var Logger = utils.NewLogger()

type Peer struct {
	conn   *conn.Conn
	outs   []chan []byte
	times  map[byte]time.Time
	last   []byte
	id     byte
	logger *logrus.Entry
	closed *abool.AtomicBool
}

func New(id byte) *Peer {
	localIp := net.ParseIP(fmt.Sprintf("127.0.0.%v", id))
	p := &Peer{
		conn:   conn.New(localIp, config.HEARTBEAT_PORT, net.ParseIP("127.255.255.255"), config.HEARTBEAT_PORT),
		times:  make(map[byte]time.Time, 1),
		id:     id,
		logger: Logger.WithField("pkg", "peers").WithField("id", id),
		closed: abool.New(),
	}

	go p.sendForever()
	go p.listenForever()

	return p
}

func (p *Peer) Subscribe(out chan []byte) {
	p.outs = append(p.outs, out)
}

func (p *Peer) Close() {
	p.closed.Set()
	p.conn.Close()
}

func (p *Peer) getHeartbeat() (byte, bool) {
	select {
	case b := <-p.conn.Receive:
		id := b[0]
		p.logger.WithField("from", id).Debug("Received heartbeat")
		return id, true
	default:
		return 0, false
	}
}

func (p *Peer) listen() {
	id, got := p.getHeartbeat()
	if got {
		p.times[id] = time.Now()
	}

	peers := make([]byte, 0)
	for id, time_ := range p.times {
		if time.Now().Sub(time_) < TIMEOUT {
			if id == 0 {
				p.logger.Panic()
			}
			peers = append(peers, id)
		}
	}

	if !(utils.Subset(peers, p.last) && utils.Subset(p.last, peers)) {
		p.logger.WithField("now", peers).WithField("past", p.last).Debug("Peers changed")
		for _, out := range p.outs {
			out <- append([]byte{}, peers...) // Go's way of deep copy ...
		}
		p.last = peers
	}
}

func (p *Peer) listenForever() {
	for {
		p.listen()
	}
}

func (p *Peer) sendForever() {
	for p.closed.IsNotSet() {
		p.conn.Send([]byte{p.id})
		p.logger.Debug("Sent heartbeat")
		time.Sleep(RESEND)
	}
}
