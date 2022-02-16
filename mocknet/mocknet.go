package mocknet

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

var Logger = utils.NewLogger()

type Broadcaster struct {
	DropPercentage chan int

	conn           *conn.Conn
	logger         *logrus.Entry
	closed         *abool.AtomicBool
	dropPercentage int
	port           int
}

func New(port int) *Broadcaster {
	rand.Seed(0)

	ip := net.ParseIP("127.255.255.255")
	p := &Broadcaster{
		DropPercentage: make(chan int, 1),
		conn:           conn.New(ip, port, nil, 0),
		logger:         Logger.WithField("addr", (&net.UDPAddr{IP: ip, Port: port}).String()).WithField("pkg", "mocknet"),
		closed:         abool.New(),
		port:           port,
	}

	go p.runForever()

	return p
}

func (b *Broadcaster) Close() {
	b.closed.Set()
	b.conn.Close()
}

func (b *Broadcaster) getDropPercentage() int {
	select {
	case dropPercentage := <-b.DropPercentage:
		if dropPercentage < 0 || dropPercentage > 100 {
			b.logger.WithField("dropPercentage", dropPercentage).Panic("invalid drop percentage! Must be in interval [0, 100]")
		} else {
			b.dropPercentage = dropPercentage
		}
	default:
	}
	return b.dropPercentage
}

func (b *Broadcaster) shouldDrop() bool {
	return rand.Intn(100) <= b.getDropPercentage()
}

func (b *Broadcaster) run() {
	select {
	case msg := <-b.conn.Receive:
		if b.shouldDrop() {
			return
		}
		for i := 0; i < 255 && b.closed.IsNotSet(); i++ {
			addr := &net.UDPAddr{IP: net.ParseIP(fmt.Sprintf("127.0.0.%v", i)), Port: b.port}
			b.conn.SendTo(msg, addr)
		}
	}
}

func (b *Broadcaster) runForever() {
	for b.closed.IsNotSet() {
		b.run()
	}
}
