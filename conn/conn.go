package conn

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/TTK4145-2022-students/project-group-78/utils"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"
)

var Logger = utils.NewLogger("conn")

const MAX_PACKET_SIZE = 1024

type Conn struct {
	Receive chan []byte

	conn       *net.UDPConn
	remoteAddr *net.UDPAddr
	logger     *logrus.Entry
	closed     *abool.AtomicBool
}

func New(localIp net.IP, localPort int, remoteIp net.IP, remotePort int) *Conn {
	localAddr := &net.UDPAddr{IP: localIp, Port: localPort}
	c := &Conn{
		Receive:    make(chan []byte, 128),
		remoteAddr: &net.UDPAddr{IP: remoteIp, Port: remotePort},
		logger:     Logger.WithField("connAddr", localAddr.String()).WithField("pkg", "conn"),
		closed:     abool.New(),
	}

	var err error
	c.conn, err = net.ListenUDP("udp", localAddr)
	if err != nil {
		c.logger.Panic(err)
	} else {
		c.logger.Debug("Listening")
	}

	go c.receiveForever()

	return c
}

func (c *Conn) receiveForever() {
	for c.closed.IsNotSet() {
		packet := make([]byte, MAX_PACKET_SIZE)
		n, addr, err := c.conn.ReadFromUDP(packet)
		if err != nil {
			if c.closed.IsSet() {
				return
			} else {
				c.logger.Panic(err)
			}
		} else {
			c.logger.WithField("from", addr.String()).WithField("size", n).Debug("Received")
		}
		c.Receive <- packet[0:n]
	}
}

func (c *Conn) SendTo(packet []byte, remoteAddr *net.UDPAddr) error {
	if len(packet) > MAX_PACKET_SIZE {
		return errors.New(fmt.Sprintf("packet size (%v) cannot exceed %v", len(packet), MAX_PACKET_SIZE))
	}

	n, err := c.conn.WriteToUDP(packet, remoteAddr)
	if err == nil {
		c.logger.WithField("to", remoteAddr.String()).WithField("size", n).Debug("Sent")
	} else {
		c.logger.WithField("to", remoteAddr.String()).WithField("size", n).Panic(err)
	}
	return err
}

func (c *Conn) Send(packet []byte) error {
	return c.SendTo(packet, c.remoteAddr)
}

func (c *Conn) Close() {
	c.closed.Set()
	time.Sleep(10 * time.Millisecond)
	err := c.conn.Close()
	if err != nil {
		c.logger.Panic(err)
	} else {
		c.logger.Debugf("Closed")
	}
}
