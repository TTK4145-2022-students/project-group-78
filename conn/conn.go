package conn

import (
	"net"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

var Logger = utils.NewLogger("conn", "addr")

const MAX_PACKET_SIZE = 1024

type Conn struct {
	Receive chan []byte

	conn   *net.UDPConn
	addr   *net.UDPAddr
	closed *abool.AtomicBool
}

func New(ip net.IP, port int) *Conn {
	addr := &net.UDPAddr{IP: ip, Port: port}
	conn, err := net.ListenUDP("udp", addr)

	c := &Conn{
		Receive: make(chan []byte, 128),
		conn:    conn,
		addr:    addr,
		closed:  abool.New(),
	}

	if err != nil {
		c.log().Panic(err)
	} else {
		c.log().Info("listening")
	}

	go c.receive()

	return c
}

func (c *Conn) log() *logrus.Entry {
	return Logger.WithField("addr", c.addr)
}

func (c *Conn) receive() {
	for c.closed.IsNotSet() {
		packet := make([]byte, MAX_PACKET_SIZE)
		n, addr, err := c.conn.ReadFromUDP(packet)
		if err != nil {
			if c.closed.IsNotSet() {
				c.log().Panic(err)
			}
		} else {
			c.Receive <- packet[0:n]
			c.log().WithFields(logrus.Fields{
				"size": n,
				"from": addr,
			}).Debug("received")
		}
	}
}

func (c *Conn) SendTo(packet []byte, ip net.IP, port int) {
	if len(packet) > MAX_PACKET_SIZE {
		c.log().WithFields(logrus.Fields{
			"size":    len(packet),
			"maxSize": MAX_PACKET_SIZE,
		}).Panic("too large")
	}

	to := &net.UDPAddr{IP: ip, Port: port}
	n, err := c.conn.WriteToUDP(packet, to)
	logger := c.log().WithFields(logrus.Fields{
		"to":   to,
		"size": n,
	})

	if err != nil {
		logger.Panic(err)
	} else {
		logger.Debug("sent")
	}
}

func (c *Conn) Close() {
	c.closed.Set()
	time.Sleep(time.Millisecond)
	err := c.conn.Close()
	if err != nil {
		c.log().Panic(err)
	} else {
		c.log().Info("closed")
	}
}
