package transport

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

var Logger = utils.NewLogger()

type Transport struct {
	Send    chan []byte
	Receive chan []byte

	id           byte
	state        state
	peers        []byte
	ackConn      *conn.Conn
	datagramConn *conn.Conn
	seq          int
	logger       *logrus.Entry
	closed       *abool.AtomicBool

	message       []byte
	messageOrigin byte
	messageAcks   []byte
	messageSent   time.Time

	stash []byte
}

func New(id byte, peers []byte) *Transport {
	localIp := net.ParseIP(fmt.Sprintf("127.0.0.%v", id))

	t := &Transport{
		Send:    make(chan []byte, 10),
		Receive: make(chan []byte, 10),

		id:           id,
		state:        idle,
		peers:        peers,
		ackConn:      conn.New(localIp, config.ACK_PORT, config.BROADCAST_IP, config.ACK_PORT),
		datagramConn: conn.New(localIp, config.DATAGRAM_PORT, config.BROADCAST_IP, config.DATAGRAM_PORT),
		seq:          1,
		logger:       Logger.WithField("id", id).WithField("pkg", "transport"),
		closed:       abool.New(),
	}

	go t.runForever()

	return t
}

func (t *Transport) Close() {
	t.closed.Set()
	time.Sleep(10 * time.Millisecond)
	t.ackConn.Close()
	t.datagramConn.Close()
}

func (t *Transport) run() {
	switch t.state {
	case idle:
		t.idle()

	case sending:
		t.sending()
	}
}

func (t *Transport) runForever() {
	for t.closed.IsNotSet() {
		t.run()
	}
}

func (t *Transport) getMessage() (message []byte, got bool) {
	if t.stash != nil {
		message = t.stash
		t.stash = nil
		got = true
	} else {
		select {
		case message = <-t.Send:
			got = true
		default:
		}
	}
	return
}

func (t *Transport) sendMessage(message []byte, origin byte) {
	t.message = message
	t.messageOrigin = origin
	t.messageSent = time.Now()

	datagram := datagram{t.seq, origin, message}
	t.datagramConn.Send(datagram.serialize())
	t.logDatagram(datagram).Trace("sent datagram")
}
