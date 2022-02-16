package mocknet

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

var Logger = utils.NewLogger("mocknet")

type Mocknet struct {
	LossPercentage chan int

	conn           *conn.Conn
	logger         *logrus.Entry
	closed         *abool.AtomicBool
	lossPercentage int
	port           int
}

func New(port int) *Mocknet {
	rand.Seed(0)
	broadcastAddr := &net.UDPAddr{IP: config.BROADCAST_IP, Port: port}
	p := &Mocknet{
		LossPercentage: make(chan int, 1),
		conn:           conn.New(config.BROADCAST_IP, port, nil, 0),
		logger:         Logger.WithField("addr", broadcastAddr.String()).WithField("pkg", "mocknet"),
		closed:         abool.New(),
		port:           port,
	}

	go p.runForever()

	return p
}

func (m *Mocknet) Close() {
	m.closed.Set()
	time.Sleep(10 * time.Millisecond)
	m.conn.Close()
}

func (m *Mocknet) getLosePercentage() int {
	select {
	case lossPercentage := <-m.LossPercentage:
		if lossPercentage < 0 || lossPercentage > 100 {
			m.logger.WithField("lossPercentage", lossPercentage).Panic("invalid loss percentage! Must be in interval [0, 100]")
		} else {
			m.lossPercentage = lossPercentage
		}
	default:
	}
	return m.lossPercentage
}

func (m *Mocknet) shouldLose() bool {
	return rand.Intn(100) <= m.getLosePercentage()
}

func (m *Mocknet) run() {
	select {
	case msg := <-m.conn.Receive:
		if m.shouldLose() {
			return
		}
		for i := 0; i < 255 && m.closed.IsNotSet(); i++ {
			addr := &net.UDPAddr{IP: net.ParseIP(fmt.Sprintf("127.0.0.%v", i)), Port: m.port}
			if !m.shouldLose() {
				m.conn.SendTo(msg, addr)
			}

		}
	default:
	}
}

func (m *Mocknet) runForever() {
	for m.closed.IsNotSet() {
		m.run()
	}
}
