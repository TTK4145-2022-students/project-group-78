package mocknet

import (
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/config"
	"github.com/TTK4145-2022-students/project-group-78/conn"
	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
)

var Logger = utils.NewLogger("mocknet", "port")

type Mocknet struct {
	conn           *conn.Conn
	port           int
	quit           chan bool
	lossPercentage *int32
}

func New(port int) *Mocknet {
	rand.Seed(0)
	p := &Mocknet{
		conn:           conn.New(config.BROADCAST_IP, port),
		port:           port,
		quit:           make(chan bool),
		lossPercentage: new(int32),
	}

	go p.run()

	return p
}

func (m *Mocknet) Close() {
	m.quit <- true
	time.Sleep(10 * time.Millisecond)
	m.conn.Close()
}

func (m *Mocknet) SetLossPercentage(percentage int) {
	atomic.StoreInt32(m.lossPercentage, int32(percentage))
}

func (m *Mocknet) log() *logrus.Entry {
	return Logger.WithField("port", m.port)
}

func (m *Mocknet) shouldLose() bool {
	return rand.Intn(100) <= int(atomic.LoadInt32(m.lossPercentage))
}

func (m *Mocknet) run() {
	for {
		select {
		case msg := <-m.conn.Receive:
			if !m.shouldLose() {
				for i := 1; i <= 255; i++ {
					ip := net.ParseIP(fmt.Sprintf("127.0.0.%v", i))
					if !m.shouldLose() {
						m.conn.SendTo(msg, ip, m.port)
					}
				}
			}
		case <-m.quit:
			return
		}
	}
}
