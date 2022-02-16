package peers

import (
	"fmt"
	"net"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
)

var TIMEOUT = time.Second
var RESEND = 100 * time.Millisecond

type Peer struct {
	conn  *conn.Conn
	outs  []chan []byte
	times map[byte]time.Time
	last  []byte
	id    byte
}

func New(id byte) *Peer {
	localIp := net.ParseIP(fmt.Sprintf("127.0.0.%v", id))
	port := 61523
	return &Peer{
		conn: conn.New(localIp, port, net.ParseIP("127.255.255.255"), port),
		id:   id,
	}

}

func (p *Peer) Subscribe(out chan []byte) {
	p.outs = append(p.outs, out)
}

func (p *Peer) listen() {
	select {
	case b := <-p.conn.Receive:
		id := b[0]
		p.times[id] = time.Now()
	}

	peers := make([]byte, len(p.last))
	for id, time_ := range p.times {
		if time.Now().Sub(time_) < TIMEOUT {
			peers = append(peers, id)
		}
	}
	if utils.Equal(peers, p.last) {
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
	for {
		p.conn.Send([]byte{p.id})
		time.Sleep(RESEND)
	}
}
