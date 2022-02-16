package transport

import (
	"fmt"
	"net"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
)

var Logger = utils.NewLogger()

// TODO: add logging

type Transport struct {
	Send    chan []byte
	Receive chan []byte

	id           int
	state        state
	peers        []int
	ackConn      *conn.Conn
	datagramConn *conn.Conn
	seq          int

	message       []byte
	messageOrigin int
	messageAcks   []int
	messageSent   time.Time

	stash []byte
}

func New(id int, peers []int) *Transport {
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
	}

	go t.runForever()

	return t
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
	for {
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

func (t *Transport) sendMessage(message []byte, origin int) {
	t.message = message
	t.messageOrigin = origin
	t.messageSent = time.Now()

	datagram := datagram{t.seq, origin, message}
	t.datagramConn.Send(datagram.serialize())
	Logger.Debugf("Sent datagram %+v", datagram)
}
